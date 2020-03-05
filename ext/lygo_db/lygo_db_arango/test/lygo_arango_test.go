package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_arango"
	"github.com/arangodb/go-driver"
	"testing"
	"time"
)



func TestSimple(t *testing.T) {

	ctext, _ := lygo_io.ReadTextFromFile("config.json")
	config := lygo_db_arango.NewArangoConfig()
	config.Parse(ctext)

	conn := lygo_db_arango.NewArangoConnection(config)
	err := conn.Open()
	if nil != err {
		// fmt.Println(err)
		t.Error(err, ctext)
		t.Fail()
		return
	}
	// print version
	fmt.Println("ARANGO SERVER", conn.Server)
	fmt.Println("ARANGO VERSION", conn.Version)
	fmt.Println("ARANGO LICENSE", conn.License)

	// remove
	conn.DropDatabase("test_sample")

	// create a db
	db, err := conn.Database("test_sample", true)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}

	if nil != db {
		fmt.Println(db.Name())
	}

	coll, err := db.Collection("not_exists", true)
	if nil != err {
		t.Error(err)
	}
	if nil == coll {
		t.Fail()
	}

	// entity
	entity := map[string]interface{}{
		"_key":     "258647",
		"name":    "Angelo",
		"surname": "Geminiani",
	}

	Key := lygo_reflect.GetString(entity, "Key")
	fmt.Println("KEY", Key)

	lygo_reflect.Set(entity, "Name", "Gian Angelo")

	doc, meta, err := coll.Upsert(entity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	entity = map[string]interface{}{
		"name":    "Marco",
		"surname": lygo_strings.Format("%s", time.Now()),
	}
	doc, meta, err = coll.Upsert(entity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)

	if true {
		return
	}

	// bew entity that test upsert used for insert
	newEntity := map[string]interface{}{
		"_key" : lygo_rnd.UuidDefault(""),
		"name":    "I'm new",
		"surname": lygo_strings.Format("%s", time.Now()),
	}
	doc, meta, err = coll.Upsert(newEntity)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println("META", meta)
	fmt.Println("DOC", doc)


	// remove
	removed, err := conn.DropDatabase("test_sample")
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
	if removed {
		fmt.Println("REMOVED", "test_sample")
	}
}

func TestInsert(t *testing.T) {
	ctext, _ := lygo_io.ReadTextFromFile("config.json")
	config := lygo_db_arango.NewArangoConfig()
	config.Parse(ctext)

	conn := lygo_db_arango.NewArangoConnection(config)
	err := conn.Open()
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
		t.Fail()
		return
	}

	// remove
	conn.DropDatabase("test_sample")

	db, err := conn.Database("test_sample", true)
	if nil != err {
		// fmt.Println(err)
		t.Error(err)
	}
	coll, err := db.Collection("coll_insert", true)
	if nil != err {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		// entity
		entity := map[string]interface{}{
			"_key" : lygo_strings.Format("key_%s", i),
			"name":    lygo_strings.Format("Name:%s", i),
			"surname": lygo_strings.Format("Surname:%s", i),
			"address": lygo_strings.Format("Address:%s", i),
		}

		doc, meta, err := coll.Insert(entity)
		if nil != err {
			t.Error(err)
		}
		fmt.Println("META", meta)
		fmt.Println("DOC", doc)
	}

	updEntity := map[string]interface{}{
		"_key" : "key_1",
		"name":    "Gian Angelo",
	}
	_, _, err = coll.Update(updEntity)
	if nil != err {
		t.Error(err)
	}

	noKeyEntity := map[string]interface{}{
		"name":    "NO KEY",
		"surname": "ZERO ZERO",
		"address": "",
	}
	doc, meta, err := coll.Insert(noKeyEntity)
	if nil != err {
		t.Error(err)
	}
	fmt.Println("---------------------")
	fmt.Println("META", meta)
	fmt.Println("DOC KEY", doc["_key"])
	fmt.Println("DOC", lygo_conv.ToString(doc))
	fmt.Println("---------------------")

	query := "FOR d IN coll_insert RETURN d"
	db.Query(query, nil, gotDocument)

	fmt.Println("---------------------")

}

func gotDocument(meta driver.DocumentMeta, doc interface{}, err error) bool{
	fmt.Println(meta, lygo_conv.ToString(doc), err)
	return false// continue
}