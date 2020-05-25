package lygo_n_commons

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// 		m e s s a g e
// ---------------------------------------------------------------------------------------------------------------------

type Message struct {
	RequestUUID string           `json:"request_uuid"`
	Lang        string           `json:"lang"`
	UID         string           `json:"uid"` // user_id
	Payload     *Command         `json:"payload"`
	Response    *MessageResponse `json:"response"`
}

type MessageResponse struct {
	Error string        `json:"error"`
	Data  []interface{} `json:"data"`
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
			instance.Response = &MessageResponse{
				Error: v.Error(),
				Data:  nil,
			}
		} else {
			instance.Response = &MessageResponse{
				Data: lygo_conv.ToArray(val),
			}
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 		c o m m a n d
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
