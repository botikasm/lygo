package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/gofiber/fiber"
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

	server.Route.Get("/get", func(ctx *fiber.Ctx) {
		ctx.SendBytes([]byte("THIS IS GET API"))
	})

	g := server.Route.Group("/api", func(ctx *fiber.Ctx) {
		ctx.SendBytes([]byte("THIS IS GROUP API"))
		ctx.Next()
	})
	g.Get("/v1", func(ctx *fiber.Ctx) {
		ctx.SendBytes([]byte("THIS IS GROUP API v1"))
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

func config() *lygo_http_server.HttpServerConfig {
	text_cfg, _ := lygo_io.ReadTextFromFile("./lygo_http_server_config.json")
	config := new(lygo_http_server.HttpServerConfig)
	config.Parse(text_cfg)

	return config
}

func onError(errCtx *lygo_http_server.HttpServerError) {
	fmt.Println(errCtx.Message, errCtx.Error.Error())
}
