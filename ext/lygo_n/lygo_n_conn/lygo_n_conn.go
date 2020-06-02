package lygo_n_conn

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/base/lygo_regex"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_host"
	"io"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type NConn struct {
	Settings *NClientSettings

	//-- private --//
	initialized  bool
	appToken     string
	statusBuffer bytes.Buffer
	events       *lygo_events.Emitter
	nio          *lygo_nio.NioClient
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNConn(args ...interface{}) *NConn {
	instance := new(NConn)

	if len(args) > 0 {
		if len(args) == 1 {
			if v, b := args[0].(*NClientSettings); b {
				instance.Settings = v //lygo_http_server.NewHttpServer(&settings.HttpServerConfig)
			}
		} else if len(args) == 2 {
			if host, b := args[0].(string); b {
				if port, b := args[1].(int); b {
					instance.Settings = new(NClientSettings)
					instance.Settings.Nio = new(lygo_nio.NioSettings)
					instance.Settings.Enabled = true
					instance.Settings.Nio.Address = fmt.Sprintf("%v:%v", host, port)
				}
			}
		}
	}

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
func (instance *NConn) GetUUID() string {
	if nil != instance {
		return instance.nio.GetUUID()
	}
	return ""
}

func (instance *NConn) GetAddress() string {
	if nil != instance {
		return instance.Settings.Nio.Address
	}
	return ""
}

func (instance *NConn) GetStatus() string {
	if nil != instance {
		return instance.statusBuffer.String()
	}
	return ""
}

func (instance *NConn) WriteStatus(w io.Writer) (int64, error) {
	if nil != instance {
		return instance.statusBuffer.WriteTo(w)
	}
	return 0, nil
}

func (instance *NConn) SetEventManager(events *lygo_events.Emitter) bool {
	if nil != instance {
		if nil != events {
			instance.events = events
		}
	}
	return false
}

func (instance *NConn) IsOpen() bool {
	if nil != instance {
		if nil != instance.nio {
			return instance.nio.IsOpen()
		}
	}
	return false
}

func (instance *NConn) Start() ([]error, []string) {
	if nil != instance {
		if !instance.initialized {
			instance.initialized = true
			return instance.init()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}, []string{}
}

func (instance *NConn) Stop() []error {
	if nil != instance {
		if instance.initialized {
			instance.initialized = false
			return instance.close()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *NConn) SendData(data map[string]interface{}) *lygo_n_commons.Response {
	if nil != instance {
		if instance.IsOpen() {
			if nil == data["app_token"] {
				data["app_token"] = instance.getAppToken()
			}
			response, err := instance.nio.Send(data)
			if nil != err {
				return instance.unmarshalMessage(nil, err)
			}
			return instance.unmarshalMessage(response, nil)
		}
	}
	return nil
}

func (instance *NConn) SendCommand(command *lygo_n_commons.Command) *lygo_n_commons.Response {
	if nil != instance {
		if instance.IsOpen() {
			response, err := instance.nio.Send(command)
			if nil != err {
				return instance.unmarshalMessage(nil, err)
			}
			return instance.unmarshalMessage(response, nil)
		}
	}
	return nil
}

func (instance *NConn) Send(commandName string, params map[string]interface{}) *lygo_n_commons.Response {
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
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NConn) init() ([]error, []string) {

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

func (instance *NConn) close() []error {
	response := make([]error, 0)

	// nio
	err := instance.nio.Close()
	if nil != err {
		response = append(response, err)
	}

	return response
}

func (instance *NConn) getAppToken() string {
	if len(instance.appToken) == 0 {
		if instance.IsOpen() {
			command := &lygo_n_commons.Command{
				AppToken:  "undefined",
				Namespace: "",
				Function:  "",
				Params:    nil,
			}
			command.SetName(lygo_n_host.CmdAppToken)
			response := instance.SendCommand(command)
			if len(response.Error) == 0 {
				s := lygo_conv.ToString(response.Data)
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

func (instance *NConn) unmarshalMessage(nm *lygo_nio.NioMessage, err error) *lygo_n_commons.Response {
	if nil != err {
		return &lygo_n_commons.Response{
			Error: err.Error(),
			Data:  nil,
		}
	}
	if nil != nm {
		data := lygo_conv.ToArrayOfByte(nm.Body)
		var m lygo_n_commons.Message
		lygo_json.Read(data, &m)
		return m.Response
	}
	return nil
}
