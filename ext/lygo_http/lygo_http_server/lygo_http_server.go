package lygo_http_server

import (
	"errors"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_types"
	"github.com/gofiber/fiber"
	"sync"
)

var (
	errorInvalidConfiguration = errors.New("configuration_invalid")
)

const maxErrors = 100

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServer struct {
	Config *lygo_http_server_config.HttpServerConfig

	//-- hooks --//
	CallbackError        lygo_http_server_types.CallbackError
	CallbackLimitReached lygo_http_server_types.CallbackLimitReached

	// ROUTING
	Route *lygo_http_server_config.HttpServerConfigRoute

	//-- private --//
	enabled    bool // len(hosts)>0
	services   map[string]*HttpServerService
	middleware []*lygo_http_server_config.HttpServerConfigRouteItem
	websocket  []*HttpServerConfigRouteWebsocket
	started    bool
	stopped    bool
	errors     []error
	muxError   sync.Mutex
	stopChan   chan bool
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServer(config *lygo_http_server_config.HttpServerConfig) *HttpServer {

	instance := new(HttpServer)
	instance.Config = config
	instance.Route = lygo_http_server_config.NewHttpServerConfigRoute()
	instance.middleware = make([]*lygo_http_server_config.HttpServerConfigRouteItem, 0)
	instance.websocket = make([]*HttpServerConfigRouteWebsocket, 0)
	instance.stopped = false
	instance.started = false
	instance.services = make(map[string]*HttpServerService)
	instance.stopChan = make(chan bool, 1)

	instance.enabled = len(config.Hosts) > 0

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) IsEnabled() bool {
	if nil != instance {
		return instance.enabled
	}
	return false
}

func (instance *HttpServer) Start() error {
	if nil != instance {
		if instance.enabled && !instance.started {
			if nil == instance.Config {
				return errorInvalidConfiguration
			}
			instance.started = true
			instance.stopped = false

			instance.initConfig()

			instance.serve()
		}
	}
	return nil
}

func (instance *HttpServer) Join() {
	if instance.enabled && !instance.stopped {
		<-instance.stopChan
	}
}

func (instance *HttpServer) Stop() {
	if instance.enabled && !instance.stopped {
		for _, service := range instance.services {
			if nil != service {
				_ = service.Shutdown()
			}
		}
		instance.stopped = true
		instance.started = false
		instance.stopChan <- true
	}
}

func (instance *HttpServer) Middleware(args ...interface{}) {
	item := new(lygo_http_server_config.HttpServerConfigRouteItem)
	switch len(args) {
	case 1:
		if v, b := args[0].(func(ctx *fiber.Ctx)); b {
			item.Path = ""
			item.Handlers = append(item.Handlers, v)
		}
	case 2:
		if path, b := args[0].(string); b {
			if f, b := args[1].(func(ctx *fiber.Ctx)); b {
				item.Path = path
				item.Handlers = append(item.Handlers, f)
			}
		}
	}

	if len(item.Handlers) > 0 {
		instance.middleware = append(instance.middleware, item)
	}
}

func (instance *HttpServer) Websocket(args ...interface{}) {
	item := new(HttpServerConfigRouteWebsocket)
	switch len(args) {
	case 1:
		if v, b := args[0].(func(ctx *HttpWebsocketConn)); b {
			item.Path = instance.Config.WebsocketEndpoint
			if len(item.Path) == 0 {
				item.Path = "/ws"
			}
			item.Handler = v
		}
	case 2:
		if path, b := args[0].(string); b {
			if f, b := args[1].(func(ctx *HttpWebsocketConn)); b {
				item.Path = path
				item.Handler = f
			}
		}
	}

	if nil != item.Handler {
		instance.websocket = append(instance.websocket, item)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) serve() {
	// configure
	config := instance.Config

	if instance.enabled {
		for _, host := range config.Hosts {
			key := host.Address
			if _, ok := instance.services[key]; !ok {
				// creates service and add to internal pool
				service := NewServerService(key,
					instance.Config, host, instance.Route, instance.middleware, instance.websocket,
					instance.onEndpointError, instance.CallbackLimitReached)
				service.Open()
				instance.services[key] = service
			}
		}
	}
}

func (instance *HttpServer) initConfig() {
	config := instance.Config
	if nil != config {
		for _, static := range config.Static {
			if len(static.Index) == 0 {
				static.Index = "index.html"
			}
		}
	}
}

func (instance *HttpServer) doError(message string, err error, ctx *fiber.Ctx) {
	instance.muxError.Lock()
	go func() {
		defer instance.muxError.Unlock()
		if nil != instance.errors && len(instance.errors) > maxErrors {
			// reset errors
			instance.errors = make([]error, 0)
		}
		instance.errors = append(instance.errors, err)
		if nil != instance.CallbackError {
			instance.CallbackError(&lygo_http_server_types.HttpServerError{
				Sender:  instance,
				Message: message,
				Context: ctx,
				Error:   err,
			})
		}
	}()
}

func (instance *HttpServer) onEndpointError(err *lygo_http_server_types.HttpServerError) {
	instance.muxError.Lock()
	go func() {
		defer instance.muxError.Unlock()
		if nil != instance.errors && len(instance.errors) > maxErrors {
			// reset errors
			instance.errors = make([]error, 0)
		}
		instance.errors = append(instance.errors, err.Error)
		if nil != instance.CallbackError {
			instance.CallbackError(err)
		}
	}()
}
