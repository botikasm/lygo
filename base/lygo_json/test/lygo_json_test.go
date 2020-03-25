package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Item struct {
	Name string `json:"name"`
}

func TestStruct(t *testing.T) {
	var item Item
	err := lygo_json.ReadFromFile("./item.json", &item)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	name := item.Name
	assert.EqualValues(t, "Angelo", name, "Unexpected value")

	var a []Item
	err = lygo_json.ReadFromFile("./array.json", &a)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	name = a[0].Name
	assert.EqualValues(t, "Angelo", name, "Unexpected value")

	s := lygo_json.Stringify(a)
	fmt.Println(s)
}

func TestMap(t *testing.T) {
	m, err := lygo_json.ReadMapFromFile("./item.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	name := m["name"]
	assert.EqualValues(t, "Angelo", name, "Unexpected value")

	s := lygo_json.Stringify(m)
	fmt.Println(s)
}

func TestArray(t *testing.T) {
	a, err := lygo_json.ReadArrayFromFile("./array.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	name := a[0]["name"]
	assert.EqualValues(t, "Angelo", name, "Unexpected value")

	s := lygo_json.Stringify(a)
	fmt.Println(s)
}
