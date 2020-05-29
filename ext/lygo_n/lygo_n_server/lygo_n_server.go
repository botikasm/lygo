package lygo_n_server

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_types"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/gofiber/fiber"
	"io"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type NServer struct {
	Settings  *NServerSettings
	OnMessage MessageFallbackHandler // handle all messages (http, nio)

	//-- private --//
	initialized  bool
	statusBuffer bytes.Buffer
	messaging    *MessagingController
	events       *lygo_events.Emitter
	server       *lygo_http_server.HttpServer
	nio          *lygo_nio.NioServer
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNServer(settings *NServerSettings) *NServer {
	instance := new(NServer)
	instance.initialized = false
	instance.Settings = settings //lygo_http_server.NewHttpServer(&settings.HttpServerConfig)

	if nil == instance.Settings {
		instance.Settings = new(NServerSettings)
		instance.Settings.Enabled = false
		instance.Settings.Http = lygo_http_server_config.NewHttpServerConfig()
		instance.Settings.Nio = new(lygo_nio.NioSettings)
	}

	instance.events = lygo_events.NewEmitter()
	instance.messaging = NewMessagingController(instance, instance.onMessage)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NServer) IsOpen() bool {
	if nil != instance {
		return nil != instance.nio && instance.nio.IsOpen()
	}
	return false
}

func (instance *NServer) GetUUID() string {
	if nil != instance {
		return instance.nio.GetUUID()
	}
	return ""
}

func (instance *NServer) GetStatus() string {
	if nil != instance {
		return instance.statusBuffer.String()
	}
	return ""
}

func (instance *NServer) WriteStatus(w io.Writer) (int64, error) {
	if nil != instance {
		return instance.statusBuffer.WriteTo(w)
	}
	return 0, nil
}

func (instance *NServer) Start() ([]error, []string) {
	if nil != instance {
		if !instance.initialized {
			return instance.init()
		}
		return []error{}, []string{}
	}
	return []error{lygo_n_commons.PanicSystemError}, []string{}
}

func (instance *NServer) Join() []error {
	if nil != instance {
		err, _ := instance.Start()
		if nil != err {
			return err
		}
		instance.server.Join()
	}
	return nil
}

func (instance *NServer) Stop() []error {
	if nil != instance {
		if nil != instance.server {
			instance.server.Stop()
		}
		return nil
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *NServer) RegisterCommand(command string, handler CommandHandler) {
	if nil != instance {
		if nil != instance.messaging {
			instance.messaging.Register(command, handler)
		}
	}
}

func (instance *NServer) AddCommandNS(namespace, function string, handler CommandHandler) {
	if nil != instance {
		if nil != instance.messaging {
			instance.messaging.RegisterNS(namespace, function, handler)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *NServer) Server() *lygo_http_server.HttpServer {
	if nil != instance {
		if nil == instance.server {
			if nil != instance.Settings.Http {
				instance.server = lygo_http_server.NewHttpServer(instance.Settings.Http)
			}
		}
		return instance.server
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NServer) init() ([]error, []string) {
	instance.initialized = true

	responseErrs := make([]error, 0)
	responseWarns := make([]string, 0)

	if instance.Settings.Enabled {
		// http
		if nil != instance.Server() {
			if instance.server.IsEnabled() {
				instance.server.CallbackError = instance.onError
				instance.server.CallbackLimitReached = instance.onLimit
				// enable websocket
				websocket := instance.server.Config.WebsocketEndpoint
				if len(websocket) > 0 {
					instance.server.Websocket(func(ws *lygo_http_server.HttpWebsocketConn) {
						ws.OnMessage(instance.handleWs)
					})
				}
				err := instance.server.Start()
				if nil == err {
					for _, host := range instance.Settings.Http.Hosts {
						instance.statusBuffer.WriteString(fmt.Sprintln("HTTP SERVER LISTENING AT:", host.Address))
					}
				} else {
					instance.statusBuffer.WriteString(fmt.Sprintln("HTTP SERVER ERROR:", err))
					responseErrs = append(responseErrs, err)
				}
				if len(websocket) > 0 {
					instance.statusBuffer.WriteString(fmt.Sprintln("WEBSOCKET RESPONDING AT:", instance.server.Config.WebsocketEndpoint))
				}
			} else {
				instance.statusBuffer.WriteString(fmt.Sprintln("WEB-SERVER IS NOT ENABLED"))
			}
		}

		// nio
		if nil != instance.Settings.Nio && len(instance.Settings.Nio.Address) > 0 && instance.Settings.Nio.Port() > 0 {
			instance.nio = lygo_nio.NewNioServer(instance.Settings.Nio.Port())
			err := instance.nio.Open()
			if nil == err {
				instance.nio.OnMessage(instance.messaging.handleNioMessage)
				instance.statusBuffer.WriteString(fmt.Sprintln("NIO SERVER LISTENING ON PORT:", instance.Settings.Nio.Port()))
			} else {
				instance.statusBuffer.WriteString(fmt.Sprintln("NIO SERVER ERROR:", err))
				responseErrs = append(responseErrs, err)
			}
		}
	} else {
		responseWarns = append(responseWarns, lygo_n_commons.ServerNotEnabledWarning.Error())
	}

	return responseErrs, responseWarns
}

func (instance *NServer) onError(errCtx *lygo_http_server_types.HttpServerError) {
	if nil != instance && nil != instance.server && nil != instance.events {
		// fmt.Println(errCtx.Message, errCtx.Error.Error())
		lygo_logs.Error(errCtx.Message, errCtx.Error.Error())
		instance.events.Emit(lygo_n_commons.EventError, lygo_n_commons.ContextWebsocket, errCtx, errCtx.Error)
	}
}

func (instance *NServer) onLimit(c *fiber.Ctx) {
	if nil != instance && nil != instance.server && nil != instance.events {
		c.Send("too many requests: limit exceeded")
		instance.events.Emit(lygo_n_commons.EventError, lygo_n_commons.ContextWebsocket, "too many requests: limit exceeded", c.Error().Error())
	}
}

func (instance *NServer) handleWs(payload *lygo_http_server.HttpWebsocketEventPayload) {
	if nil != instance && nil != instance.server && nil != instance.messaging {
		instance.messaging.handleWsMessage(payload)
	}
}

func (instance *NServer) onMessage(method string, message *lygo_n_commons.Message) (interface{}, bool) {
	if nil != instance.OnMessage {
		return instance.OnMessage(method, message)
	}
	return nil, false
}
