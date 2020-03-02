package lygo_reflect

import (
	"fmt"
	"testing"
)

type MyDoc struct {
	Name string
	Date string
}

func TestSimple(t *testing.T) {

	// works with struct
	doc := &MyDoc{
		Name: "Foo",
		Date: "25/2/2020",
	}
	name := Get(doc, "Name")
	if name != "Foo" {
		t.Fail()
		t.Errorf("Expected 'Foo', but got '%v'", name)
	}
	if _, b := Set(doc, "Name", "Test"); b {
		t.Fail()
	}
	name = Get(doc, "Name")
	if name != "Test" {
		t.Fail()
		t.Errorf("Expected 'Test', but got '%v'", name)
	}

	// works with map
	var mdoc interface{}
	mdoc = map[string]interface{}{
		"Name": "Foo",
	}
	name = Get(mdoc, "Name")
	if name != "Foo" {
		t.Fail()
	}
	_, b := Set(mdoc, "Name", "Test")
	if !b {
		t.Fail()
	}
	name = Get(mdoc, "Name")
	if name != "Test" {
		t.Fail()
		t.Errorf("Expected 'Test', but got '%v'", name)
	}
	fmt.Println("NAME", name)
}
