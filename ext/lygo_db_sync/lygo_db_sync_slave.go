package lygo_db_sync

import (
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/base/lygo_sys"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"strings"
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
	events  *lygo_events.Emitter
	mux     sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncSlave(config *DBSyncConfig) *DBSyncSlave {
	instance := new(DBSyncSlave)
	instance.config = config
	instance.tickers = make([]*DBSync, 0)
	instance.events = lygo_events.NewEmitter()

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

func (instance *DBSyncSlave) OnError(callback func(e *lygo_events.Event)) {
	if nil != instance {
		instance.events.On("error", callback)
	}
}

func (instance *DBSyncSlave) OffError() {
	if nil != instance {
		instance.events.Off("error")
	}
}

func (instance *DBSyncSlave) OnConnect(callback func(e *lygo_events.Event)) {
	if nil != instance {
		instance.events.On("connect", callback)
	}
}

func (instance *DBSyncSlave) OffConnect() {
	if nil != instance {
		instance.events.Off("connect")
	}
}

func (instance *DBSyncSlave) OnDisconnect(callback func(e *lygo_events.Event)) {
	if nil != instance {
		instance.events.On("disconnect", callback)
	}
}

func (instance *DBSyncSlave) OffDisconnect() {
	if nil != instance {
		instance.events.Off("disconnect")
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncSlave) init() {
	if len(instance.config.Uuid) > 0 {
		instance.UID = instance.config.Uuid
	} else {
		instance.UID, _ = lygo_sys.ID()
	}
	instance.client = lygo_nio.NewNioClient(instance.config.Host(), instance.config.Port())
	instance.client.OnConnect(instance.doConnect)
	instance.client.OnDisconnect(instance.doDisconnect)
}

func (instance *DBSyncSlave) startTickers() {
	items := instance.config.Sync
	for _, config := range items {
		ticker := NewDBSync(instance.UID, instance.client, instance.config.Database, config)
		ticker.OnError(instance.onActionSyncError)
		ticker.OnSync(instance.onActionSync)
		instance.tickers = append(instance.tickers, ticker)
		_ = ticker.Open()
	}
}

func (instance *DBSyncSlave) stopTickers() {
	for _, ticker := range instance.tickers {
		_ = ticker.Close()
	}
}

func (instance *DBSyncSlave) doError(err error) {
	if nil != instance {
		instance.events.EmitAsync("error", err)
	}
}

func (instance *DBSyncSlave) doConnect(e *lygo_events.Event) {
	if nil != instance {
		instance.events.EmitAsync(e.Name)
	}
}

func (instance *DBSyncSlave) doDisconnect(e *lygo_events.Event) {
	if nil != instance {
		instance.events.EmitAsync(e.Name)
	}
}

func (instance *DBSyncSlave) onActionSyncError(sender *DBSync, err error) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		// fmt.Println(err)
		lygo_logs.Error(err)
		instance.doError(err)
	}
}

func (instance *DBSyncSlave) onActionSync(message *DBSyncMessage) map[string]interface{} {
	if nil != instance {
		if instance.client.IsOpen() {
			response, err := instance.client.Send(message)
			if nil != err {
				lygo_logs.Error(err)
			}
			if nil != response {
				s := string(response.Body.([]byte))
				if strings.Index(s, "{") == 0 {
					var entity map[string]interface{}
					err := lygo_json.Read(s, &entity)
					if nil == err {
						return entity
					}
				}
			}
		}
	}
	return nil // rollback transaction
}
