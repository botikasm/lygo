package _test_test

import (
	"errors"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/ext/lygo_n"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
	"github.com/gofiber/fiber"
	"testing"
	"time"
)

func TestSimpleCommunication(t *testing.T) {

	server := lygo_n_server.NewNServer(configSrv())

	server.OnMessage = onMessage
	server.RegisterCommand("get.version", func(message *lygo_n_commons.Command) interface{} {
		return "1.0.2"
	})
	server.RegisterCommand("get.boolean", func(message *lygo_n_commons.Command) interface{} {
		return true
	})
	server.RegisterCommand("get.error", func(message *lygo_n_commons.Command) interface{} {
		return errors.New("ERROR SIMULATION")
	})
	server.RegisterCommand("get.file", func(message *lygo_n_commons.Command) interface{} {
		data, err := lygo_io.ReadBytesFromFile(message.GetParamAsString("file"))
		if nil!=err{
			return err
		}
		return data
	})

	initialize(server)

	errs, _ := server.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	client := lygo_n_client.NewNClient(configCli())
	errs, _ = client.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	request := map[string]interface{}{
		"namespace": "get",
		"function":  "version",
	}
	response, err := client.SendData(request)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("sys.version", len(response), string(response))

	response, err = client.Send("get.boolean", nil)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("get.boolean", len(response), string(response))

	response, err = client.Send("get.error", nil)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("get.error", len(response), string(response))

	response, err = client.Send("get.file", map[string]interface{}{"file":"./client.config.json"})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("get.file", len(response), string(response))

	// invoke internal command
	command := &lygo_n_commons.Command{
		AppToken:  "",
		Namespace: "n",
		Function:  "sys_app_token",
		Params:    nil,
	}
	response, err = client.SendCommand(command)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	appToken := string(response)
	fmt.Println("n.sys_app_token", appToken)


	// app.Join()
	fmt.Println("EXITING...")
}

func TestNode(t *testing.T) {
	node := lygo_n.NewNode(config())
	if nil==node{
		t.FailNow()
	}

	// open node
	errs := node.Start()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}

	// register a command handler
	node.RegisterCommand("sys.version", func(message *lygo_n_commons.Command) interface{} {
		return "v1.0.1"
	})

	// test command
	response, err :=node.Send("sys.version", nil)
	if nil!=err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Response to command:", "sys.version", string(response))

	errs = node.Stop()
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	time.Sleep(5*time.Second)
	lygo_logs.Close()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func config() *lygo_n.NSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.node.json")
	config := new(lygo_n.NSettings)
	config.Parse(text_cfg)

	return config
}

func configSrv() *lygo_n_server.NServerSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.server.json")
	config := new(lygo_n_server.NServerSettings)
	config.Parse(text_cfg)

	return config
}

func configCli() *lygo_n_client.NClientSettings {
	text_cfg, _ := lygo_io.ReadTextFromFile("./config.client.json")
	config := new(lygo_n_client.NClientSettings)
	config.Parse(text_cfg)

	return config
}

func initialize(app *lygo_n_server.NServer) {
	app.Server().Route.Get("*", func(ctx *fiber.Ctx) {
		ctx.Write("ROOT\n")
		ctx.Write(ctx.BaseURL())
		ctx.Write(ctx.OriginalURL())
		ctx.Next()
	})
	g := app.Server().Route.Group("/api", func(ctx *fiber.Ctx) {
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

func onMessage(method string, message *lygo_n_commons.Message) (interface{}, bool) {
	var response interface{} = nil
	handled := true
	switch method {
	case "sys.version":
		response = "SHOULD NOT HANDLE THIS!!!!"
	default:
		handled = false // command not found
	}
	return response, handled
}
