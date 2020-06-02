package lygo_n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
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

type NHttp struct {
	Settings           *lygo_n_commons.NHostSettings
	SendCommandHandler func(commandName string, params map[string]interface{}) *lygo_n_commons.Response

	//-- private --//
	initialized  bool
	statusBuffer bytes.Buffer
	events       *lygo_events.Emitter
	http         *lygo_http_server.HttpServer
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNHttp(settings *lygo_n_commons.NHostSettings) *NHttp {
	instance := new(NHttp)
	instance.initialized = false
	instance.Settings = settings //lygo_http_server.NewHttpServer(&settings.HttpServerConfig)

	if nil == instance.Settings {
		instance.Settings = new(lygo_n_commons.NHostSettings)
		instance.Settings.Enabled = false
		instance.Settings.Http = lygo_http_server_config.NewHttpServerConfig()
	}

	instance.events = lygo_events.NewEmitter()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NHttp) IsOpen() bool {
	if nil != instance {
		return nil != instance.http && instance.http.IsEnabled()
	}
	return false
}

func (instance *NHttp) GetStatus() string {
	if nil != instance {
		return instance.statusBuffer.String()
	}
	return ""
}

func (instance *NHttp) WriteStatus(w io.Writer) (int64, error) {
	if nil != instance {
		return instance.statusBuffer.WriteTo(w)
	}
	return 0, nil
}

func (instance *NHttp) SetEventManager(events *lygo_events.Emitter) bool {
	if nil != instance {
		if nil != events {
			instance.events = events
		}
	}
	return false
}

func (instance *NHttp) Start() ([]error, []string) {
	if nil != instance {
		if !instance.initialized {
			return instance.init()
		}
		return []error{}, []string{}
	}
	return []error{lygo_n_commons.PanicSystemError}, []string{}
}

func (instance *NHttp) Join() []error {
	if nil != instance {
		err, _ := instance.Start()
		if nil != err {
			return err
		}
		instance.http.Join()
	}
	return nil
}

func (instance *NHttp) Stop() []error {
	if nil != instance {
		if nil != instance.http {
			instance.http.Stop()
		}
		return nil
	}
	return []error{lygo_n_commons.PanicSystemError}
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *NHttp) Http() *lygo_http_server.HttpServer {
	if nil != instance {
		return instance.getHttp()
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NHttp) init() ([]error, []string) {
	instance.initialized = true

	responseErrs := make([]error, 0)
	responseWarns := make([]string, 0)

	if instance.Settings.Enabled {
		// http
		http := instance.getHttp()
		if nil != http {
			if http.IsEnabled() {
				http.CallbackError = instance.onError
				http.CallbackLimitReached = instance.onLimit
				// enable websocket
				websocket := http.Config.WebsocketEndpoint
				if len(websocket) > 0 {
					http.Websocket(func(ws *lygo_http_server.HttpWebsocketConn) {
						ws.OnMessage(instance.handleWs)
					})
				}
				err := http.Start()
				if nil == err {
					for _, host := range instance.Settings.Http.Hosts {
						instance.statusBuffer.WriteString(fmt.Sprintln("HTTP SERVER LISTENING AT:", host.Address))
					}
				} else {
					instance.statusBuffer.WriteString(fmt.Sprintln("HTTP SERVER ERROR:", err))
					responseErrs = append(responseErrs, err)
				}
				if len(websocket) > 0 {
					instance.statusBuffer.WriteString(fmt.Sprintln("WEBSOCKET RESPONDING AT:", instance.http.Config.WebsocketEndpoint))
				}
			} else {
				instance.statusBuffer.WriteString(fmt.Sprintln("WEB-SERVER IS NOT ENABLED"))
			}
		}
	} else {
		responseWarns = append(responseWarns, lygo_n_commons.ServerNotEnabledWarning.Error())
	}

	return responseErrs, responseWarns
}

func (instance *NHttp) getHttp() *lygo_http_server.HttpServer {
	if nil != instance {
		if nil == instance.http {
			if nil != instance.Settings.Http {
				instance.http = lygo_http_server.NewHttpServer(instance.Settings.Http)
			}
		}
		return instance.http
	}
	return nil
}

func (instance *NHttp) onError(errCtx *lygo_http_server_types.HttpServerError) {
	if nil != instance && nil != instance.http && nil != instance.events {
		// fmt.Println(errCtx.Message, errCtx.Error.Error())
		lygo_logs.Error(errCtx.Message, errCtx.Error.Error())
		instance.events.Emit(lygo_n_commons.EventError, lygo_n_commons.ContextWebsocket, errCtx, errCtx.Error)
	}
}

func (instance *NHttp) onLimit(c *fiber.Ctx) {
	if nil != instance && nil != instance.http && nil != instance.events {
		c.Send("too many requests: limit exceeded")
		instance.events.Emit(lygo_n_commons.EventError, lygo_n_commons.ContextWebsocket, "too many requests: limit exceeded", c.Error().Error())
	}
}

func (instance *NHttp) handleWs(payload *lygo_http_server.HttpWebsocketEventPayload) {
	if nil != instance && nil != instance.http && nil != instance.SendCommandHandler {
		if nil != payload {
			ws := payload.Websocket
			if nil != ws && ws.IsAlive() && len(payload.Message.Data) > 0 {
				var m lygo_n_commons.Message
				err := json.Unmarshal(payload.Message.Data, &m)
				if nil == err && m.IsValid() {
					if m.IsValid() {
						// response := instance.execute(&m)
						// m.SetResponse(response)
						commandName := m.Payload.Name()
						params := m.Payload.Params.(map[string]interface{})
						response := instance.SendCommandHandler(commandName, params)
						m.SetResponse(response)
						if ws.IsAlive() {
							ws.SendData(m.Marshal())
						}
					} else {
						// invalid message
					}
				} else {
					instance.events.Emit(lygo_n_commons.EventError, lygo_n_commons.ContextWebsocket, payload, err)
				}
			}
		}
	}
}
