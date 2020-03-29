package lygo_http_server_service

import (
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpWebsocket struct {

	//-- private --//
	app           *fiber.Fiber
	configService *lygo_http_server_config.HttpServerConfigHost
	configRoutes  []*lygo_http_server_config.HttpServerConfigRouteWebsocket
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpWebsocket
//----------------------------------------------------------------------------------------------------------------------

func NewHttpWebsocket(app *fiber.Fiber,
	configService *lygo_http_server_config.HttpServerConfigHost,
	configRoutes []*lygo_http_server_config.HttpServerConfigRouteWebsocket) *HttpWebsocket {

	instance := new(HttpWebsocket)
	instance.app = app
	instance.configService = configService
	instance.configRoutes = configRoutes

	return instance
}

func (instance *HttpWebsocket) Init() {
	app := instance.app
	routes := instance.configRoutes
	configService := instance.configService
	if nil != configService.Websocket && configService.Websocket.Enabled {
		settings := configService.Websocket
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
		for _, route := range routes {
			if len(route.Path) > 0 {
				app.Get(route.Path, websocket.New(route.Handler, config))
			}
		}
	}
}
