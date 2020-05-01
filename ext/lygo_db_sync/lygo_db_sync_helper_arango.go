package lygo_db_sync

import (
	"github.com/arangodb/go-driver"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_arango"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

// ARANGO DB HELPER
// Use helper to insert/update data setting automatically all sync flags.
type DBSyncHelperArango struct {

	//-- private --//
	config *DBSyncDatabaseConfig
	conn   *lygo_db_arango.ArangoConnection
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewDBSyncHelperArango(config *DBSyncDatabaseConfig) *DBSyncHelperArango {
	instance := new(DBSyncHelperArango)
	instance.config = config

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncHelperArango) Open() error {
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

func (instance *DBSyncHelperArango) Close() error {
	if nil != instance {
		if nil != instance.conn {
			instance.conn = nil
		}
	}
	return nil
}

func (instance *DBSyncHelperArango) Query(database string, query string, params map[string]interface{}) ([]interface{}, error) {
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

func (instance *DBSyncHelperArango) Upsert(database, collection string, item map[string]interface{}) (interface{}, error) {
	if nil != instance {
		item[FLD_FLAG] = true // enable sync
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return nil, err
		}
		coll, err := db.Collection(collection, true)
		if nil != err {
			return nil, err
		}
		response, _, err := coll.Upsert(item)
		return response, err
	}
	return nil, nil
}

func (instance *DBSyncHelperArango) CountQuery(database, query string, bindVars map[string]interface{}) (int64, error) {
	if nil != instance {
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return -1, err
		}
		count, err := db.Count(query, bindVars)
		if nil != err {
			return -1, err
		}
		return count, err
	}
	return -1, nil
}

func (instance *DBSyncHelperArango) Count(database, collection string) (int64, error) {
	if nil != instance {
		db, err := instance.conn.Database(database, true)
		if nil != err {
			return -1, err
		}
		coll, err := db.Collection(collection, true)
		if nil != err {
			return -1, err
		}
		return coll.Count()
	}
	return -1, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
