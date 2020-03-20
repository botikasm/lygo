package lygo_db_bolt

import (
	"encoding/json"
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
	return nil
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
	return response, nil
}

func (instance *BoltCollection) Get(key string) (interface{}, error) {
	var response interface{}
	if nil != instance && nil != instance.db {
		err := instance.db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(instance.name))
			if nil != b {
				buf := b.Get([]byte(key))
				err := json.Unmarshal(buf, &response)

				return err
			} else {
				return ErrCollectionDoesNotExists
			}
		})
		return response, err
	}
	return nil, nil
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
	return nil
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
	return response, nil
}

func (instance *BoltCollection) Find(query *BoltQuery) ([]interface{}, error) {
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
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltCollection) getKey(entity interface{}) []byte {
	return []byte(lygo_reflect.GetString(entity, "_key"))
}
