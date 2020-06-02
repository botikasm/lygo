package lygo_n

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_bolt"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"path"
	"sync"
	"time"
)

const (
	COLL_PUBLISHERS = "publishers"
	COLL_NODES      = "nodes"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NStorage struct {
	config *lygo_db_bolt.BoltConfig
	db     *lygo_db_bolt.BoltDatabase
	mux    sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNodeStorage(workspace string) *NStorage {
	instance := new(NStorage)
	instance.config = lygo_db_bolt.NewBoltConfig()

	instance.config.Name = path.Join(workspace, "storage", "db.dat")
	lygo_paths.Mkdir(instance.config.Name)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NStorage) Start() error {
	if nil != instance {
		instance.db = lygo_db_bolt.NewBoltDatabase(instance.config)
		return instance.db.Open()
	}
	return lygo_n_commons.PanicSystemError
}

func (instance *NStorage) Stop() error {
	if nil != instance {
		if nil != instance.db {
			return instance.db.Close()
		}
	}
	return lygo_n_commons.PanicSystemError
}

//----------------------------------------------------------------------------------------------------------------------
//	q u e r y
//----------------------------------------------------------------------------------------------------------------------

func (instance *NStorage) CountPublishersAll() int64 {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_PUBLISHERS, true)
		if nil == err {
			count, err := coll.Count()
			if nil == err {
				return count
			}
		}
	}
	return 0
}

func (instance *NStorage) QueryPublishersAll() []map[string]interface{} {
	response := make([]map[string]interface{}, 0)
	if nil != instance {
		coll, err := instance.db.Collection(COLL_PUBLISHERS, true)
		if nil == err {
			coll.ForEach(func(k, v []byte) bool {
				var entity map[string]interface{}
				err := json.Unmarshal(v, &entity)
				if nil == err {
					response = append(response, entity)
				}
				return false
			})
		}
	}
	return response
}

func (instance *NStorage) CountNodes(networkId string) int64 {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			count, err := coll.CountByFieldValue("network_id", networkId)
			if nil == err {
				return count
			}
		}
	}
	return 0
}

func (instance *NStorage) QueryNodes(networkId string) []map[string]interface{} {
	response := make([]map[string]interface{}, 0)
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			coll.ForEach(func(k, v []byte) bool {
				var entity map[string]interface{}
				err := json.Unmarshal(v, &entity)
				if nil == err {
					id := lygo_reflect.GetString(entity, "network_id")
					if id == networkId {
						response = append(response, entity)
					}
				}
				return false
			})
		}
	}
	return response
}

func (instance *NStorage) IsLockedNode(key string) bool {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			entity, err := coll.Get(key)
			if nil == err && nil != entity {
				return lygo_reflect.GetBool(entity, "locked")
			}
		}
	}
	return false
}

func (instance *NStorage) LockNode(key string) {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			entity, err := coll.Get(key)
			if nil == err && nil != entity {
				lygo_reflect.Set(entity, "locked", true)
				err = coll.Upsert(entity)
			}
		}
	}
}

func (instance *NStorage) UnlockNode(key string) {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			entity, err := coll.Get(key)
			if nil == err && nil != entity {
				lygo_reflect.Set(entity, "locked", false)
				err = coll.Upsert(entity)
			}
		}
	}
}

func (instance *NStorage) UnlockNodes(networkId string) {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			nodes := instance.QueryNodes(networkId)
			for _, node := range nodes {
				id := lygo_reflect.GetString(node, "network_id")
				if id == networkId {
					lygo_reflect.Set(node, "locked", false)
					err = coll.Upsert(node)
				}
			}
		}
	}
}

func (instance *NStorage) AddPublisher(address lygo_n_commons.NAddress) map[string]interface{} {
	if nil != instance {
		return instance.addItem(COLL_PUBLISHERS, address, "")
	}
	return nil
}

func (instance *NStorage) RemovePublisher(address lygo_n_commons.NAddress) error {
	if nil != instance {
		return instance.removeItem(COLL_PUBLISHERS, address)
	}
	return nil
}

func (instance *NStorage) AddNode(address lygo_n_commons.NAddress, networkId string) map[string]interface{} {
	if nil != instance {
		return instance.addItem(COLL_NODES, address, networkId)
	}
	return nil
}

func (instance *NStorage) RemoveNode(address lygo_n_commons.NAddress) error {
	if nil != instance {
		return instance.removeItem(COLL_NODES, address)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NStorage) init() {

}

func (instance *NStorage) addItem(collName string, address lygo_n_commons.NAddress, networkId string) map[string]interface{} {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		coll, err := instance.db.Collection(collName, true)
		if nil == err {
			item, err := coll.Get(address.String())
			if nil == err {
				if nil != item {
					// update timestamp
					lygo_reflect.Set(item, "timestamp", time.Now().Unix())
				} else {
					// creates new item
					item = map[string]interface{}{
						"_key":       address.String(),
						"timestamp":  time.Now().Unix(), // last ping
						"locked":     false,             // if node is locked from a request
						"network_id": networkId,
					}
				}
				err = coll.Upsert(item)
				if nil == err {
					return item.(map[string]interface{})
				}
			}
		}
	}
	return nil
}

func (instance *NStorage) removeItem(collName string, address lygo_n_commons.NAddress) error {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		coll, err := instance.db.Collection(collName, true)
		if nil == err {
			err = coll.Remove(address.String())
		}
		return err
	}
	return nil
}
