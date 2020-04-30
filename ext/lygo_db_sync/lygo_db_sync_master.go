package lygo_db_sync

import (
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncMaster struct {

	//-- private --//
	config *DBSyncConfig
	server *lygo_nio.NioServer
	mux    sync.Mutex // global mutex to synchronize updates
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncMaster(config *DBSyncConfig) *DBSyncMaster {
	instance := new(DBSyncMaster)
	instance.config = config
	instance.init()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncMaster) Open() error {
	if nil != instance {
		err := instance.server.Open()
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *DBSyncMaster) Close() error {
	if nil != instance {
		err := instance.server.Close()
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *DBSyncMaster) Join() {
	if nil != instance {
		instance.server.Join()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncMaster) init() {
	instance.server = lygo_nio.NewNioServer(instance.config.Port())
	instance.server.OnMessage(instance.onMessage)
}

func (instance *DBSyncMaster) onMessage(message *lygo_nio.NioMessage) interface{} {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if v, b := message.Body.([]byte); b {
		var body DBSyncMessage
		err := lygo_json.Read(v, &body)
		if nil == err {
			driver := GetDriver(body.Driver, instance.config.Database)
			if nil != driver {
				err := driver.Open()
				if nil==err{
					uid := body.UID
					database := body.Database
					collection := body.Collection
					item := body.Data
					if nil != item && len(uid) > 0 && len(database) > 0 && len(collection) > 0 {
						if entity,b:=item.(map[string]interface{});b{
							entity, err = driver.Merge(database, collection, entity)
							if nil==err{
								return entity
							}
						}
					}
				}
			}
		}
	}
	return "false" // custom response
}
