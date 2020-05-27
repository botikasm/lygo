package lygo_n_client

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/base/lygo_regex"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type NClient struct {
	Settings *NClientSettings

	//-- private --//
	initialized bool
	appToken    string
	events      *lygo_events.Emitter
	nio         *lygo_nio.NioClient
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNClient(settings *NClientSettings) *NClient {
	instance := new(NClient)
	instance.Settings = settings //lygo_http_server.NewHttpServer(&settings.HttpServerConfig)

	if nil == instance.Settings {
		instance.Settings = new(NClientSettings)
		instance.Settings.Nio = new(lygo_nio.NioSettings)
	}

	instance.events = lygo_events.NewEmitter()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NClient) IsOpen() bool {
	if nil != instance {
		if nil != instance.nio {
			return instance.nio.IsOpen()
		}
	}
	return false
}

func (instance *NClient) Start() ([]error, []string) {
	if nil != instance {
		if !instance.initialized {
			instance.initialized = true
			return instance.init()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}, []string{}
}

func (instance *NClient) Stop() []error {
	if nil != instance {
		if instance.initialized {
			instance.initialized = false
			return instance.close()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *NClient) SendData(data map[string]interface{}) ([]byte, error) {
	if nil != instance {
		if instance.IsOpen() {
			if nil == data["app_token"] {
				data["app_token"] = instance.getAppToken()
			}
			response, err := instance.nio.Send(data)
			if nil != err {
				return nil, err
			}
			return lygo_conv.ToArrayOfByte(response.Body), nil
		}
	}
	return nil, nil
}

func (instance *NClient) SendCommand(command *lygo_n_commons.Command) ([]byte, error) {
	if nil != instance {
		if instance.IsOpen() {
			response, err := instance.nio.Send(command)
			if nil != err {
				return nil, err
			}
			return lygo_conv.ToArrayOfByte(response.Body), nil
		}
	}
	return nil, nil
}

func (instance *NClient) Send(commandName string, params map[string]interface{}) ([]byte, error) {
	if nil != instance {
		tokens := strings.Split(commandName, ".")
		command := &lygo_n_commons.Command{
			AppToken:  instance.getAppToken(),
			Namespace: lygo_array.GetAt(tokens, 0, "").(string),
			Function:  lygo_array.GetAt(tokens, 1, "").(string),
			Params:    params,
		}
		return instance.SendCommand(command)
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NClient) init() ([]error, []string) {

	responseErrs := make([]error, 0)
	responseWarns := make([]string, 0)

	if instance.Settings.Enabled {
		// nio
		nioHost := instance.Settings.Nio.Host()
		nioPort := instance.Settings.Nio.Port()
		if len(nioHost) > 0 && nioPort > 0 {
			instance.nio = lygo_nio.NewNioClient(nioHost, nioPort)
			err := instance.nio.Open()
			if nil != err {
				responseErrs = append(responseErrs, err)
			}
		}
	} else {
		responseWarns = append(responseWarns, lygo_n_commons.ClientNotEnabledWarning.Error())
	}

	return responseErrs, responseWarns
}

func (instance *NClient) close() []error {
	response := make([]error, 0)

	// nio
	err := instance.nio.Close()
	if nil != err {
		response = append(response, err)
	}

	return response
}

func (instance *NClient) getAppToken() string {
	if len(instance.appToken) == 0 {
		if instance.IsOpen() {
			command := &lygo_n_commons.Command{
				AppToken:  "undefined",
				Namespace: "",
				Function:  "",
				Params:    nil,
			}
			command.SetName(lygo_n_server.CmdAppToken)
			response, err := instance.SendCommand(command)
			if nil == err {
				s := lygo_conv.ToString(response)
				if len(s) > 0 {
					if !lygo_regex.IsValidJsonObject(s) {
						instance.appToken = s
					}
				}
			}
		}
	}
	return instance.appToken
}
