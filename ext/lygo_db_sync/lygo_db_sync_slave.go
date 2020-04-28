package lygo_db_sync

import "github.com/botikasm/lygo/base/lygo_nio"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncSlave struct {

	//-- private --//
	config *DBSyncConfig
	client *lygo_nio.NioClient
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncSlave(config *DBSyncConfig) *DBSyncSlave {
	instance := new(DBSyncSlave)
	instance.config = config
	instance.init()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncSlave) Open() error {
	if nil != instance {
		err := instance.client.Open()
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *DBSyncSlave) Join(){
	if nil!=instance{
		instance.client.Join()
	}
}


//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncSlave) init() {
	instance.client = lygo_nio.NewNioClient(instance.config.Host(), instance.config.Port())

}
