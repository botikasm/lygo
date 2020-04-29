package lygo_db_sync

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/pkg/errors"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const DRIVER_ARANGO = "arango"

//-- DB-SYNC FIELDS --//
const (
	FLD_FLAG      = "_db_sync_flag" // if equals true, record is synchronized
	FLD_TIMESTAMP = "_db_sync_timestamp"
)

//-- ERRORS --//
var (
	MissingDatabaseConfigurationError = errors.New("missing_database_configuration")
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
	Collection(database, collection string) (bool, error)
	Execute(database string, query string, params map[string]interface{}) ([]interface{}, error)
}

type DBSyncMessage struct {
	UID        string
	Driver     string
	Database   string
	Collection string
	UniqueKey  []string
	Data       interface{}
}

type DBSync struct {
	UID string

	//-- private --//
	dbConfig     *DBSyncDatabaseConfig
	config       *DBSyncConfigSync
	ticker       *lygo_events.EventTicker
	conn         DBSyncNetConnection
	errorHandler func(sender *DBSync, err error)
	syncHandler  func(message *DBSyncMessage)
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSync(uid string, conn DBSyncNetConnection, dbConfig *DBSyncDatabaseConfig, config *DBSyncConfigSync) *DBSync {
	instance := new(DBSync)
	instance.UID = uid
	instance.dbConfig = dbConfig
	instance.config = config
	instance.conn = conn
	instance.ticker = lygo_events.NewEventTicker(config.Interval*time.Second, instance.onLoop)

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

func (instance *DBSync) OnError(callback func(sender *DBSync, err error)) {
	if nil != instance {
		instance.errorHandler = callback
	}
}
func (instance *DBSync) OnSync(callback func(message *DBSyncMessage)) {
	if nil != instance {
		instance.syncHandler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSync) triggerError(context, message string) {
	if nil != instance.errorHandler {
		instance.errorHandler(instance, errors.New("["+context+"] "+message))
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
			instance.triggerError("lygo_db_sync.onLoop", message)
		}
	}()

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
				uniqueKey := action.UniqueKey // unique key index (need to update an existing record)
				// ensure collection exists
				if _, err := driver.Collection(localDb, localCollection); nil == err {
					// prepare query
					params := map[string]interface{}{
						"timestamp": time.Now().Unix(),
					}
					localQuery := driver.BuildQuery(localCollection, filter, params)
					localData, err := driver.Execute(localDb, localQuery, params)
					if nil != err {
						instance.triggerError("driver.Execute", err.Error())
					} else if len(localData) > 0 {
						for _, item := range localData {
							instance.sync(driverName, instance.config.RemoteDBName, action.RemoteCollection, uniqueKey, item)
						}
						// fmt.Println(instance.UID, len(localData))
					}
				}
			}
		} else {
			instance.triggerError("driver.Open", err.Error())
		}
	}
}

func (instance *DBSync) sync(driver, remoteDatabase, remoteCollection string, uniqueKey []string, data interface{}) {
	if nil != instance.syncHandler {
		message := &DBSyncMessage{
			UID:        instance.UID,
			Driver:     driver,
			Database:   remoteDatabase,
			Collection: remoteCollection,
			UniqueKey:  uniqueKey,
			Data:       data,
		}
		instance.syncHandler(message)
	}
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
