package lygo_n

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_net"
	"time"
)

const (
	nodeLimit = 100
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NDiscovery struct {
	initialized bool
	uuid        string
	config      *lygo_n_commons.NDiscoverySettings
	events      *lygo_events.Emitter
	storage     *NStorage
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNodeDiscovery(events *lygo_events.Emitter, uuid string, config *lygo_n_commons.NDiscoverySettings) *NDiscovery {
	instance := new(NDiscovery)
	instance.config = config
	instance.events = events
	instance.initialized = false
	instance.uuid = uuid

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NDiscovery) String() string {
	if nil != instance {
		return instance.config.String()
	}
	return ""
}

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
			// add defaults
			for _, item := range instance.config.Publishers {
				response = append(response, item.String())
			}
		}
	}
	return response
}

func (instance *NDiscovery) Nodes() []string {
	response := make([]string, 0)
	if nil != instance {
		if nil != instance.storage {
			list := instance.storage.QueryNodes(instance.config.NetworkId)
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
	}
	return response
}

func (instance *NDiscovery) IsNetworkOfNodesEnabled() bool {
	if nil != instance {
		return len(instance.config.NetworkId) > 0
	}
	return false
}

func (instance *NDiscovery) AcquireNode() *lygo_n_net.NConn {
	var response *lygo_n_net.NConn = nil
	if instance.IsNetworkOfNodesEnabled() {
		nodes := instance.Nodes() // get network nodes
		if len(nodes) > 0 {
			for _, node := range nodes {
				na := lygo_n_commons.NewNAddress(node)
				host := na.Host()
				port := na.Port()
				conn := lygo_n_net.NewNConn(host, port)
				errs, _ := conn.Start()
				if len(errs) == 0 {
					if instance.storage.IsLockedNode(node) {
						// add locked node only if no node have still been assigned
						if nil == response {
							response = conn
						}
					} else {
						if nil != response {
							response.Stop()
						}
						response = conn
						break
					}
				} else {
					instance.removeNode(na)
				}
			}
		}
		if nil != response {
			instance.storage.LockNode(response.GetAddress())
		}
	}

	return response
}

func (instance *NDiscovery) ReleaseNode(conn *lygo_n_net.NConn) {
	if nil != conn {
		defer conn.Stop()
		instance.storage.UnlockNode(conn.GetAddress())
	}
}

func (instance *NDiscovery) NewPublisherConnection() *lygo_n_net.NConn {
	nodes := instance.Publishers()
	for _, node := range nodes {
		na := lygo_n_commons.NewNAddress(node)
		host := na.Host()
		port := na.Port()
		conn := lygo_n_net.NewNConn(host, port)
		errs, _ := conn.Start()
		if len(errs) == 0 {
			return conn
		} else {
			instance.removeNode(na)
		}
	}
	return nil
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
			instance.storage.UnlockNodes(instance.config.NetworkId) // release existing nodes

			// add self as a publisher
			if instance.config.Publisher.Enabled && instance.config.Publish.HasAddress() {
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
		networkId := instance.config.NetworkId
		if len(networkId) == 0 {
			break
		}

		// start broadcast
		publishers := instance.Publishers()
		for _, publisher := range publishers {
			// avoid itself
			if !instance.config.Publish.IsAddress(publisher) {
				na := lygo_n_commons.NewNAddress(publisher)
				go instance.syncWithPublisher(networkId, na)
			}
		}

		// interval
		time.Sleep(3 * time.Second)
	}
	instance.events.EmitAsync(lygo_n_commons.EventQuitDiscovery)
}

func (instance *NDiscovery) syncWithPublisher(networkId string, publisherAddress lygo_n_commons.NAddress) {
	host := publisherAddress.Host()
	port := publisherAddress.Port()
	publisherConn := lygo_n_net.NewNConn(host, port)
	errs, _ := publisherConn.Start()
	defer publisherConn.Stop()

	if len(errs) == 0 {
		// publish itself if enabled
		err := instance.registerAsNode(publisherConn, networkId)
		if nil != err {
			// publisher response error
			instance.removePublisher(publisherAddress)
			return
		}

		// synchronize nodes
		err = instance.syncNodes(publisherConn, networkId)
		if nil != err {
			// publisher response error
			instance.removePublisher(publisherAddress)
			return
		}
	} else {
		// publisher offline
		instance.removePublisher(publisherAddress)
	}
}
func (instance *NDiscovery) registerAsNode(publisherConn *lygo_n_net.NConn, networkId string) error {
	if instance.config.Publish.IsEnabled() {
		response := instance.sendCommand(publisherConn, lygo_n_commons.CMD_REGISTER_NODE, map[string]interface{}{
			"address":    instance.config.Publish.Address.String(),
			"network_id": networkId,
		})
		return response.GetError()
	}
	return nil
}

func (instance *NDiscovery) syncNodes(publisherConn *lygo_n_net.NConn, networkId string) error {
	// get list of nodes
	response, err := instance.sendCommandRetMap(publisherConn, lygo_n_commons.CMD_GET_NODE_LIST, map[string]interface{}{
		"network_id": networkId,
	})
	// fmt.Println("discoverAddress", string(response), err)
	if nil == err && len(response) > 0 {
		p := response["publishers"]
		if nil != p {
			instance.onNewPublisher(p.([]interface{}))
		}
		n := response["nodes"]
		if nil != n {
			instance.onNewNode(n.([]interface{}))
		}
	} else {
		if nil != err {
			return err
		}
		return lygo_n_commons.EmptyResponseError
	}
	return nil
}

func (instance *NDiscovery) removePublisher(na lygo_n_commons.NAddress) {
	_ = instance.storage.RemovePublisher(na)
	instance.events.EmitAsync(lygo_n_commons.EventRemovedPublisher, na.String())
}

func (instance *NDiscovery) removeNode(na lygo_n_commons.NAddress) {
	_ = instance.storage.RemoveNode(na)
	instance.events.EmitAsync(lygo_n_commons.EventRemovedNode, na.String())
}

func (instance *NDiscovery) onNewPublisher(items []interface{}) {
	if nil != items && len(items) > 0 {
		for _, item := range items {
			count := instance.storage.CountPublishersAll()
			if count < nodeLimit {
				// add node to database
				key := lygo_reflect.GetString(item, "_key")
				instance.storage.AddPublisher(lygo_n_commons.NewNAddress(key))
				instance.events.EmitAsync(lygo_n_commons.EventNewPublisher, count, key)
			}
		}
	}
}

func (instance *NDiscovery) onNewNode(items []interface{}) {
	if nil != items && len(items) > 0 {
		for _, item := range items {
			key := lygo_reflect.GetString(item, "_key")
			networkId := lygo_reflect.GetString(item, "network_id")
			count := instance.storage.CountNodes(networkId)
			if count < nodeLimit {
				// add node to database
				instance.storage.AddNode(lygo_n_commons.NewNAddress(key), networkId)
				instance.events.EmitAsync(lygo_n_commons.EventNewNode, count, key, networkId)
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m m a n d    s e n d e r
//----------------------------------------------------------------------------------------------------------------------

func (instance *NDiscovery) sendCommand(conn *lygo_n_net.NConn, command string, params map[string]interface{}) *lygo_n_commons.Response {
	return conn.Send(command, params)
}

func (instance *NDiscovery) sendCommandRetMap(conn *lygo_n_net.NConn, command string, params map[string]interface{}) (map[string]interface{}, error) {
	response := conn.Send(command, params)
	if response.HasError() {
		return nil, response.GetError()
	}
	var items map[string]interface{}
	err := lygo_json.Read(response.Data, &items)
	if nil != err {
		return nil, err
	}
	return items, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m m a n d    h a n d l e r s
//----------------------------------------------------------------------------------------------------------------------

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

// CMD_REGISTER_NODE
func (instance *NDiscovery) registerNode(message *lygo_n_commons.Command) interface{} {
	if nil != instance.storage {
		address := message.GetParamAsString("address")
		networkId := message.GetParamAsString("network_id")
		return instance.storage.AddNode(lygo_n_commons.NewNAddress(address), networkId)
	}
	return nil
}
