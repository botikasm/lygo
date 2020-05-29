package lygo_n

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NDiscovery struct {
	initialized bool
	uuid        string
	config      *NDiscoverySettings
	storage     *NStorage
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNodeDiscovery(uuid string, config *NDiscoverySettings) *NDiscovery {
	instance := new(NDiscovery)
	instance.config = config
	instance.initialized = false
	instance.uuid = uuid

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NDiscovery) IsEnabled() bool {
	if nil != instance {
		return len(instance.config.Publishers) > 0 || instance.config.Publisher.Enabled
	}
	return false
}

func (instance *NDiscovery) Start() error {
	if nil != instance {
		return instance.init()
	}
	return lygo_n_commons.PanicSystemError
}

func (instance *NDiscovery) Stop() error {
	if nil != instance {
		return instance.finish()
	}
	return lygo_n_commons.PanicSystemError
}

func (instance *NDiscovery) Publishers() []string {
	response := make([]string, 0)
	if nil != instance {
		if nil != instance.storage {
			list := instance.storage.QueryPublishersAll()
			if len(list) > 0 {
				for _, item := range list {
					key := lygo_reflect.GetString(item, "_key")
					if len(key) > 0 {
						if lygo_array.IndexOf(key, response) == -1 {
							response = append(response, key)
						}
					}
				}
			}
		}
		if len(response) == 0 {
			for _, item := range instance.config.Publishers {
				response = append(response, item.String())
			}
		}
	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NDiscovery) init() error {
	if !instance.initialized {
		instance.initialized = true

		var err error

		if instance.IsEnabled() {

			// storage
			instance.storage = NewNodeStorage(lygo_paths.GetWorkspacePath())
			err = instance.storage.Start()
			if nil != err {
				goto exit
			}

			// add self as a publisher
			if instance.config.Publisher.Enabled && len(instance.config.Publish.Address) > 0 {
				instance.storage.AddPublisher(instance.config.Publish.Address)
			}

			// main loop
			go instance.discover()
		}

		// exit procedure
	exit:
		return err
	}
	return nil
}

func (instance *NDiscovery) finish() error {
	if instance.initialized && nil != instance.storage {
		instance.initialized = false

		if instance.IsEnabled() {
			var err error

			// storage
			if nil != instance.storage {
				err = instance.storage.Stop()
				instance.storage = nil
			}

			return err
		}
	}
	return nil
}

func (instance *NDiscovery) discover() {
	for {
		if !instance.initialized || !instance.IsEnabled() {
			break
		}

		publishers := instance.Publishers()
		if len(publishers) == 0 {
			break
		}

		networkId := instance.config.NetworkId
		if len(networkId) == 0 {
			break
		}

		// start broadcast
		for _, publisher := range publishers {
			na := NewNAddress(publisher)
			host := na.Host()
			port := na.Port()
			conn := lygo_n_client.NewNClient(host, port)
			errs, _ := conn.Start()
			if len(errs) == 0 {
				// get list of nodes
				response, err := conn.Send(CMD_GET_NODE_LIST, map[string]interface{}{
					"network_id": networkId,
				})
				if nil == err {
					var items map[string]interface{}
					err = lygo_json.Read(string(response), &items)
					if nil == err {
						instance.onNewPublisher(items["publishers"].([]interface{}))
						instance.onNewNode(items["nodes"].([]interface{}))
					}
				}
				_ = conn.Stop()
			}
			time.Sleep(100 * time.Millisecond)
		}

		// interval
		time.Sleep(3 * time.Second)
	}
}

// CMD_GET_NODE_LIST
func (instance *NDiscovery) getNodeList(message *lygo_n_commons.Command) interface{} {
	response := make(map[string]interface{})
	if nil != instance.storage {
		networkId := message.GetParamAsString("network_id")
		response["publishers"] = instance.storage.QueryPublishersAll()
		response["nodes"] = instance.storage.QueryNodes(networkId)
	}
	return response
}

func (instance *NDiscovery) onNewPublisher(items []interface{}) {
	if nil != items && len(items) > 0 {
		for _, item := range items {
			// add node to database
			key := lygo_reflect.GetString(item, "_key")
			instance.storage.AddPublisher(NewNAddress(key))
		}
	}
}

func (instance *NDiscovery) onNewNode(items []interface{}) {
	if nil != items && len(items) > 0 {

	}
}
