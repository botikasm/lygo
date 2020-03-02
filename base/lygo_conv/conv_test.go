package lygo_conv

import "testing"

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