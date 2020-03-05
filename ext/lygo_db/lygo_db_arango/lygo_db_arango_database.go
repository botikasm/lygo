package lygo_db_arango

import (
	"context"
	"github.com/arangodb/go-driver"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type ArangoDatabase struct {
	database driver.Database
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArangoDatabase) Native() driver.Database {
	if nil!=instance && instance.IsReady() {
		return instance.database
	}
	return nil
}

func (instance *ArangoDatabase) IsReady() bool {
	return nil!=instance && nil != instance.database
}

func (instance *ArangoDatabase) Name() string {
	if nil!=instance && instance.IsReady() {
		return instance.database.Name()
	}
	return ""
}

func (instance *ArangoDatabase) Drop() (bool, error) {
	if nil!=instance && instance.IsReady() {
		ctx := context.Background()
		err := instance.database.Remove(ctx)
		return nil == err, err
	}
	return false, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionNames() ([]string, error) {
	response := make([]string, 0)
	if nil!=instance && instance.IsReady() {
		ctx := context.Background()
		collections, err := instance.database.Collections(ctx)
		if nil != err {
			return response, err
		}
		for _, coll := range collections {
			response = append(response, coll.Name())
		}
	}
	return response, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionExists(name string) (bool, error) {
	if nil!=instance && instance.IsReady() {
		ctx := context.Background()
		exists, err := instance.database.CollectionExists(ctx, name)
		if nil != err {
			return false, err
		}
		return exists, nil
	}
	return false, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) CollectionAutoCreate(name string) (*ArangoCollection, error) {
	return instance.Collection(name, true)
}

func (instance *ArangoDatabase) Collection(name string, createIfNotExists bool) (*ArangoCollection, error) {
	if nil!=instance && instance.IsReady() {
		ctx := context.Background()
		exists, err := instance.database.CollectionExists(ctx, name)
		if nil != err {
			return nil, err
		}

		if !exists && createIfNotExists {
			_, err := instance.database.CreateCollection(ctx, name, nil)
			if nil != err {
				return nil, err
			}
		}

		collection, err := instance.database.Collection(ctx, name)
		if nil != err {
			return nil, err
		}
		if nil == collection {
			return nil, ErrCollectionDoesNotExists
		}
		response := new(ArangoCollection)
		response.collection = collection

		return response, nil
	}
	return nil, ErrDatabaseDoesNotExists
}

func (instance *ArangoDatabase) Query(query string, bindVars map[string]interface{}, callback QueryCallback) (*ArangoCollection, error) {
	if nil!=instance && instance.IsReady() {
		ctx := context.Background()
		cursor, err := instance.database.Query(ctx, query, bindVars)
		if nil != err {
			return nil, err
		}

		defer cursor.Close()
		for {
			var doc interface{}
			meta, err := cursor.ReadDocument(ctx, &doc)
			if driver.IsNoMoreDocuments(err) {
				break
			} else {
				if nil != callback {
					exit := callback(meta, doc, err)
					if exit {
						break
					}
				}
			}
		}
	}
	return nil, ErrDatabaseDoesNotExists
}



//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
