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
	UID    string
	Config *DBSyncConfig

	//-- private --//
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
	instance.Config = config
	instance.tickers = make([]*DBSync, 0)
	instance.events = lygo_events.NewEmitter()
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncSlave) Open() error {
	if nil != instance {
		var err error
		// init client
		instance.init()
		// start client
		err = instance.client.Open()
		// start tickers also if server is offline
		instance.startTickers()
		return err
	}
	return nil
}

// downloads data from server and update slave with server data
// Return array of int64 and array of errors. One item for each collection
func (instance *DBSyncSlave) Reverse() ([]int64, []error) {
	if nil != instance {
		return instance.reverseSync()
	}
	return nil, nil
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
	if nil != instance {
		if len(instance.Config.Uuid) > 0 {
			instance.UID = instance.Config.Uuid
		} else {
			instance.UID, _ = lygo_sys.ID()
		}
		if nil == instance.client {
			instance.client = lygo_nio.NewNioClient(instance.Config.Host(), instance.Config.Port())
			instance.client.OnConnect(instance.doConnect)
			instance.client.OnDisconnect(instance.doDisconnect)
		}
	}
}

func (instance *DBSyncSlave) startTickers() {
	items := instance.Config.Sync
	for _, config := range items {
		ticker := NewDBSync(instance.UID, instance.client, instance.Config.Database, config, &instance.mux)
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

func (instance *DBSyncSlave) pauseTickers() {
	for _, ticker := range instance.tickers {
		_ = ticker.Pause()
	}
}

func (instance *DBSyncSlave) resumeTickers() {
	for _, ticker := range instance.tickers {
		_ = ticker.Resume()
	}
}

func (instance *DBSyncSlave) reverseSync() ([]int64, []error) {
	// pause all tickers
	instance.pauseTickers()
	defer instance.resumeTickers()

	errs := make([]error, 0)
	totals := make([]int64, 0)
	// ready for slave update
	for _, ticker := range instance.tickers {
		count, err := ticker.ReverseSync()
		if nil != err {
			errs = append(errs, err)
			totals = append(totals, 0)
		} else {
			errs = append(errs, nil)
			totals = append(totals, count)
		}
	}
	return totals, errs
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

func (instance *DBSyncSlave) onActionSync(message *DBSyncMessage) []map[string]interface{} {
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
						return []map[string]interface{}{entity}
					}
				} else if strings.Index(s, "[") == 0 {
					var entities []map[string]interface{}
					err := lygo_json.Read(s, &entities)
					if nil == err {
						return entities
					}
				}
			}
		}
	}
	return nil // rollback transaction
}
