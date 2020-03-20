package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_bolt"
	"testing"
)

func TestDatabase(t *testing.T) {

	drop_on_exit := false

	config, err := getConfig()
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	db := lygo_db_bolt.NewBoltDatabase(config)
	err = db.Open()
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	defer db.Close()

	coll, err := db.Collection("my-coll", false)
	if nil != coll {
		t.Error("COLLECTION SHOULD BE NULL")
		t.Fail()
	}
	fmt.Println("Test OK:", err)

	coll, err = db.Collection("my-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	// insert item
	item := &map[string]interface{}{
		"_key": "1",
		"name": "Mario",
		"age":  22,
	}
	err = coll.Upsert(item)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	item2 := &map[string]interface{}{
		"_key": "2",
		"name": "Giorgio",
		"age":  22,
	}
	err = coll.Upsert(item2)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	item3 := &map[string]interface{}{
		"_key": "3",
		"name": "Mirko",
		"age":  45,
	}
	err = coll.Upsert(item3)
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	count, err := coll.CountByFieldValue("age", 22)
	fmt.Println(count)

	item_des, err := coll.Get("1")
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(item_des)

	data, err := coll.GetByFieldValue("age", 22)
	fmt.Println(data)

	// remove collection
	err = coll.Drop()
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	coll, err = db.Collection("my-coll", false)
	if nil != coll {
		t.Error("COLLECTION SHOULD BE NULL")
		t.Fail()
	}

	// remove database
	if drop_on_exit {
		err = db.Drop()
		if nil != err {
			t.Error(err)
			t.Fail()
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getConfig() (*lygo_db_bolt.BoltConfig, error) {
	text_cfg, err := lygo_io.ReadTextFromFile("./config.json")
	if nil != err {
		return nil, err
	}
	config := lygo_db_bolt.NewBoltConfig()
	err = config.Parse(text_cfg)

	return config, err
}
