package test

import (
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_db_sync"
	"testing"
)

func TestSimple(t *testing.T) {

	psw, _ := lygo_io.ReadTextFromFile("./psw.txt")

	// master
	mc := new(lygo_db_sync.DBSyncConfig)
	mc.Load("./master.json")
	mc.Database.Authentication.Password = psw
	master := lygo_db_sync.NewDBSyncMaster(mc)
	err := master.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// slave
	sc := new(lygo_db_sync.DBSyncConfig)
	sc.Load("./slave.json")
	sc.Database.Authentication.Password = psw
	slave := lygo_db_sync.NewDBSyncSlave(sc)
	err = slave.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// wait
	slave.Join()
}
