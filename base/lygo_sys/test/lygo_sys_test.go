package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_sys"
	"testing"
)

func TestShutdown(t *testing.T) {
	err := lygo_sys.Shutdown("1234567890")
	if nil != err {
		t.Error(err)
		t.Fail()
	}
}

func TestGetOS(t *testing.T) {
	fmt.Println("GOOS: ", lygo_sys.GetOS())
}
