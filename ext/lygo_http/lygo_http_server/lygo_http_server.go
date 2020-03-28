package lygo_http_server

import (
	"crypto/tls"
	"errors"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/gofiber/compression"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/limiter"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"
	"github.com/gofiber/websocket"
	"sync"
	"time"
)

var (
	errorInvalidConfiguration = errors.New("configuration_invalid")
)

const maxErrors = 100

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServer struct {
	Config *HttpServerConfig

	//-- hooks --//
	CallbackError        ErrorHookCallback
	CallbackLimitReached func(ctx *fiber.Ctx)

	// ROUTING
	Route *HttpServerConfigRoute

	//-- private --//
	_hosts     map[string]*fiber.Fiber
	middleware []*HttpServerConfigRouteItem
	websocket  []*HttpServerConfigRouteWebsocket
	started    bool
	stopped    bool
	errors     []error
	muxError   sync.Mutex
}

type HttpServerError struct {
	Server  *HttpServer
	Message string
	Error   error
	Context *fiber.Ctx
}

type ErrorHookCallback func(*HttpServerError)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServer(config *HttpServerConfig) *HttpServer {
	instance := new(HttpServer)
	instance.Config = config
	instance.Route = NewHttpServerConfigRoute()
	instance.middleware = make([]*HttpServerConfigRouteItem, 0)
	instance.websocket = make([]*HttpServerConfigRouteWebsocket, 0)
	instance.stopped = false
	instance.started = false
	instance._hosts = make(map[string]*fiber.Fiber)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Start() error {
	if !instance.started {
		if nil == instance.Config {
			return errorInvalidConfiguration
		}
		instance.started = true
		instance.stopped = false

		instance.initConfig()

		instance.serve()

	}
	return nil
}

func (instance *HttpServer) Join() {
	if !instance.stopped {
		for !instance.stopped {
			time.Sleep(10 * time.Second)
		}
	}
}

func (instance *HttpServer) Stop() {
	if !instance.stopped {
		for _, host := range instance._hosts {
			if nil != host {
				_ = host.Shutdown()
			}
		}
		instance.stopped = true
		instance.started = false
	}
}

func (instance *HttpServer) Middleware(args ...interface{}) {
	item := new(HttpServerConfigRouteItem)
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
		if v, b := args[0].(func(ctx *websocket.Conn)); b {
			item.Path = "/ws"
			item.Handler = v
		}
	case 2:
		if path, b := args[0].(string); b {
			if f, b := args[1].(func(ctx *websocket.Conn)); b {
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

	for _, host := range config.Hosts {
		key := host.Address
		if _, ok := instance._hosts[key]; !ok {
			app := fiber.New()
			instance.configure(app, host, config)
			instance._hosts[key] = app
			go instance.listen(host, app)
		}
	}
}

func (instance *HttpServer) listen(host *HttpServerConfigHost, app *fiber.Fiber) {
	var tlsConfig *tls.Config
	if host.TLS && len(host.SslKey) > 0 && len(host.SslCert) > 0 {
		cer, err := tls.LoadX509KeyPair(host.SslCert, host.SslKey)
		if err != nil {
			instance.doError("Error loading Certificates", err, nil)
		}
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	}

	if nil == tlsConfig {
		if err := app.Listen(host.Address); err != nil {
			instance.doError(lygo_strings.Format("Error Opening channel: '%s'", host.Address), err, nil)
		}
	} else {
		if err := app.Listen(host.Address, tlsConfig); err != nil {
			instance.doError(lygo_strings.Format("Error Opening TLS channel: '%s'", host.Address), err, nil)
		}
	}
}

func (instance *HttpServer) configure(app *fiber.Fiber, host *HttpServerConfigHost, config *HttpServerConfig) {
	if nil != app {
		cfg := recover.Config{
			Handler: instance.onServerError,
		}
		app.Use(recover.New(cfg))

		if config.EnableRequestId {
			app.Use(requestid.New())
		}

		app.Settings.ServerHeader = config.ServerHeader
		app.Settings.Prefork = config.Prefork
		app.Settings.CaseSensitive = config.CaseSensitive
		app.Settings.StrictRouting = config.StrictRouting
		app.Settings.Immutable = config.Immutable
		if config.BodyLimit > 0 {
			app.Settings.BodyLimit = config.BodyLimit
		}
		if config.ReadTimeout > 0 {
			app.Settings.ReadTimeout = config.ReadTimeout * time.Millisecond
		}
		if config.WriteTimeout > 0 {
			app.Settings.WriteTimeout = config.WriteTimeout * time.Millisecond
		}
		if config.IdleTimeout > 0 {
			app.Settings.IdleTimeout = config.IdleTimeout * time.Millisecond
		}

		// CORS
		initCORS(app, instance.Config.CORS)

		// compression
		initCompression(app, instance.Config.Compression)

		// limiter
		initLimiter(app, instance.Config.Limiter, instance.onLimitReached)

		// Middleware
		if len(instance.middleware) > 0 {
			initMiddleware(app, instance.middleware)
		}

		// Route
		if nil != instance.Route {
			initRoute(app, instance.Route, nil)
		}

		// websocket
		initSocket(app, host, instance.websocket)

		// Static
		if len(config.Static) > 0 {
			for _, static := range config.Static {
				if static.Enabled && len(static.Prefix) > 0 && len(static.Root) > 0 {
					app.Static(static.Prefix, static.Root, fiber.Static{
						Compress:  static.Compress,
						ByteRange: static.ByteRange,
						Browse:    static.Browse,
						Index:     static.Index,
					})
				}
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
			instance.CallbackError(&HttpServerError{
				Server:  instance,
				Message: message,
				Context: ctx,
				Error:   err,
			})
		}
	}()
}

func (instance *HttpServer) onServerError(c *fiber.Ctx, err error) {
	c.SendString(err.Error())
	c.SendStatus(500)
}

func (instance *HttpServer) onLimitReached(c *fiber.Ctx) {
	// request limit reached
	if nil != instance.CallbackLimitReached {
		instance.CallbackLimitReached(c)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	s t a t i c
//----------------------------------------------------------------------------------------------------------------------

func initSocket(app *fiber.Fiber, item *HttpServerConfigHost, routes []*HttpServerConfigRouteWebsocket) {
	if nil != item.Websocket && item.Websocket.Enabled {
		settings := item.Websocket
		config := websocket.Config{}
		config.EnableCompression = settings.EnableCompression
		if settings.HandshakeTimeout > 0 {
			config.HandshakeTimeout = settings.HandshakeTimeout * time.Millisecond
		}
		if len(settings.Origins) > 0 {
			config.Origins = settings.Origins
		} else {
			config.Origins = []string{"*"}
		}
		if len(settings.Subprotocols) > 0 {
			config.Subprotocols = settings.Subprotocols
		}
		if settings.ReadBufferSize > 0 {
			config.ReadBufferSize = settings.ReadBufferSize
		}
		if settings.WriteBufferSize > 0 {
			config.WriteBufferSize = settings.WriteBufferSize
		}

		// open websocket handlers
		for _, route:=range routes{
			if len(route.Path)>0{
				app.Get(route.Path, websocket.New(route.Handler, config))
			}
		}
	}
}

func initCORS(app *fiber.Fiber, corsCfg *HttpServerConfigCORS) {
	if nil != corsCfg && corsCfg.Enabled {
		config := cors.Config{}
		if corsCfg.MaxAge > 0 {
			config.MaxAge = corsCfg.MaxAge
		}
		if corsCfg.AllowCredentials {
			config.AllowCredentials = true
		}
		if len(corsCfg.AllowMethods) > 0 {
			config.AllowMethods = corsCfg.AllowMethods
		}
		if len(corsCfg.AllowOrigins) > 0 {
			config.AllowOrigins = corsCfg.AllowOrigins
		}
		if len(corsCfg.ExposeHeaders) > 0 {
			config.ExposeHeaders = corsCfg.ExposeHeaders
		}
		app.Use(cors.New(config))
	}
}

func initCompression(app *fiber.Fiber, cfg *HttpServerConfigCompression) {
	if nil != cfg && cfg.Enabled {
		config := compression.Config{}
		config.Level = cfg.Level

		app.Use(compression.New(config))
	}
}

func initLimiter(app *fiber.Fiber, cfg *HttpServerConfigLimiter, handler func(ctx *fiber.Ctx)) {
	if nil != cfg && cfg.Enabled {
		config := limiter.Config{}
		if cfg.Timeout > 0 {
			config.Timeout = cfg.Timeout
		}
		if cfg.Max > 0 {
			config.Max = cfg.Max
		}
		if len(cfg.Message) > 0 {
			config.Message = cfg.Message
		}
		if cfg.StatusCode > 0 {
			config.StatusCode = cfg.StatusCode
		}

		config.Handler = handler

		app.Use(limiter.New(config))
	}
}

func initMiddleware(app *fiber.Fiber, items []*HttpServerConfigRouteItem) {
	for _, item := range items {
		path := item.Path
		if len(path) == 0 {
			app.Use(item.Handlers[0])
		} else {
			app.Use(path, item.Handlers[0])
		}
	}
}

func initRoute(app *fiber.Fiber, route *HttpServerConfigRoute, parent *fiber.Group) {
	for k, i := range route.data {
		initRouteItem(app, k, i, parent)
	}
}

func initGroup(app *fiber.Fiber, group *HttpServerConfigGroup, parent *fiber.Group) {
	var g *fiber.Group
	if nil == parent {
		g = app.Group(group.Path, group.Handlers...)
	} else {
		g = parent.Group(group.Path, group.Handlers...)
	}
	if nil != g && len(group.Children) > 0 {
		for _, c := range group.Children {
			if cc, b := c.(*HttpServerConfigGroup); b {
				// children is a group
				initGroup(app, cc, g)
			} else if cc, b := c.(*HttpServerConfigRoute); b {
				// children is route
				initRoute(app, cc, g)
			}
		}
	}
}

func initRouteItem(app *fiber.Fiber, key string, item interface{}, parent *fiber.Group) {
	method := lygo_strings.SplitAndGetAt(key, "_", 0)
	switch method {
	case "GROUP":
		v := item.(*HttpServerConfigGroup)
		initGroup(app, v, parent)
	case "ALL":
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.All(v.Path, v.Handlers...)
		} else {
			parent.All(v.Path, v.Handlers...)
		}
	case fiber.MethodGet:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Get(v.Path, v.Handlers...)
		} else {
			parent.Get(v.Path, v.Handlers...)
		}
	case fiber.MethodPost:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Post(v.Path, v.Handlers...)
		} else {
			parent.Post(v.Path, v.Handlers...)
		}
	case fiber.MethodOptions:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Options(v.Path, v.Handlers...)
		} else {
			parent.Options(v.Path, v.Handlers...)
		}
	case fiber.MethodPut:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Put(v.Path, v.Handlers...)
		} else {
			parent.Put(v.Path, v.Handlers...)
		}
	case fiber.MethodHead:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Head(v.Path, v.Handlers...)
		} else {
			parent.Head(v.Path, v.Handlers...)
		}
	case fiber.MethodPatch:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Patch(v.Path, v.Handlers...)
		} else {
			parent.Patch(v.Path, v.Handlers...)
		}
	case fiber.MethodDelete:
		v := item.(*HttpServerConfigRouteItem)
		if nil == parent {
			app.Delete(v.Path, v.Handlers...)
		} else {
			parent.Delete(v.Path, v.Handlers...)
		}
	}
}
