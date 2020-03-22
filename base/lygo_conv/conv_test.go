package lygo_conv

import (
	"github.com/stretchr/testify/assert"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.val); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToArray(t *testing.T) {
	arr := ToArray([]string{"1", "2"})
	assert.EqualValues(t, []interface{}{"1", "2"}, arr, "Unexpected value")

	arr = ToArray([]int{1, 2})
	assert.EqualValues(t, []interface{}{1, 2}, arr, "Unexpected value")
}

func TestToArrayOfString(t *testing.T) {
	arr := ToArrayOfString([]string{"1", "2"})
	assert.EqualValues(t, []string{"1", "2"}, arr, "Unexpected value")

	arr = ToArrayOfString([]int{1, 2})
	assert.EqualValues(t, []string{"1", "2"}, arr, "Unexpected value")

	arr = ToArrayOfString([]interface{}{1, 2})
	assert.EqualValues(t, []string{"1", "2"}, arr, "Unexpected value")
}