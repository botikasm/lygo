package test

import (
	"encoding/json"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_mu"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/base/lygo_stopwatch"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_db/lygo_db_bolt"
	"testing"
)

type Entity struct {
	Key string `json:"_key"`
	Age int    `json:"age"`
}

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

func TestQuery(t *testing.T) {
	item := &map[string]interface{}{
		"_key": "1",
		"name": "Mario",
		"age":  22,
	}

	query, err := lygo_db_bolt.NewQueryFromFile("./query.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	match := query.MatchFilter(item)
	if !match {
		t.Error("Query do not match")
		t.FailNow()
	}

}

func TestBigData(t *testing.T) {

	drop_on_exit := true

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

	coll, err := db.Collection("big-coll", true)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	arrayData := generateArray()

	watch := lygo_stopwatch.New()
	watch.Start()

	//-- START LOOP TO ADD RECORDS --//
	cur, _ := coll.Count()
	cur++
	for i := cur; i < cur+10; i++ {
		item := &map[string]interface{}{
			"_key": lygo_strings.Format("%s", i),
			"name": "NAME " + lygo_strings.Format("%s", i),
			"age":  i,
			"x":    arrayData,
			"y":    arrayData,
		}
		err = coll.Upsert(item)
		if nil != err {
			t.Error(err)
			t.Fail()
			break
		}
	}
	watch.Stop()
	fmt.Println("ELAPSED FOR CREATION: ", watch.Seconds(), "seconds")

	watch.Start()
	count, err := coll.Count()
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	watch.Stop()
	fmt.Println("RECORDS: ", count)
	fmt.Println("ELAPSED FOR COUNT: ", watch.Seconds(), "seconds")

	watch.Start()
	query, err := lygo_db_bolt.NewQueryFromFile("./query.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	query.Conditions[0].Filters[0].Value = 1
	//var entity Entity
	//data, err := coll.Find(query, &entity)
	// data, err := coll.Find(query)
	data := make([]map[string]interface{}, 0)
	err = coll.ForEach(func(k, v []byte) bool {
		key := string(k)
		fmt.Println("KEY", key)
		var e Entity
		json.Unmarshal(v, &e)
		if query.MatchFilter(e) {
			var m map[string]interface{}
			json.Unmarshal(v, &m)
			data = append(data, m)
		}
		return false // continue
	})

	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	watch.Stop()
	fmt.Println("QUERY: ", len(data))
	fmt.Println("ELAPSED FOR QUERY: ", watch.Seconds(), "seconds")
	if len(data) > 0 {
		for _, item := range data {
			fmt.Println("AGE:", lygo_reflect.Get(item, "age"))
		}
	}

	size, _ := db.Size()
	fmt.Println("FILE SIZE: ", lygo_mu.ToMegaBytes(size), "Mb")

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

func generateArray() []float64 {
	size := 250 * 60 * 10 // 10 minutes of 250Hz data sequence
	data := make([]float64, size)
	for i := 0; i < size; i++ {
		val := float64(i * 2)
		data[i] = val
	}
	return data
}
