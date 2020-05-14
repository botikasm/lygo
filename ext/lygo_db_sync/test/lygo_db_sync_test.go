package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/ext/lygo_db_sync"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"testing"
	"time"
)

func TestSyncFromServerToClient(t *testing.T) {
	// init workspace and logging
	lygo_paths.SetWorkspaceParent("./test")
	lygo_logs.SetOutput(lygo_logs.OUTPUT_FILE)

	// master
	master, err := startMaster()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// slave
	slave, err := startSlave()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("Master: ", master.Config.Address)
	fmt.Println("Slave: ", slave.Config.Address)

	counts, errs := slave.Reverse()
	fmt.Println("Collections to update:", len(counts))
	for i, count := range counts {
		err := errs[i]
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println(i, "-", count)
	}
}

func TestSyncFromClientToServerWithMockData(t *testing.T) {

	// init workspace and logging
	lygo_paths.SetWorkspaceParent("./test")
	lygo_logs.SetOutput(lygo_logs.OUTPUT_FILE)

	// master
	master, err := startMaster()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// slave
	slave, err := startSlave()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// client count
	fmt.Println("Slave data:", countClientData(slave.Config))
	// server count
	fmt.Println("Master data:", countServerData(master.Config))

	startAddingEntities(t, slave.Config.Database, "sync_slave", "coll_slave")

	// client count
	fmt.Println("Slave data:", countClientData(slave.Config))
	// server count
	fmt.Println("Master data:", countServerData(master.Config))

	// wait server ends replicate
	fmt.Println("WAITING END SYNC")
	clientCount := countClientData(slave.Config)
	for {
		time.Sleep(5 * time.Second)
		serverCount := countServerData(master.Config)
		if clientCount == serverCount {
			break
		}
	}
	fmt.Println("EXITING...")
	fmt.Println("Master data:", countServerData(master.Config))
	// wait
	//slave.Join()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func startMaster() (*lygo_db_sync.DBSyncMaster, error) {
	psw, _ := lygo_io.ReadTextFromFile("./psw.txt")
	// master
	mc := new(lygo_db_sync.DBSyncConfig)
	mc.Load("./master.json")
	mc.Database.Authentication.Password = psw
	master := lygo_db_sync.NewDBSyncMaster(mc)
	err := master.Open()

	return master, err
}

func startSlave() (*lygo_db_sync.DBSyncSlave, error) {
	psw, _ := lygo_io.ReadTextFromFile("./psw.txt")
	sc := new(lygo_db_sync.DBSyncConfig)
	sc.Load("./slave.json")
	sc.Database.Authentication.Password = psw
	slave := lygo_db_sync.NewDBSyncSlave(sc)
	slave.OnError(handleEvent)
	slave.OnConnect(handleEvent)
	slave.OnDisconnect(handleEvent)
	err := slave.Open()
	return slave, err
}

func countClientData(sc *lygo_db_sync.DBSyncConfig) int64 {
	clientdb := lygo_db_sync.NewDBSyncHelperArango(sc.Database)
	clientdb.Open()
	clientCount, _ := clientdb.Count("sync_slave", "coll_slave")
	clientdb.Close()
	return clientCount
}

func countServerData(mc *lygo_db_sync.DBSyncConfig) int64 {
	serverdb := lygo_db_sync.NewDBSyncHelperArango(mc.Database)
	serverdb.Open()
	serverCount, _ := serverdb.Count("sync_master", "coll_master")
	serverdb.Close()
	return serverCount
}

func handleEvent(e *lygo_events.Event) {
	fmt.Println("handleEvent", e.Name, e.Argument(0))
}

func startAddingEntities(t *testing.T, config *lygo_db_sync.DBSyncDatabaseConfig, database, collection string) {
	arango := lygo_db_sync.NewDBSyncHelperArango(config)
	err := arango.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("START ADDING/UPDATING DATA TO REPLICATE")

	item := map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}
	for i := 0; i < 1000; i++ {
		item["_key"] = fmt.Sprintf("item_%v", i)
		item["now"] = time.Now()
		item["driver"] = config.Driver
		item["username"] = config.Authentication.Username
		item["password"] = config.Authentication.Password
		item["endpoints"] = config.Endpoints
		// add random fields
		random := make(map[string]interface{})
		item["random"] = random
		for ii := 0; ii < 10; ii++ {
			k := lygo_rnd.RndId()
			random[k] = lygo_rnd.RndId()
		}

		_, err := arango.Upsert(database, collection, item)
		if nil != err {
			t.Error(err)
			t.FailNow()
			break
		}
	}

	fmt.Println("END ADDING ENTITIES")

}
