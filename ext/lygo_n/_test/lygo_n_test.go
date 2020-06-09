package _test_test

import (
	"errors"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/botikasm/lygo/ext/lygo_n"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_net"
	"github.com/gofiber/fiber"
	"testing"
	"time"
)

func TestSimpleCommunication(t *testing.T) {

	n := lygo_n.NewNode(config())
	n.Settings.Name = "SINGLE NODE"

	registerCommands(n)

	// http handler
	initializeHttp(n.Http())

	errs := n.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	client := lygo_n_net.NewNConn(configCli())
	errs, _ = client.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	request := map[string]interface{}{
		"namespace": "get",
		"function":  "version",
	}
	response := client.SendData(request)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	body := response.GetDataAsString()
	fmt.Println("sys.version", "len:", len(body), "data:", body, "FROM:", response.Info.Name)

	response = client.Send("get.boolean", nil)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	body = response.GetDataAsString()
	fmt.Println("get.boolean", "len:", len(body), "data:", lygo_conv.ToBool(body), "FROM:", response.Info.Name)

	response = client.Send("get.array", nil)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	body = response.GetDataAsString()
	fmt.Println("get.array", len(body), string(body), "FROM:", response.Info.Name)

	response = client.Send("get.error", nil)
	if response.HasError() {
		fmt.Println("get.error", len(body), response.Error, "FROM:", response.Info.Name)
	} else {
		t.Error("Expecting an error here!!!")
		t.FailNow()
	}

	response = client.Send("get.file", map[string]interface{}{"file": "./config.client.json"})
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	body = response.GetDataAsString()
	fmt.Println("get.file", len(body), string(body), "FROM:", response.Info.Name)

	// invoke internal command
	command := &lygo_n_commons.Command{
		AppToken:  "",
		Namespace: "n",
		Function:  "sys_app_token",
		Params:    nil,
	}
	response = client.SendCommand(command)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	body = response.GetDataAsString()
	appToken := string(body)
	fmt.Println("n.sys_app_token", appToken, "FROM:", response.Info.Name)

	// app.Join()
	fmt.Println("EXITING...")
}

func TestNodeNoNetworks(t *testing.T) {
	node := lygo_n.NewNode(config())
	if nil == node {
		t.FailNow()
	}
	node.Settings.Server.Http.Enabled = false // no http support
	node.Settings.Name = "SELF-HOSTED"

	// open node
	errs := node.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	registerCommands(node)

	// test LOCAL command
	response := node.Execute("get.version", nil)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	if nil != response {
		body := response.GetDataAsString()
		fmt.Println("Response to LOCAL command:", "sys.version", string(body), "FROM: LOCAL")
	} else {
		t.Error("Missing LOCAL response")
		t.FailNow()
	}

	// wait internal server starts
	time.Sleep(3 * time.Second)

	// test node command
	response = node.Send("get.version", nil)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	if nil != response {
		body := response.GetDataAsString()
		fmt.Println("Response to command:", "sys.version", string(body), "FROM:", response.Info.Name)
	}



	errs = node.Stop()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	// time.Sleep(1 * time.Second)
}

func TestMultipleNodes(t *testing.T) {

	// publisher node
	fmt.Println("* PUBLISHER", "node10010")
	node10010 := lygo_n.NewNode(config())
	node10010.Settings.Name = "node10010"
	node10010.Settings.Discovery.Publisher.Enabled = true
	node10010.Settings.Discovery.NetworkId = ""
	node10010.Settings.Discovery.Publish.Enabled = false
	node10010.Settings.Discovery.Publish.Address = "localhost:10010"
	node10010.Settings.Workspace = "./_workspace/10010"
	node10010.Settings.Server.Http.Hosts = nil // disable HTTP
	node10010.Settings.Server.Nio.Address = ":10010"
	node10010.Events().On(lygo_n_commons.EventQuitDiscovery, func(event *lygo_events.Event) {
		fmt.Println("node10010", lygo_n_commons.EventQuitDiscovery, event.ArgumentAsString(0))
	})
	errs := node10010.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	// fmt.Println(node10010.GetStatus())

	NETWORK_ID := "net_01"

	// simple node
	fmt.Println("* NODE", "node10001")
	node10001 := lygo_n.NewNode(config())
	node10001.Settings.Name = "node10001"
	node10001.Settings.Discovery.Publisher.Enabled = false
	node10001.Settings.Discovery.NetworkId = NETWORK_ID
	node10001.Settings.Discovery.Publishers = []lygo_n_commons.NAddress{"localhost:10010"}
	node10001.Settings.Discovery.Publish.Enabled = true
	node10001.Settings.Discovery.Publish.Address = "localhost:10001"
	node10001.Settings.Workspace = "./_workspace/10001"
	node10001.Settings.Server.Http.Hosts = nil
	node10001.Settings.Server.Nio.Address = ":10001"
	node10001.Events().On(lygo_n_commons.EventNewPublisher, func(event *lygo_events.Event) {
		fmt.Println("node10001", lygo_n_commons.EventNewPublisher, event.ArgumentAsString(0), event.ArgumentAsString(1), event.ArgumentAsString(2))
	})
	node10001.Events().On(lygo_n_commons.EventRemovedPublisher, func(event *lygo_events.Event) {
		fmt.Println("node10001", lygo_n_commons.EventRemovedPublisher, event.ArgumentAsString(0))
	})
	node10001.Events().On(lygo_n_commons.EventQuitDiscovery, func(event *lygo_events.Event) {
		fmt.Println("node10001", lygo_n_commons.EventQuitDiscovery, event.ArgumentAsString(0))
	})
	errs = node10001.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	// fmt.Println(node10001.GetStatus())
	registerCommands(node10001)

	// simple node
	fmt.Println("* NODE", "node10002")
	node10002 := lygo_n.NewNode(config())
	node10002.Settings.Name = "node10002"
	node10002.Settings.Discovery.Publisher.Enabled = false
	node10002.Settings.Discovery.NetworkId = NETWORK_ID
	node10002.Settings.Discovery.Publishers = []lygo_n_commons.NAddress{"localhost:10010"}
	node10002.Settings.Discovery.Publish.Enabled = true
	node10002.Settings.Discovery.Publish.Address = "localhost:10002"
	node10002.Settings.Workspace = "./_workspace/10002"
	node10002.Settings.Server.Http.Hosts = nil
	node10002.Settings.Server.Nio.Address = ":10002"
	node10002.Events().On(lygo_n_commons.EventNewPublisher, func(event *lygo_events.Event) {
		fmt.Println("node10002", lygo_n_commons.EventNewPublisher, event.ArgumentAsString(0), event.ArgumentAsString(1), event.ArgumentAsString(2))
	})
	node10002.Events().On(lygo_n_commons.EventRemovedPublisher, func(event *lygo_events.Event) {
		fmt.Println("node10002", lygo_n_commons.EventRemovedPublisher, event.ArgumentAsString(0))
	})
	node10002.Events().On(lygo_n_commons.EventQuitDiscovery, func(event *lygo_events.Event) {
		fmt.Println("node10002", lygo_n_commons.EventQuitDiscovery, event.ArgumentAsString(0))
	})
	errs = node10002.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	registerCommands(node10002)
	//fmt.Println(node10002.GetStatus())

	fmt.Println("* NODE", "node10003")
	node10003 := lygo_n.NewNode(config())
	node10003.Settings.Name = "node10003"
	node10003.Settings.Discovery.Publisher.Enabled = false
	node10003.Settings.Discovery.NetworkId = NETWORK_ID
	node10003.Settings.Discovery.Publishers = []lygo_n_commons.NAddress{"localhost:10010"}
	node10003.Settings.Discovery.Publish.Enabled = true
	node10003.Settings.Discovery.Publish.Address = "localhost:10003"
	node10003.Settings.Workspace = "./_workspace/10003"
	// node10003.Settings.Server.Http.Hosts = nil // http support
	node10003.Settings.Server.Nio.Address = ":10003"
	errs = node10003.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	registerCommands(node10003)

	// wait node sync
	time.Sleep(5 * time.Second)

	// stop publisher
	// fmt.Println("STOPPING PUBLISHER")
	// node10010.Stop()

	// node1 send a command and node 2 should respond
	response := node10001.Send("n.sys_version", nil)
	if response.HasError() {
		t.Error(response.Error)
		t.FailNow()
	}
	if nil == response {
		t.Error("Expecting a response from node10001")
		t.FailNow()
	}
	body := response.GetDataAsString()
	fmt.Println("n.sys_version", string(body), "from", response.Info.Name)

	// detach handlers
	fmt.Println("DETACH EVENT HANDLERS...")
	node10001.Events().Off(lygo_n_commons.EventNewPublisher)
	node10001.Events().Off(lygo_n_commons.EventQuitDiscovery)
	node10001.Events().Off(lygo_n_commons.EventRemovedPublisher)
	node10002.Events().Off(lygo_n_commons.EventNewPublisher)
	node10002.Events().Off(lygo_n_commons.EventQuitDiscovery)
	node10002.Events().Off(lygo_n_commons.EventRemovedPublisher)

	fmt.Println("START LOOP ASYNC COMMAND CALL....")
	count1 := 0
	count2 := 0
	count3 := 0
	go func() {
		for {
			time.Sleep(lygo_rnd.BetweenDuration(50, 500) * time.Millisecond)
			go func() {
				response := node10001.Send("n.sys_version", nil)
				//fmt.Println(response.Info.Name, response.GetDataAsString())
				if response.Info.Name == "node10001" {
					count1++
				} else if response.Info.Name == "node10002" {
					count2++
				} else {
					count3++
				}
			}()
		}
	}()

	go func() {
		for {
			time.Sleep(lygo_rnd.BetweenDuration(100, 300) * time.Millisecond)
			go func() {
				response := node10001.Send("n.sys_version", nil)
				// fmt.Println(response.Info.Name, response.GetDataAsString())
				if response.Info.Name == "node10001" {
					count1++
				} else if response.Info.Name == "node10002" {
					count2++
				} else {
					count3++
				}
			}()
		}
	}()

	fmt.Println("WAITING 10 SECONDS...")
	time.Sleep(10 * time.Second)
	fmt.Println("\tnode10001:", count1)
	fmt.Println("\tnode10002", count2)
	fmt.Println("\tnode10003", count3)

	time.Sleep(10 * time.Minute)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func config() *lygo_n_commons.NSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.node.json")
	config := new(lygo_n_commons.NSettings)
	config.Parse(text_cfg)

	return config
}

func configSrv() *lygo_n_commons.NHostSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.server.json")
	config := new(lygo_n_commons.NHostSettings)
	config.Parse(text_cfg)

	return config
}

func configCli() *lygo_n_commons.NConnSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.client.json")
	config := new(lygo_n_commons.NConnSettings)
	config.Parse(text_cfg)

	return config
}

func initializeHttp(http *lygo_http_server.HttpServer) {
	http.Route.Get("*", func(ctx *fiber.Ctx) {
		ctx.Write("ROOT\n")
		ctx.Write(ctx.BaseURL())
		ctx.Write(ctx.OriginalURL())
		ctx.Next()
	})
	g := http.Route.Group("/api", func(ctx *fiber.Ctx) {
		id := ctx.Params("id")

		ctx.Write("/api\n")
		ctx.Write(id)
		//ctx.SendBytes([]byte("THIS IS GROUP API\n"))
		ctx.Next()
	})
	g.Get("/v1/:id", func(ctx *fiber.Ctx) {
		id := ctx.Params("id")
		ctx.Write("/v1\n")
		ctx.Write("THIS IS v1")
		ctx.Write(id)
	})
}

func registerCommands(n *lygo_n.N) {
	n.RegisterCommand("get.version", func(message *lygo_n_commons.Command) interface{} {
		return "1.0.2"
	})
	n.RegisterCommand("get.boolean", func(message *lygo_n_commons.Command) interface{} {
		return true
	})
	n.RegisterCommand("get.array", func(message *lygo_n_commons.Command) interface{} {
		return []interface{}{"1", "2", "3", 4, 5, 6, true, false}
	})
	n.RegisterCommand("get.error", func(message *lygo_n_commons.Command) interface{} {
		return errors.New("ERROR SIMULATION")
	})
	n.RegisterCommand("get.file", func(message *lygo_n_commons.Command) interface{} {
		data, err := lygo_io.ReadBytesFromFile(message.GetParamAsString("file"))
		if nil != err {
			return err
		}
		return data
	})
}
