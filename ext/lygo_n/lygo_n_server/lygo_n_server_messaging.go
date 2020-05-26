package lygo_n_server

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// 		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type CommandHandler func(message *lygo_n_commons.Command) interface{}
type MessageFallbackHandler func(commandName string, message *lygo_n_commons.Message) (interface{}, bool)

type MessagingController struct {
	app             *NServer
	events          *lygo_events.Emitter
	fallback        MessageFallbackHandler // fallback for custom/unregisterd commands (namespace.function)
	commandHandlers map[string]CommandHandler
}

// ---------------------------------------------------------------------------------------------------------------------
// 		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewMessagingController(app *NServer, callback MessageFallbackHandler) *MessagingController {
	instance := new(MessagingController)
	instance.app = app
	instance.events = app.events
	instance.commandHandlers = make(map[string]CommandHandler)
	instance.fallback = callback

	registerInternalCommands(instance)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
// 		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MessagingController) Register(commandName string, handler CommandHandler) {
	if nil != instance {
		instance.commandHandlers[commandName] = handler
	}
}

func (instance *MessagingController) RegisterNS(namespace, function string, handler CommandHandler) {
	if nil != instance {
		commandName := namespace + "." + function
		instance.commandHandlers[commandName] = handler
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MessagingController) handleWsMessage(payload interface{}) {
	if nil != payload {
		if v, b := payload.(*lygo_http_server.HttpWebsocketEventPayload); b {
			instance.handleWs(v)
		}
	}
}

func (instance *MessagingController) handleNioMessage(nioMessage *lygo_nio.NioMessage) interface{} {
	// convert message body to string format
	body := lygo_conv.ToString(nioMessage.Body)
	if strings.Index(body, "{") > -1 {
		message := new(lygo_n_commons.Message)
		lygo_json.Read(body, &message.Payload)
		return instance.execute(message)
	}
	return lygo_n_commons.UnsupportedMessageTypeError // custom response
}

func (instance *MessagingController) handleWs(payload *lygo_http_server.HttpWebsocketEventPayload) {
	if nil != instance && nil != instance.events {
		ws := payload.Websocket
		if nil != ws && ws.IsAlive() && len(payload.Message.Data) > 0 {
			var m lygo_n_commons.Message
			err := json.Unmarshal(payload.Message.Data, &m)
			if nil == err && m.IsValid() {
				if m.IsValid() {
					response := instance.execute(&m)
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

func (instance *MessagingController) execute(message *lygo_n_commons.Message) interface{} {
	if nil != instance {
		commandName := message.Payload.Namespace + "." + message.Payload.Function
		if commandName != CmdAppToken {
			// check token
			if lygo_n_commons.AppToken != message.Payload.AppToken {
				return lygo_n_commons.InvalidTokenError
			}
		}
		handler := instance.getHandler(commandName)
		if nil != handler {
			return handler(message.Payload)
		} else if nil != instance.fallback {
			value, handled := instance.fallback(commandName, message)
			if handled {
				return value
			}
		}
		return lygo_n_commons.CommandNotFoundError
	}
	return lygo_n_commons.PanicSystemError
}

func (instance *MessagingController) getHandler(command string) CommandHandler {
	if nil != instance {
		if v, b := instance.commandHandlers[command]; b {
			return v
		}
	}
	return nil
}
