package lygo_db_sync

import "github.com/botikasm/lygo/base/lygo_nio"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncMaster struct {

	//-- private --//
	config *DBSyncConfig
	server *lygo_nio.NioServer
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

	return nil
}
