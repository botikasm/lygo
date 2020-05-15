package lygo_strings

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	type args struct {
		s      string
		params []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"text %s %s", []interface{}{"1", 2}}, "text 1 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.s, tt.args.params...); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		slice   []string
		trimVal string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestTrimSpaces(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestCapitalizeFirst(t *testing.T) {
	response := CapitalizeFirst("lower words")
	if response!="Lower words"{
		t.Error("Failed")
		t.FailNow()
	}

	response = CapitalizeFirst("lower")
	if response!="Lower"{
		t.Error("Failed")
		t.FailNow()
	}

	response = CapitalizeFirst("")
	if response!=""{
		t.Error("Failed")
		t.FailNow()
	}
}

func TestFill(t *testing.T) {
	s := FillLeft("123", 10, '0')
	fmt.Println(s)

	s = FillLeft("1234567890123", 10, '-')
	fmt.Println(s)

	s = FillLeft("123456789", 10, '0')
	fmt.Println(s)

	s = FillRight("123", 10, '*')
	fmt.Println(s)
}
