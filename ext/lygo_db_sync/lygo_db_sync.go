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

type DBSyncDriver interface {
	Open() error
	Close() error
	BuildQuery(collection string, filter string, params map[string]interface{}) string
	Collection(database, collection string) (bool, error)
	Execute(database string, query string, params map[string]interface{}) ([]interface{}, error)
}

type DBSync struct {
	UID string

	//-- private --//
	dbConfig     *DBSyncDatabaseConfig
	config       *DBSyncConfigSync
	ticker       *lygo_events.EventTicker
	errorHandler func(sender *DBSync, err error)
	syncHandler  func(sender *DBSync, driver, remoteDatabase, remoteCollection string, uniqueKey []string, data interface{})
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSync(uid string, dbConfig *DBSyncDatabaseConfig, config *DBSyncConfigSync) *DBSync {
	instance := new(DBSync)
	instance.UID = uid
	instance.dbConfig = dbConfig
	instance.config = config
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
func (instance *DBSync) OnSync(callback func(sender *DBSync, driver, remoteDatabase, remoteCollection string, uniqueKey []string, data interface{})) {
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
		instance.syncHandler(instance, driver, remoteDatabase, remoteCollection, uniqueKey, data)
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
