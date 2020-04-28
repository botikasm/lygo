package lygo_db_sync

import (
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/base/lygo_sys"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncSlave struct {
	UID string
	//-- private --//
	config  *DBSyncConfig
	client  *lygo_nio.NioClient
	tickers []*DBSync
	mux     sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncSlave(config *DBSyncConfig) *DBSyncSlave {
	instance := new(DBSyncSlave)
	instance.config = config
	instance.tickers = make([]*DBSync, 0)
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
		instance.startTickers()
	}
	return nil
}

func (instance *DBSyncSlave) Close() error {
	if nil != instance {
		instance.stopTickers()
		err := instance.client.Close()
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *DBSyncSlave) Join() {
	if nil != instance {
		instance.client.Join()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncSlave) init() {
	instance.UID, _ = lygo_sys.ID()
	instance.client = lygo_nio.NewNioClient(instance.config.Host(), instance.config.Port())

}

func (instance *DBSyncSlave) startTickers() {
	items := instance.config.Sync
	for _, config := range items {
		ticker := NewDBSync(instance.UID, instance.config.Database, config)
		ticker.OnError(instance.onTickerError)
		ticker.OnSync(instance.onTickerSync)
		instance.tickers = append(instance.tickers, ticker)
		_ = ticker.Open()
	}
}

func (instance *DBSyncSlave) stopTickers() {
	for _, ticker := range instance.tickers {
		_ = ticker.Close()
	}
}

func (instance *DBSyncSlave) onTickerError(sender *DBSync, err error) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		// fmt.Println(err)
		lygo_logs.Error(err)
	}
}

func (instance *DBSyncSlave) onTickerSync(sender *DBSync, driver, remoteDatabase, remoteCollection string, uniqueKey []string, data interface{}) {

	// TODO: HANDLE SYNC
	lygo_logs.Info(driver, remoteDatabase)

}
