package lygo_csv

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {

	in := `first_name;last_name;username
"Rob";"Pike";rob
# lines beginning with a # character are ignored
Ken;Thompson;ken
"Robert";"Griesemer";"gri"`

	options := NewCsvOptionsDefaults()

	data, err := ReadAll(in, options)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(data)

	// missing a column in data
	in = `first_name;last_name;username
"Rob";"Pike";rob
# lines beginning with a # character are ignored
Ken;Thompson;ken
"Robert";"Griesemer"`
	data, err = ReadAll(in, options)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(data)
}

func TestHeaders(t *testing.T) {

	in := `first_name;last_name;username
"Rob";"Pike";rob
# lines beginning with a # character are ignored
Ken;Thompson;ken`

	options := NewCsvOptionsDefaults()
	options.FirstRowHeader = true

	data, err := ReadAll(in, options)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println("rows", len(data), data)

	options.FirstRowHeader = false
	data, err = ReadAll(in, options)
	if nil != err {
		t.Error(err)
		t.Fail()
	}
	fmt.Println("rows", len(data), data)
}
