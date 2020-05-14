package lygo_db_sync

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/pkg/errors"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const DRIVER_ARANGO = "arango"

//-- DB-SYNC FIELDS --//
const (
	FLD_UUID      = "_db_sync_uuid" // client UID
	FLD_FLAG      = "_db_sync_flag" // if equals true, record is synchronized
	FLD_TIMESTAMP = "_db_sync_timestamp"
)

//-- ERRORS --//
var (
	MissingDatabaseConfigurationError = errors.New("missing_database_configuration")
	ConnectionIsClosedError           = errors.New("connection_is_closed")
	DriverNotFoundError               = errors.New("driver_not_found")
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------
type DBSyncNetConnection interface {
	IsOpen() bool
}

type DBSyncDriver interface {
	Open() error
	Close() error
	BuildQuery(collection string, filter string, params map[string]interface{}) string
	BuildQueryReverse(collection string, filter string, params map[string]interface{}) string
	Collection(database, collection string) (bool, error)
	Execute(database string, query string, params map[string]interface{}) ([]interface{}, error)
	SetNeedUpdateFlag(database, collection string, raw_entity interface{}) error
	Merge(database, collection string, item map[string]interface{}) (map[string]interface{}, error)
}

type DBSyncMessage struct {
	UID        string      `json:"uid"` //  client uid
	Driver     string      `json:"driver"`
	Database   string      `json:"database"`
	Collection string      `json:"collection"`
	Data       interface{} `json:"data"`
}

type DBSync struct {
	UID string

	//-- private --//
	dbConfig     *DBSyncDatabaseConfig
	config       *DBSyncConfigSync
	ticker       *lygo_events.EventTicker
	conn         DBSyncNetConnection
	errorHandler func(sender *DBSync, err error)
	syncHandler  func(message *DBSyncMessage) []map[string]interface{}
	mux          *sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSync(uid string, conn DBSyncNetConnection, dbConfig *DBSyncDatabaseConfig, config *DBSyncConfigSync, mux *sync.Mutex) *DBSync {
	instance := new(DBSync)
	instance.UID = uid
	instance.dbConfig = dbConfig
	instance.config = config
	instance.conn = conn
	instance.ticker = lygo_events.NewEventTicker(config.Interval*time.Second, instance.onLoop)
	instance.mux = mux

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSync) Open() error {
	if nil != instance {
		instance.ticker.Start()
	}
	return nil
}

func (instance *DBSync) Close() error {
	if nil != instance {
		instance.ticker.Stop()
	}
	return nil
}

func (instance *DBSync) Pause() error {
	if nil != instance {
		instance.ticker.Pause()
	}
	return nil
}

func (instance *DBSync) Resume() error {
	if nil != instance {
		instance.ticker.Resume()
	}
	return nil
}

func (instance *DBSync) ReverseSync() (int64, error) {
	if nil != instance {
		return instance.reverseSync()
	}
	return 0, nil
}

func (instance *DBSync) OnError(callback func(sender *DBSync, err error)) {
	if nil != instance {
		instance.errorHandler = callback
	}
}
func (instance *DBSync) OnSync(callback func(message *DBSyncMessage) []map[string]interface{}) {
	if nil != instance {
		instance.syncHandler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSync) triggerErrorAsync(context, message string) {
	if nil != instance.errorHandler {
		go instance.errorHandler(instance, errors.New("["+context+"] "+message))
	} else {
		fmt.Println(message)
	}
}

func (instance *DBSync) onLoop(ticker *lygo_events.EventTicker) {
	// RECOVERY
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			message := lygo_strings.Format("TICKER %s ERROR: %s", instance.config.Uid, r)
			instance.triggerErrorAsync("lygo_db_sync.onLoop", message)
		}
	}()

	// synchronize
	instance.mux.Lock()
	defer instance.mux.Unlock()

	// only if connection is open
	if !instance.conn.IsOpen() {
		return
	}

	driverName := instance.dbConfig.Driver
	driver := GetDriver(driverName, instance.dbConfig)
	if nil != driver {
		err := driver.Open()
		if nil == err {
			localDb := instance.config.LocalDBName
			for _, action := range instance.config.Actions {
				localCollection := action.LocalCollection
				filter := action.Filter
				// ensure collection exists
				if _, err := driver.Collection(localDb, localCollection); nil == err {
					// prepare query
					params := map[string]interface{}{
						"uuid":      instance.UID,
						"timestamp": time.Now().Unix(),
					}
					localQuery := driver.BuildQuery(localCollection, filter, params)
					localData, err := driver.Execute(localDb, localQuery, params)
					if nil != err {
						instance.triggerErrorAsync("driver.Execute", err.Error())
					} else if len(localData) > 0 {
						for _, item := range localData {
							syncResp := instance.sync(driverName, instance.config.RemoteDBName, action.RemoteCollection, item)
							if nil == syncResp || len(syncResp) == 0 {
								// ROLLBACK
								// sync error
								// should rollback transaction
								err = driver.SetNeedUpdateFlag(localDb, localCollection, item)
								if nil != err {
									instance.triggerErrorAsync("sync.Rollback", err.Error())
								}
							} else {
								// UPDATE LOCAL
								for _, entity := range syncResp {
									_, err = driver.Merge(localDb, localCollection, entity)
									if nil != err {
										instance.triggerErrorAsync("driver.Merge", err.Error())
									}
								}
							}
						}
						// fmt.Println(instance.UID, len(localData))
					}
				}
			}
		} else {
			instance.triggerErrorAsync("driver.Open", err.Error())
		}
	}
}

func (instance *DBSync) sync(driver, remoteDatabase, remoteCollection string, data interface{}) []map[string]interface{} {
	if nil != instance.syncHandler {
		message := &DBSyncMessage{
			UID:        instance.UID,
			Driver:     driver,
			Database:   remoteDatabase,
			Collection: remoteCollection,
			Data:       data,
		}
		return instance.syncHandler(message)
	}
	return nil
}

func (instance *DBSync) syncBack(driver, remoteDatabase, remoteCollection string, remoteQuery string) []map[string]interface{} {
	return instance.sync(driver, remoteDatabase, remoteCollection, remoteQuery)
}

func (instance *DBSync) reverseSync() (int64, error) {
	var count int64 = 0
	// synchronize
	instance.mux.Lock()
	defer instance.mux.Unlock()

	// only if connection is open
	if !instance.conn.IsOpen() {
		return 0, ConnectionIsClosedError
	}

	driverName := instance.dbConfig.Driver
	driver := GetDriver(driverName, instance.dbConfig)
	if nil != driver {
		err := driver.Open()
		if nil == err {
			localDb := instance.config.LocalDBName
			for _, action := range instance.config.Actions {
				localCollection := action.LocalCollection
				filter := action.Filter
				// ensure collection exists
				if _, err := driver.Collection(localDb, localCollection); nil == err {
					params := map[string]interface{}{
						"uuid":  instance.UID,
						"skip":  0,
						"limit": 100,
					}
					skip := 0
					for {
						params["skip"] = skip
						skip += 100
						// paged query
						remoteQuery := driver.BuildQueryReverse(action.RemoteCollection, filter, params)
						remoteData := instance.syncBack(driverName, instance.config.RemoteDBName, action.RemoteCollection, remoteQuery)
						if nil != err || nil == remoteData || len(remoteData) == 0 {
							break
						}
						// update items
						for _, entity := range remoteData {
							if nil != entity {
								_, err = driver.Merge(localDb, localCollection, entity)
								if nil != err {
									instance.triggerErrorAsync("driver.Merge", err.Error())
									break
								} else {
									count++
								}
							}
						}
					}
				}
			}
		} else {
			return count, err
		}
	} else {
		return count, errors.WithMessage(DriverNotFoundError, driverName)
	}

	return count, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func GetDriver(name string, config *DBSyncDatabaseConfig) DBSyncDriver {
	switch name {
	case DRIVER_ARANGO:
		return DBSyncDriver(NewDBSyncDriverArango(config))
	}
	return nil
}
