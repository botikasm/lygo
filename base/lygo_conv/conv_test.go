package lygo_conv

import (
	"fmt"
	"github.com/arangodb/go-velocypack/test"
	"testing"
)


func TestToString(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil", args{nil}, ""},
		{"string", args{"hola"}, "hola"},
		{"int", args{123}, "123"},
		{"float32", args{123.456}, "123.456"},
		{"boolean", args{true}, "true"},
		{"array", args{[]int{1,2,3}}, "[1,2,3]"},
		{"[]byte", args{[]byte("hello")}, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.val); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	var v interface{} = 12345678.456
	fmt.Println(fmt.Sprintf("%v", ToInt(v)))
}

func TestToArray(t *testing.T) {
	arr := ToArray([]string{"1", "2"})
	test.ASSERT_EQ([]interface{}{"1", "2"}, arr, t)

	arr = ToArray([]int{1, 2})
	test.ASSERT_EQ([]interface{}{1, 2}, arr, t)
}

func TestToArrayOfString(t *testing.T) {
	arr := ToArrayOfString([]string{"1", "2"})
	test.ASSERT_EQ([]string{"1", "2"}, arr, t)

	arr = ToArrayOfString([]int{1, 2})
	test.ASSERT_EQ([]string{"1", "2"}, arr, t)

	arr = ToArrayOfString([]interface{}{1, 2})
	test.ASSERT_EQ([]string{"1", "2"}, arr, t)
}