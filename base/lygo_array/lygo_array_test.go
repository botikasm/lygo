package lygo_array

import (
	"fmt"
	"testing"
)

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
