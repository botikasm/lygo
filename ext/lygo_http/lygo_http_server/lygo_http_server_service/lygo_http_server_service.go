package lygo_http_server_service

import (
	"crypto/tls"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_types"
	"github.com/gofiber/compression"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/limiter"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServerService struct {
	Key string

	//-- private --//
	app *fiber.Fiber

	config           *lygo_http_server_config.HttpServerConfig
	configHost       *lygo_http_server_config.HttpServerConfigHost
	configRoute      *lygo_http_server_config.HttpServerConfigRoute
	configMiddleware []*lygo_http_server_config.HttpServerConfigRouteItem
	configWebsocket  []*lygo_http_server_config.HttpServerConfigRouteWebsocket

	callbackError        lygo_http_server_types.CallbackError
	callbackLimitReached lygo_http_server_types.CallbackLimitReached
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerEndpoint
//----------------------------------------------------------------------------------------------------------------------

func NewServerService(key string,
	config *lygo_http_server_config.HttpServerConfig,
	host *lygo_http_server_config.HttpServerConfigHost,
	route *lygo_http_server_config.HttpServerConfigRoute,
	middleware []*lygo_http_server_config.HttpServerConfigRouteItem,
	websocket []*lygo_http_server_config.HttpServerConfigRouteWebsocket,
	callbackError lygo_http_server_types.CallbackError,
	callbackLimitReached lygo_http_server_types.CallbackLimitReached) *HttpServerService {

	instance := new(HttpServerService)
	instance.Key = key

	instance.app = fiber.New()
	instance.config = config
	instance.configHost = host
	instance.configRoute = route
	instance.configMiddleware = middleware
	instance.configWebsocket = websocket

	instance.callbackError = callbackError
	instance.callbackLimitReached = callbackLimitReached

	return instance
}

func (instance *HttpServerService) Shutdown() error {
	if nil != instance && nil != instance.app {
		return instance.app.Shutdown()
	}
	return nil
}

func (instance *HttpServerService) Open() {
	if nil != instance && nil != instance.app {
		instance.init()
		go instance.listen()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServerService) init() {
	app := instance.app
	config := instance.config
	host := instance.configHost
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
		initCORS(app, instance.config.CORS)

		// compression
		initCompression(app, instance.config.Compression)

		// limiter
		initLimiter(app, instance.config.Limiter, instance.onLimitReached)

		// Middleware
		if len(instance.configMiddleware) > 0 {
			initMiddleware(app, instance.configMiddleware)
		}

		// Route
		if nil != instance.configRoute {
			initRoute(app, instance.configRoute, nil)
		}

		// websocket
		socket := NewHttpWebsocket(app, host, instance.configWebsocket)
		socket.Init()

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

func (instance *HttpServerService) listen() {
	app := instance.app
	host := instance.configHost
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

func (instance *HttpServerService) onServerError(c *fiber.Ctx, err error) {
	c.SendString(err.Error())
	c.SendStatus(500)
}

func (instance *HttpServerService) onLimitReached(c *fiber.Ctx) {
	// request limit reached
	if nil != instance.callbackLimitReached {
		instance.callbackLimitReached(c)
	}
}

func (instance *HttpServerService) doError(message string, err error, ctx *fiber.Ctx) {
	go func() {
		if nil != instance.callbackError {
			instance.callbackError(&lygo_http_server_types.HttpServerError{
				Sender:  instance,
				Message: message,
				Context: ctx,
				Error:   err,
			})
		}
	}()
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func initCORS(app *fiber.Fiber, corsCfg *lygo_http_server_config.HttpServerConfigCORS) {
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

func initCompression(app *fiber.Fiber, cfg *lygo_http_server_config.HttpServerConfigCompression) {
	if nil != cfg && cfg.Enabled {
		config := compression.Config{}
		config.Level = cfg.Level

		app.Use(compression.New(config))
	}
}

func initLimiter(app *fiber.Fiber, cfg *lygo_http_server_config.HttpServerConfigLimiter, handler func(ctx *fiber.Ctx)) {
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

func initMiddleware(app *fiber.Fiber, items []*lygo_http_server_config.HttpServerConfigRouteItem) {
	for _, item := range items {
		path := item.Path
		if len(path) == 0 {
			app.Use(item.Handlers[0])
		} else {
			app.Use(path, item.Handlers[0])
		}
	}
}

func initRoute(app *fiber.Fiber, route *lygo_http_server_config.HttpServerConfigRoute, parent *fiber.Group) {
	for k, i := range route.Data {
		initRouteItem(app, k, i, parent)
	}
}

func initGroup(app *fiber.Fiber, group *lygo_http_server_config.HttpServerConfigGroup, parent *fiber.Group) {
	var g *fiber.Group
	if nil == parent {
		g = app.Group(group.Path, group.Handlers...)
	} else {
		g = parent.Group(group.Path, group.Handlers...)
	}
	if nil != g && len(group.Children) > 0 {
		for _, c := range group.Children {
			if cc, b := c.(*lygo_http_server_config.HttpServerConfigGroup); b {
				// children is a group
				initGroup(app, cc, g)
			} else if cc, b := c.(*lygo_http_server_config.HttpServerConfigRoute); b {
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
		v := item.(*lygo_http_server_config.HttpServerConfigGroup)
		initGroup(app, v, parent)
	case "ALL":
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.All(v.Path, v.Handlers...)
		} else {
			parent.All(v.Path, v.Handlers...)
		}
	case fiber.MethodGet:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Get(v.Path, v.Handlers...)
		} else {
			parent.Get(v.Path, v.Handlers...)
		}
	case fiber.MethodPost:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Post(v.Path, v.Handlers...)
		} else {
			parent.Post(v.Path, v.Handlers...)
		}
	case fiber.MethodOptions:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Options(v.Path, v.Handlers...)
		} else {
			parent.Options(v.Path, v.Handlers...)
		}
	case fiber.MethodPut:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Put(v.Path, v.Handlers...)
		} else {
			parent.Put(v.Path, v.Handlers...)
		}
	case fiber.MethodHead:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Head(v.Path, v.Handlers...)
		} else {
			parent.Head(v.Path, v.Handlers...)
		}
	case fiber.MethodPatch:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Patch(v.Path, v.Handlers...)
		} else {
			parent.Patch(v.Path, v.Handlers...)
		}
	case fiber.MethodDelete:
		v := item.(*lygo_http_server_config.HttpServerConfigRouteItem)
		if nil == parent {
			app.Delete(v.Path, v.Handlers...)
		} else {
			parent.Delete(v.Path, v.Handlers...)
		}
	}
}
