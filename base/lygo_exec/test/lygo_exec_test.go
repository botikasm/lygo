package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_exec"
	"runtime"
	"testing"
)

func TestCmd(t *testing.T) {
	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		out, err = lygo_exec.Run("tasklist")
	} else {
		out, err = lygo_exec.Run("ls", "-lah")
	}
	if nil != err {
		t.Errorf("Error: %v", err)
		t.FailNow()
	}
	fmt.Println(string(out))
}

func TestDocx(t *testing.T) {
	err := lygo_exec.Open("./simple.docx")
	if nil != err {
		t.Errorf("Error: %v", err)
		t.FailNow()
	}
}
