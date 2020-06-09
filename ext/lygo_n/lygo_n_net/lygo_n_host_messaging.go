package lygo_n_net

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// 		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type CommandHandler func(message *lygo_n_commons.Command) interface{}
type MessageFallbackHandler func(commandName string, message *lygo_n_commons.Message) (interface{}, bool)

type MessagingController struct {
	app             *NHost
	events          *lygo_events.Emitter
	fallback        MessageFallbackHandler // fallback for custom/unregisterd commands (namespace.function)
	commandHandlers map[string]CommandHandler
}

// ---------------------------------------------------------------------------------------------------------------------
// 		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewMessagingController(app *NHost, callback MessageFallbackHandler) *MessagingController {
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

func (instance *MessagingController) Execute(commandName string, params map[string]interface{}) *lygo_n_commons.Response {
	if nil != instance {
		tokens := strings.Split(commandName, ".")
		command := &lygo_n_commons.Command{
			AppToken:  lygo_n_commons.AppToken,
			Namespace: lygo_array.GetAt(tokens, 0, "").(string),
			Function:  lygo_array.GetAt(tokens, 1, "").(string),
			Params:    params,
		}
		message := new(lygo_n_commons.Message)
		message.Payload = command

		// execute
		instance.execute(message)

		return message.Response
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// 		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MessagingController) handleNioMessage(nioMessage *lygo_nio.NioMessage) interface{} {
	message := new(lygo_n_commons.Message)
	// convert message body to string format
	body := lygo_conv.ToString(nioMessage.Body)
	if strings.Index(body, "{") > -1 {
		err := lygo_json.Read(body, &message.Payload)
		if nil != err {
			message.Response = &lygo_n_commons.Response{
				Info:  instance.app.Info,
				Error: err.Error(),
				Data:  nil,
			}
		} else {
			instance.execute(message)
		}
		// return instance.execute(message)
	} else {
		message.SetResponse(lygo_n_commons.UnsupportedMessageTypeError)
	}
	//return lygo_n_commons.UnsupportedMessageTypeError // custom response
	return message
}

func (instance *MessagingController) execute(message *lygo_n_commons.Message) {
	if nil != instance {
		commandName := message.Payload.Namespace + "." + message.Payload.Function

		// check if command is token protected and validate
		if commandName != CmdAppToken {
			// check token
			if lygo_n_commons.AppToken != message.Payload.AppToken {
				message.SetResponse(lygo_n_commons.InvalidTokenError)
				// return lygo_n_commons.InvalidTokenError
				goto done
			}
		}

		handler := instance.getHandler(commandName)
		if nil != handler {
			value := handler(message.Payload)
			message.SetResponse(value)
			// return value
			goto done
		} else if nil != instance.fallback {
			value, handled := instance.fallback(commandName, message)
			if handled {
				message.SetResponse(value)
				// return value
				goto done
			}
		}
		message.SetResponse(lygo_n_commons.CommandNotFoundError)
		// return lygo_n_commons.CommandNotFoundError
		goto done
	}
	message.SetResponse(lygo_n_commons.PanicSystemError)
	// return lygo_n_commons.PanicSystemError

done:
	if nil != message.Response {
		message.Response.Info = instance.app.Info
	}
}

func (instance *MessagingController) getHandler(command string) CommandHandler {
	if nil != instance {
		if v, b := instance.commandHandlers[command]; b {
			return v
		}
	}
	return nil
}
