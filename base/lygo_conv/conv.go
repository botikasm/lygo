package lygo_conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func ToArray(val interface{}) []interface{} {
	if nil == val {
		return nil
	}
	_, aa := toArray(val)
	return aa
}

func ToArrayOfString(val interface{}) []string {
	if nil == val {
		return nil
	}
	_, aa := toArrayOfString(val)
	return aa
}

func ToString(val interface{}) string {
	if nil == val {
		return ""
	}
	// string
	s, ss := val.(string)
	if ss {
		return s
	}
	// integer
	i, ii := val.(int)
	if ii {
		return strconv.Itoa(i)
	}
	// float32
	f, ff := val.(float32)
	if ff {
		return fmt.Sprintf("%g", f) // Exponent as needed, necessary digits only
	}
	// float 64
	F, FF := val.(float64)
	if FF {
		return fmt.Sprintf("%g", F) // Exponent as needed, necessary digits only
		// return strconv.FormatFloat(F, 'E', -1, 64)
	}

	// boolean
	b, bb := val.(bool)
	if bb {
		return strconv.FormatBool(b)
	}

	// array
	if aa, tt := IsArray(val); aa {
		response := []string{}
		// array := make([]interface{}, tt.Len())
		for i := 0; i < tt.Len(); i++ {
			v := tt.Index(i).Interface()
			s := ToString(v)
			response = append(response, s)
		}
		return "[" + strings.Join(response, ",") + "]"
	}

	// map
	if b, _ := IsMap(val); b {
		data, err := json.Marshal(val)
		if nil == err {
			return string(data)
		}
	}

	// undefined value
	return fmt.Sprintf("%v", val)
}

func ToInt(val interface{}) int {
	return ToIntDef(val, -1)
}

func ToIntDef(val interface{}, def int) int {
	s := ToString(val)
	v, err := strconv.Atoi(s)
	if nil == err {
		return v
	}
	return def
}

func ToMap(val interface{}) map[string]interface{} {
	if b, _ := IsString(val); b {
		s := ToString(val)
		var m map[string]interface{}
		err := json.Unmarshal([]byte(s), &m)
		if nil == err {
			return m
		}
	}
	if b, _ := IsMap(val); b {
		return toMap(val)
	}
	return nil
}

func IsString(val interface{}) (bool, string) {
	v, vv := val.(string)
	if vv {
		return true, v
	}
	return false, ""
}

func IsInt(val interface{}) (bool, int) {
	v, vv := val.(int)
	if vv {
		return true, v
	}
	return false, 0
}

func IsBool(val interface{}) (bool, bool) {
	v, vv := val.(bool)
	if vv {
		return true, v
	}
	return false, false
}

func IsFloat32(val interface{}) (bool, float32) {
	v, vv := val.(float32)
	if vv {
		return true, v
	}
	return false, 0
}

func IsFloat64(val interface{}) (bool, float64) {
	v, vv := val.(float64)
	if vv {
		return true, v
	}
	return false, 0
}

func IsArray(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Slice:
		return true, rt
	case reflect.Array:
		return true, rt
	default:
		return false, rt
	}
}

func IsMap(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Map:
		return true, rt
	default:
		return false, rt
	}
}

func IsStruct(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Struct:
		return true, rt
	default:
		return false, rt
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toArray(val interface{}) (bool, []interface{}) {
	response := []interface{}{}
	aa, tt := IsArray(val)
	if aa {
		for i := 0; i < tt.Len(); i++ {
			v := tt.Index(i).Interface()
			response = append(response, v)
		}
	}
	return aa, response
}

func toArrayOfString(val interface{}) (bool, []string) {
	response := []string{}
	aa, tt := IsArray(val)
	if aa {
		for i := 0; i < tt.Len(); i++ {
			v := tt.Index(i).Interface()
			response = append(response, ToString(v))
		}
	}
	return aa, response
}

func toMap(val interface{}) map[string]interface{} {
	if m, b := val.(map[string]interface{}); b {
		return m
	}

	// warning: this change the pointer to original object
	data, err := json.Marshal(val)
	if nil == err {
		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		if nil == err {
			return m
		}
	}

	return nil
}
