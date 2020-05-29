package lygo_n

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_bolt"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"path"
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

func (instance *NStorage) QueryNodes(networkId string) []map[string]interface{} {
	response := make([]map[string]interface{}, 0)
	if nil != instance {
		coll, err := instance.db.Collection(COLL_NODES, true)
		if nil == err {
			coll.ForEach(func(k, v []byte) bool {
				var entity map[string]interface{}
				err := json.Unmarshal(v, &entity)
				if nil != err {
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

func (instance *NStorage) AddPublisher(address NAddress) map[string]interface{} {
	if nil != instance {
		coll, err := instance.db.Collection(COLL_PUBLISHERS, true)
		if nil == err {
			item, err := coll.Get(address.String())
			if nil == err {
				if nil != item {
					// update timestamp
					lygo_reflect.Set(item, "timestamp", time.Now().Unix())
				} else {
					// creates new item
					item = map[string]interface{}{
						"_key":      address.String(),
						"timestamp": time.Now().Unix(),
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

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NStorage) init() {

}
