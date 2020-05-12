package lygo_array

import (
	"fmt"
	"testing"
	"time"
)

type MyDoc struct {
	Name string
	Date string
}

func (instance *MyDoc) String() string {
	return instance.Name
}

func TestShuffle(t *testing.T) {
	array := []string{"1", "2", "3", "4", "5"}
	Shuffle(array)
	fmt.Println(array)
}

func TestSub(t *testing.T) {
	array := []string{"1", "2", "3", "4", "5"}
	n := Sub(array, 1, 1)
	fmt.Println(array, n)
}

func TestIndexOf(t *testing.T) {
	array := []string{"1", "2", "3", "4", "5"}
	i := IndexOf("2", array)
	fmt.Println(array, i)
	i = IndexOf("2", &array)
	fmt.Println(array, i)
}

func TestSort(t *testing.T) {
	array := []string{"7", "2", "4", "5", "3", "1"}
	Sort(array)
	fmt.Println("SORTED ARRAY", array)
	Reverse(array)
	fmt.Println("REVERTED ARRAY", array)

	intArr := make([]int, 0)
	for i := 0; i < 1000000; i++ {
		intArr = append(intArr, i)
	}
	// reverse
	now := time.Now()
	Reverse(intArr)
	diff := time.Now().Sub(now)
	fmt.Println("Reverse elapsed", diff.Milliseconds(), intArr[0:20])

	// sort
	now = time.Now()
	Sort(intArr)
	diff = time.Now().Sub(now)
	fmt.Println("Sort elapsed", diff.Milliseconds(), intArr[0:20])

	// shuffle
	now = time.Now()
	Shuffle(intArr)
	diff = time.Now().Sub(now)
	fmt.Println("Shuffle elapsed", diff.Milliseconds(), intArr[0:20])

	// sort desc
	now = time.Now()
	SortDesc(intArr)
	diff = time.Now().Sub(now)
	fmt.Println("Sort desc", diff.Milliseconds(), intArr[0:20])
}

func TestAppendUnique(t *testing.T) {
	array1 := []interface{}{"1", "2", "3", "4", "5", 1}
	array2 := []string{"1", "4", "6", "7"}
	arrayX := AppendUnique(&array1, array2).([]interface{})
	if len(arrayX) != 8 {
		t.Error("Invalid number of items")
		t.FailNow()
	}
	fmt.Println("NEW ARRAY", arrayX)

	arrayX = AppendUnique(arrayX, "hello").([]interface{})
	arrayX = AppendUnique(arrayX, "1").([]interface{})
	fmt.Println("NEW ARRAY", arrayX)
}

func TestAppendUniqueFunc(t *testing.T) {
	doc1 := &MyDoc{
		Name: "Mario",
		Date: "now",
	}
	doc2 := &MyDoc{
		Name: "Mario",
		Date: "now",
	}
	doc3 := &MyDoc{
		Name: "Ivan",
		Date: "now",
	}
	slice1 := []*MyDoc{doc1}
	slice2 := []*MyDoc{doc1, doc2, doc3}

	slice1 = AppendUniqueFunc(slice1, slice2, func(t, s interface{}) bool {
		n1 := t.(*MyDoc).Name
		n2 := s.(*MyDoc).Name
		return n1 != n2
	}).([]*MyDoc)
	fmt.Println(len(slice1), slice1)

	doc4 := &MyDoc{
		Name: "Podovsky",
		Date: "now",
	}
	slice1 = AppendUniqueFunc(slice1, doc4, func(t, s interface{}) bool {
		n1 := t.(*MyDoc).Name
		n2 := s.(*MyDoc).Name
		return n1 != n2
	}).([]*MyDoc)
	fmt.Println(len(slice1), slice1)
}
