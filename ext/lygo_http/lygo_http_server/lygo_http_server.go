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
	"log"
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
	_http    *fiber.Fiber
	_https   *fiber.Fiber
	started  bool
	stopped  bool
	errors   []error
	muxError sync.Mutex
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
	instance.stopped = false
	instance.started = false

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
		if nil != instance._http {
			_ = instance._http.Shutdown()
		}
		if nil != instance._https {
			_ = instance._https.Shutdown()
		}
		instance.stopped = true
		instance.started = false
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) http() *fiber.Fiber {
	if nil == instance._http {
		instance._http = fiber.New()
		instance.configure(instance._http, instance.Config)
	}
	return instance._http
}

func (instance *HttpServer) https() *fiber.Fiber {
	if nil == instance._https {
		instance._https = fiber.New()
		instance.configure(instance._https, instance.Config)
	}
	return instance._https
}

func (instance *HttpServer) configure(app *fiber.Fiber, config *HttpServerConfig) {
	if nil != app {
		cfg := recover.Config{
			Handler: instance.onServerError,
		}
		app.Use(recover.New(cfg))

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

		// Route
		if nil != instance.Route {
			initRoute(app, instance.Route, nil)
		}

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

func (instance *HttpServer) serve() {
	// configure
	config := instance.Config

	if len(config.Address) > 0 {
		// http
		go instance.listenAndServe(config.Address)
	}
	if len(config.AddressTLS) > 0 {
		// https
		go instance.listenAndServeTLS(config.AddressTLS, config.SslCert, config.SslKey)
	}
}

func (instance *HttpServer) listenAndServe(addr string) {
	if err := instance.http().Listen(addr); err != nil {
		instance.doError(lygo_strings.Format("Error Opening channel: '%s'", addr), err, nil)
	}
}

func (instance *HttpServer) listenAndServeTLS(addr string, certFile string, keyFile string) {
	if len(certFile) > 0 && len(keyFile) > 0 {

		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		// use files
		if err := instance.https().Listen(addr, config); err != nil {
			instance.doError(lygo_strings.Format("Error Opening TLS channel: '%s'", addr), err, nil)
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
