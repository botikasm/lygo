package lygo_db_bolt

import (
	"encoding/json"
	"errors"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"go.etcd.io/bbolt"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltCollection struct {

	//-- private --//
	name string
	db   *bbolt.DB
}

type ForEachCallback func(k, v []byte) bool

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltCollection) Drop() error {
	if nil != instance && nil != instance.db {
		err := instance.db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				return tx.DeleteBucket([]byte(instance.name))
			}
			return nil
		})
		return err
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Count() (int64, error) {
	var response int64
	response = 0
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					response++
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) CountByFieldValue(fieldName string, fieldValue interface{}) (int64, error) {
	var response int64
	response = 0
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						value := lygo_reflect.Get(entity, fieldName)
						if lygo_conv.Equals(fieldValue, value) {
							response++
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Get(key string) (interface{}, error) {
	var response interface{}
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				buf := b.Get([]byte(key))
				if nil!=buf{
					err := json.Unmarshal(buf, &response)
					return err
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return nil, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Upsert(entity interface{}) error {
	if nil != instance && nil != instance.db {
		err := instance.db.Update(func(tx *bbolt.Tx) error {
			key := instance.getKey(entity)
			if len(key) == 0 {
				return ErrMissingDocumentKey
			}
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				buf, err := json.Marshal(entity)
				if nil == err {
					b.Put(key, buf)
				}
				return err
			} else {
				return ErrCollectionDoesNotExists
			}
		})
		return err
	}
	return ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) GetByFieldValue(fieldName string, fieldValue interface{}) ([]interface{}, error) {
	response := make([]interface{}, 0)
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						value := lygo_reflect.Get(entity, fieldName)
						if lygo_conv.Equals(fieldValue, value) {
							response = append(response, entity)
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) Find(query *BoltQuery) ([]interface{}, error) {
	response := make([]interface{}, 0)
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var entity map[string]interface{}
					err := json.Unmarshal(v, &entity)
					if nil == err {
						if query.MatchFilter(entity) {
							response = append(response, entity)
						}
					}
				}
			} else {
				return ErrCollectionDoesNotExists
			}
			return nil
		})
		return response, err
	}
	return response, ErrDatabaseIsNotConnected
}

func (instance *BoltCollection) ForEach(callback ForEachCallback) error {
	if nil != instance && nil != instance.db {
		if nil != callback {
			err := instance.db.View(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(instance.name))
				if nil != b {
					_ = b.ForEach(func(k, v []byte) error {
						exit := callback(k, v)
						if exit {
							return errors.New("exit")
						}
						return nil
					})
				} else {
					return ErrCollectionDoesNotExists
				}
				return nil
			})
			return err
		} else {
			return nil
		}
	}
	return ErrDatabaseIsNotConnected
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltCollection) getKey(entity interface{}) []byte {
	return []byte(lygo_reflect.GetString(entity, "_key"))
}

func unmarshal(v []byte, entity interface{}) (interface{}, error) {
	if nil != entity {
		err := json.Unmarshal(v, entity)
		if nil != err {
			return nil, err
		}
		//var resp struct{}
		//lygo_reflect.Copy(reflect.ValueOf(entity).Elem().Interface(), &resp)

		return entity, nil
	} else {
		var resp map[string]interface{}
		err := json.Unmarshal(v, &resp)
		if nil != err {
			return nil, err
		}
		return resp, nil
	}
}
