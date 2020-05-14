package lygo_db_sync

import (
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_arango"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	TplQueryLocal = "FOR doc IN @@collection " +
		"FILTER (doc.$FLD_FLAG==NULL || doc.$FLD_FLAG==true) $FILTER" +
		// "LIMIT 1000 " +
		"UPDATE doc WITH { $FLD_FLAG:false, $FLD_TIMESTAMP:@timestamp, $FLD_UUID:@uuid } " +
		"IN @@collection OPTIONS { ignoreErrors: true } " +
		"RETURN NEW"
	TplQueryReverse = "FOR doc IN @@collection " +
		"FILTER doc.$FLD_UUID==@uuid $FILTER" +
		"LIMIT @skip, @limit " +
		"RETURN doc"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncDriverArango struct {

	//-- private --//
	config *DBSyncDatabaseConfig
	conn   *lygo_db_arango.ArangoConnection
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncDriverArango(config *DBSyncDatabaseConfig) *DBSyncDriverArango {
	instance := new(DBSyncDriverArango)
	instance.config = config

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncDriverArango) Open() error {
	if nil != instance {
		if nil != instance.config {
			cfg := lygo_db_arango.NewArangoConfig()
			cfg.Parse(instance.config.String())
			conn := lygo_db_arango.NewArangoConnection(cfg)
			err := conn.Open()
			if nil != err {
				return err
			}
			instance.conn = conn
		} else {
			return MissingDatabaseConfigurationError
		}
	}
	return nil
}

func (instance *DBSyncDriverArango) Close() error {
	if nil != instance {
		if nil != instance.conn {
			instance.conn = nil
		}
	}
	return nil
}

func (instance *DBSyncDriverArango) Collection(database, collection string) (bool, error) {
	if nil != instance {
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return false, err
		}
		_, err = db.CollectionAutoCreate(collection)
		return nil == err, err
	}
	return false, nil
}

// FOR doc IN @@collection FILTER (doc._db_sync_last==NULL) UPDATE doc WITH { "_db_sync_last":@timestamp } IN @@collection RETURN NEW
func (instance *DBSyncDriverArango) BuildQuery(collection string, filter string, params map[string]interface{}) string {
	response := strings.Replace(TplQueryLocal, "@@collection", collection, -1)
	if len(filter) > 0 {
		response = strings.Replace(response, "$FILTER", "&& ("+filter+") ", -1)
	} else {
		response = strings.Replace(response, "$FILTER", filter, -1)
	}
	response = strings.Replace(response, "$FLD_FLAG", FLD_FLAG, -1)
	response = strings.Replace(response, "$FLD_TIMESTAMP", FLD_TIMESTAMP, -1)
	response = strings.Replace(response, "$FLD_UUID", FLD_UUID, -1)
	return response
}

func (instance *DBSyncDriverArango) BuildQueryReverse(collection string, filter string, params map[string]interface{}) string {
	response := strings.Replace(TplQueryReverse, "@@collection", collection, -1)
	if len(filter) > 0 {
		response = strings.Replace(response, "$FILTER", "&& ("+filter+") ", -1)
	} else {
		response = strings.Replace(response, "$FILTER", filter, -1)
	}
	response = strings.Replace(response, "$FLD_FLAG", FLD_FLAG, -1)
	response = strings.Replace(response, "$FLD_TIMESTAMP", FLD_TIMESTAMP, -1)
	response = strings.Replace(response, "$FLD_UUID", FLD_UUID, -1)
	response = strings.Replace(response, "@@collection", collection, -1)
	for k, v := range params {
		response = strings.Replace(response, "@"+k, lygo_conv.ToStringQuoted(v), -1)
	}

	return response
}

func (instance *DBSyncDriverArango) Execute(database string, query string, params map[string]interface{}) ([]interface{}, error) {
	if nil != instance {
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return nil, err
		}

		response := make([]interface{}, 0)
		err = db.Query(query, params, func(meta driver.DocumentMeta, item interface{}, err error) bool {
			if nil != err {
				return true // exit
			}
			response = append(response, item)
			return false // continue
		})
		return response, err
	}
	return nil, nil
}

func (instance *DBSyncDriverArango) Merge(database, collection string, item map[string]interface{}) (map[string]interface{}, error) {
	if nil != instance {
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return nil, err
		}
		coll, err := db.Collection(collection, true)
		if nil == err {
			key := item["_key"].(string)
			if v, _ := coll.Exists(key); v && len(key) > 0 {
				// get existing
				existing, _, err := coll.Read(key)
				if nil != err {
					return nil, err
				}
				// merge existing (avoid override fields with zero values)
				for k, v := range item {
					if k != "_key" {
						if lygo_reflect.IsZero(existing[k]) {
							existing[k] = v
						} else if !lygo_reflect.IsZero(v) {
							existing[k] = v
						}
					}
				}
				// UPDATE
				response, _, err := coll.Update(existing)
				return response, err
			} else {
				// INSERT
				response, _, err := coll.Insert(item)
				return response, err
			}
		}
		return nil, err
	}
	return nil, nil
}

func (instance *DBSyncDriverArango) SetNeedUpdateFlag(database, collection string, raw_entity interface{}) error {
	if nil != instance {
		entity := lygo_conv.ForceMap(raw_entity)
		if nil != entity {
			entity[FLD_FLAG] = true
			db, err := instance.conn.Database(database, true)
			if nil != err {
				return err
			}
			coll, err := db.Collection(collection, true)
			if nil != err {
				return err
			}
			if nil != coll {
				_, _, err = coll.Upsert(entity)
				return err
			}
		} else {
			return errors.New("entity_not_valid")
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
