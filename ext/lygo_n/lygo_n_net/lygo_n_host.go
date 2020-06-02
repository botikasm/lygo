package lygo_n_net

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"io"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type NHost struct {
	Info      *lygo_n_commons.NHostInfo
	Settings  *lygo_n_commons.NHostSettings
	OnMessage MessageFallbackHandler // handle all messages (http, nio)

	//-- private --//
	initialized  bool
	statusBuffer bytes.Buffer
	messaging    *MessagingController
	events       *lygo_events.Emitter
	nio          *lygo_nio.NioServer
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNHost(settings *lygo_n_commons.NHostSettings) *NHost {
	instance := new(NHost)
	instance.Info = new(lygo_n_commons.NHostInfo)
	instance.Info.Name = "" // assigned later from outside
	instance.Info.UUID = "" // assigned from nio server
	instance.Info.Timestamp = time.Now().Unix()
	instance.initialized = false
	instance.Settings = settings //lygo_http_server.NewHttpServer(&settings.HttpServerConfig)

	if nil == instance.Settings {
		instance.Settings = new(lygo_n_commons.NHostSettings)
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

func (instance *NHost) IsOpen() bool {
	if nil != instance {
		return nil != instance.nio && instance.nio.IsOpen()
	}
	return false
}

func (instance *NHost) GetUUID() string {
	if nil != instance {
		return instance.nio.GetUUID()
	}
	return ""
}

func (instance *NHost) GetStatus() string {
	if nil != instance {
		return instance.statusBuffer.String()
	}
	return ""
}

func (instance *NHost) WriteStatus(w io.Writer) (int64, error) {
	if nil != instance {
		return instance.statusBuffer.WriteTo(w)
	}
	return 0, nil
}

func (instance *NHost) SetEventManager(events *lygo_events.Emitter) bool {
	if nil != instance {
		if nil != events {
			instance.events = events
		}
	}
	return false
}

func (instance *NHost) Start() ([]error, []string) {
	if nil != instance {
		if !instance.initialized {
			return instance.init()
		}
		return []error{}, []string{}
	}
	return []error{lygo_n_commons.PanicSystemError}, []string{}
}

func (instance *NHost) Join() []error {
	if nil != instance {
		err, _ := instance.Start()
		if nil != err {
			return err
		}
		instance.nio.Join()
	}
	return nil
}

func (instance *NHost) Stop() []error {
	if nil != instance {
		if nil != instance.nio {
			err := instance.nio.Close()
			if nil != err {
				return []error{err}
			}
		}
		return nil
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *NHost) RegisterCommand(command string, handler CommandHandler) {
	if nil != instance {
		if nil != instance.messaging {
			instance.messaging.Register(command, handler)
		}
	}
}

func (instance *NHost) AddCommandNS(namespace, function string, handler CommandHandler) {
	if nil != instance {
		if nil != instance.messaging {
			instance.messaging.RegisterNS(namespace, function, handler)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NHost) init() ([]error, []string) {
	instance.initialized = true

	responseErrs := make([]error, 0)
	responseWarns := make([]string, 0)

	if instance.Settings.Enabled {

		// nio
		if nil != instance.Settings.Nio && len(instance.Settings.Nio.Address) > 0 && instance.Settings.Nio.Port() > 0 {
			instance.nio = lygo_nio.NewNioServer(instance.Settings.Nio.Port())
			instance.Info.UUID = instance.nio.GetUUID()
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

func (instance *NHost) onMessage(method string, message *lygo_n_commons.Message) (interface{}, bool) {
	if nil != instance.OnMessage {
		return instance.OnMessage(method, message)
	}
	return nil, false
}
