package lygo_http_server_config

import (
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type HttpServerConfigRoute struct {
	Data map[string]interface{}
}

type HttpServerConfigRouteItem struct {
	Path     string
	Handlers []func(ctx *fiber.Ctx)
}

type HttpServerConfigGroup struct {
	Path     string
	Handlers []func(ctx *fiber.Ctx)
	Children []interface{}
}

type HttpServerConfigRouteWebsocket struct {
	Path     string
	Handler func(c *websocket.Conn)
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerConfigRoute
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServerConfigRoute() *HttpServerConfigRoute {
	instance := new(HttpServerConfigRoute)
	instance.Data = make(map[string]interface{})
	return instance
}

func (instance *HttpServerConfigRoute) Group(path string, handlers ...func(ctx *fiber.Ctx)) *HttpServerConfigGroup {
	m := instance.Data
	g := &HttpServerConfigGroup{
		Path:     path,
		Handlers: handlers,
	}
	g.Children = make([]interface{}, 0)
	m[buildKey("GROUP", path)] = g
	return g
}

func (instance *HttpServerConfigRoute) All(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.Data
	m[buildKey("ALL", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *HttpServerConfigRoute) Get(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.Data
	m[buildKey("GET", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *HttpServerConfigRoute) Post(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.Data
	m[buildKey("POST", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerConfigGroup
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServerConfigGroup) All(path string, handlers ...func(ctx *fiber.Ctx)) {
	g := NewHttpServerConfigRoute()
	m := g.Data
	m[buildKey("ALL", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
	instance.Children = append(instance.Children, g)
}

func (instance *HttpServerConfigGroup) Get(path string, handlers ...func(ctx *fiber.Ctx)) {
	g := NewHttpServerConfigRoute()
	m := g.Data
	m[buildKey("GET", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
	instance.Children = append(instance.Children, g)
}


//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func buildKey(method, path string) string {
	return method + "_" + lygo_crypto.MD5(path)
}
