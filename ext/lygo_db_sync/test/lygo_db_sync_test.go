package test

import (
	"github.com/botikasm/lygo/ext/lygo_db_sync"
	"testing"
)

func TestSimple(t *testing.T) {

	// master
	mc := new (lygo_db_sync.DBSyncConfig)
	mc.Load("./master.json")
	master := lygo_db_sync.NewDBSyncMaster(mc)
	err := master.Open()
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}

	// slave
	sc := new (lygo_db_sync.DBSyncConfig)
	sc.Load("./slave.json")
	slave := lygo_db_sync.NewDBSyncSlave(sc)
	err = slave.Open()
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}

	// wait
	slave.Join()
}