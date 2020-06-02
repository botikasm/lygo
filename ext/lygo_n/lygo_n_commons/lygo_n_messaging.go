package lygo_n_commons

import (
	"encoding/json"
	"errors"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_json"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// 		Message
// ---------------------------------------------------------------------------------------------------------------------

type Message struct {
	RequestUUID string    `json:"request_uuid"`
	Lang        string    `json:"lang"`
	UID         string    `json:"uid"` // user_id
	Payload     *Command  `json:"payload"`
	Response    *Response `json:"response"`
}

func (instance *Message) Marshal() []byte {
	data, _ := json.Marshal(instance)
	return data
}

func (instance *Message) String() string {
	return string(instance.Marshal())
}

func (instance *Message) IsValid() bool {
	if nil != instance {
		if len(instance.UID) > 0 && nil != instance.Payload {
			return true
		}
	}
	return false
}

func (instance *Message) SetResponse(val interface{}) {
	if nil != instance {
		if v, b := val.(error); b {
			instance.Response = &Response{
				Error: v.Error(),
				Data:  nil,
			}
		} else {
			instance.Response = &Response{
				Data: lygo_conv.ToArrayOfByte(val),
			}
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 		Command
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	AppToken  string      `json:"app_token"`
	Namespace string      `json:"namespace"`
	Function  string      `json:"function"`
	Params    interface{} `json:"params"`
}

func (instance *Command) Name() string {
	if nil != instance {
		return instance.Namespace + "." + instance.Function
	}
	return ""
}

func (instance *Command) SetName(name string) {
	tokens := strings.Split(name, ".")
	instance.Namespace = lygo_array.GetAt(tokens, 0, "").(string)
	instance.Function = lygo_array.GetAt(tokens, 1, "").(string)
}

func (instance *Command) GetParam(name string) interface{} {
	if nil != instance && nil != instance.Params {
		params := lygo_conv.ToMap(instance.Params)
		if v, b := params[name]; b {
			return v
		}
	}
	return nil
}

func (instance *Command) GetParamAsString(name string) string {
	if nil != instance && nil != instance.Params {
		return lygo_conv.ToString(instance.GetParam(name))
	}
	return ""
}

func (instance *Command) GetParamAsInt(name string) int {
	if nil != instance && nil != instance.Params {
		return lygo_conv.ToInt(instance.GetParam(name))
	}
	return 0
}

func (instance *Command) GetParamAsMap(name string) map[string]interface{} {
	if nil != instance && nil != instance.Params {
		return lygo_conv.ForceMap(instance.GetParam(name))
	}
	return nil
}

func (instance *Command) GetParamAsMapArray(name string) []map[string]interface{} {
	if nil != instance && nil != instance.Params {
		response := make([]map[string]interface{}, 0)
		value := instance.GetParam(name)
		if nil != value {
			arr := lygo_conv.ToArray(value)
			if nil != arr {
				for _, v := range arr {
					m := lygo_conv.ForceMap(v)
					if nil != m {
						response = append(response, m)
					}
				}
			}
		}
		return response
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// 		NServerInfo
// ---------------------------------------------------------------------------------------------------------------------

type NHostInfo struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	Timestamp int64  `json:"timestamp"`
}

// ---------------------------------------------------------------------------------------------------------------------
// 		Response
// ---------------------------------------------------------------------------------------------------------------------

type Response struct {
	Info  *NHostInfo `json:"info"`
	Error string     `json:"error"`
	Data  []byte     `json:"data"`
}

func (instance *Response) HasError() bool {
	return nil != instance && len(instance.Error) > 0
}

func (instance *Response) GetError() error {
	if nil != instance && len(instance.Error) > 0 {
		return errors.New(instance.Error)
	}
	return nil
}

func (instance *Response) GetDataAsString() string {
	if nil != instance && len(instance.Error) == 0 {
		return string(instance.Data)
	}
	return ""
}

func (instance *Response) GetDataAsMap() map[string]interface{} {
	if nil != instance && len(instance.Error) == 0 {
		var response map[string]interface{}
		err := lygo_json.Read(instance.Data, &response)
		if nil == err {
			return response
		}
	}
	return nil
}

func (instance *Response) GetDataAsArrayOfMap() []map[string]interface{} {
	if nil != instance && len(instance.Error) == 0 {
		var response []map[string]interface{}
		err := lygo_json.Read(instance.Data, &response)
		if nil == err {
			return response
		}
	}
	return nil
}
