package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_types"
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {

	// load configuration
	config := config()
	if nil == config {
		t.Errorf("Configuration is not valid")
	}

	server := lygo_http_server.NewHttpServer(config)
	server.CallbackError = onError

	server.Route.Get("*", func(ctx *fiber.Ctx) {
		ctx.Write("ROOT\n")
		ctx.Next()
	})

	server.Route.Get("/get", func(ctx *fiber.Ctx) {
		ctx.Write(fmt.Sprintf("Hi, I'm worker #%v", os.Getpid()))
		// ctx.SendBytes([]byte("THIS IS GET API"))
	})

	g := server.Route.Group("/api", func(ctx *fiber.Ctx) {
		ctx.SendBytes([]byte("THIS IS GROUP API\n"))
		ctx.Next()
	})
	g.Get("/v1", func(ctx *fiber.Ctx) {
		ctx.Write("/v1\n")
		ctx.Write("THIS IS v1")
	})

	server.Middleware("/", func(ctx *fiber.Ctx){
		fmt.Println("First middleware")
		ctx.Append("middleware", "First middleware")
		ctx.Next()
	})

	server.Middleware(func(ctx *fiber.Ctx){
		fmt.Println("Second middleware")
		ctx.Append("middleware", "Second middleware")
		ctx.Next()
	})

	server.Middleware("/api", func(ctx *fiber.Ctx){
		fmt.Println("API middleware")
		ctx.Write("API middleware")
		ctx.Append("middleware", "API middleware")
		ctx.Next()
	})

	server.Websocket("/", func(c *websocket.Conn) {
		fmt.Println("LOCALS", c.Locals("Hello")) // "World"
		// Websocket stuff
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("ERROR read:", err)
				break
			}
			fmt.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				fmt.Println("ERROR write:", err)
				break
			}
		}
	})

	server.Middleware("/yoda", func(ctx *fiber.Ctx){
		// NOT FOUND
		b, _ :=lygo_io.ReadBytesFromFile("./www/yoda.jpeg")
		ctx.SendBytes(b)
	})

	// start server
	err := server.Start()
	if nil != err {
		t.Error(err)
	}

	// Wait forever.
	server.Join()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func config() *lygo_http_server_config.HttpServerConfig {
	text_cfg, _ := lygo_io.ReadTextFromFile("./lygo_http_server_config.json")
	config := new(lygo_http_server_config.HttpServerConfig)
	config.Parse(text_cfg)

	return config
}

func onError(errCtx *lygo_http_server_types.HttpServerError) {
	fmt.Println(errCtx.Message, errCtx.Error.Error())
}
