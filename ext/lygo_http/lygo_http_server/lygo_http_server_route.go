package lygo_http_server

import (
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/gofiber/fiber"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type HttpServerConfigRoute struct {
	data map[string]interface{}
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

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerConfigRoute
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServerConfigRoute() *HttpServerConfigRoute {
	instance := new(HttpServerConfigRoute)
	instance.data = make(map[string]interface{})
	return instance
}

func (instance *HttpServerConfigRoute) Group(path string, handlers ...func(ctx *fiber.Ctx)) *HttpServerConfigGroup {
	m := instance.data
	g := &HttpServerConfigGroup{
		Path:     path,
		Handlers: handlers,
	}
	g.Children = make([]interface{}, 0)
	m[buildKey("GROUP", path)] = g
	return g
}

func (instance *HttpServerConfigRoute) All(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.data
	m[buildKey("ALL", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *HttpServerConfigRoute) Get(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.data
	m[buildKey("GET", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *HttpServerConfigRoute) Post(path string, handlers ...func(ctx *fiber.Ctx)) {
	m := instance.data
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
	m := g.data
	m[buildKey("ALL", path)] = &HttpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
	instance.Children = append(instance.Children, g)
}

func (instance *HttpServerConfigGroup) Get(path string, handlers ...func(ctx *fiber.Ctx)) {
	g := NewHttpServerConfigRoute()
	m := g.data
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
