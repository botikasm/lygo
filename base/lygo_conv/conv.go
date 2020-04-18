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

func ToArray(val ...interface{}) []interface{} {
	if nil == val {
		return nil
	}
	aa := toArray(val...)
	return aa
}

func ToArrayOfString(val ...interface{}) []string {
	if nil == val {
		return nil
	}
	aa := toArrayOfString(val...)
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
		// byte array??
		if ba, b :=val.([]byte);b{
			return string(ba)
		} else {
			response := []string{}
			// array := make([]interface{}, tt.Len())
			for i := 0; i < tt.Len(); i++ {
				v := tt.Index(i).Interface()
				s := ToString(v)
				response = append(response, s)
			}
			return "[" + strings.Join(response, ",") + "]"
		}
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

func Int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func ToInt(val interface{}) int {
	return ToIntDef(val, -1)
}

func ToIntDef(val interface{}, def int) int {
	if b, s := IsString(val); b {
		v, err := strconv.Atoi(s)
		if nil == err {
			return v
		}
	}
	switch i := val.(type) {
	case float32:
		return int(i)
	case float64:
		return int(i)
	case int:
		return i
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	}

	return def
}

func ToFloat32(val interface{}) float32 {
	return ToFloat32Def(val, -1)
}

func ToFloat32Def(val interface{}, defVal float32) float32 {
	if b, s := IsString(val); b {
		v, err := strconv.ParseFloat(s, 32)
		if nil == err {
			return float32(v)
		}
	}
	switch i := val.(type) {
	case float32:
		return i
	case float64:
		return float32(i)
	case int:
		return float32(i)
	case int8:
		return float32(i)
	case int16:
		return float32(i)
	case int32:
		return float32(i)
	case int64:
		return float32(i)
	}
	return defVal
}

func ToFloat64(val interface{}) float64 {
	return ToFloat64Def(val, -1.0)
}

func ToFloat64Def(val interface{}, defVal float64) float64 {
	if b, s := IsString(val); b {
		v, err := strconv.ParseFloat(s, 64)
		if nil == err {
			return v
		}
	}
	switch i := val.(type) {
	case float32:
		return float64(i)
	case float64:
		return i
	case int:
		return float64(i)
	case int8:
		return float64(i)
	case int16:
		return float64(i)
	case int32:
		return float64(i)
	case int64:
		return float64(i)
	}
	return defVal
}

func ToBool(val interface{}) bool {
	if b, s := IsString(val); b {
		v, err := strconv.ParseBool(s)
		if nil == err {
			return v
		}
	}
	v, vv := val.(bool)
	if vv {
		return v
	}
	return false
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

func ToMapOfString(val interface{}) map[string]string {
	if b, _ := IsString(val); b {
		s := ToString(val)
		var m map[string]string
		err := json.Unmarshal([]byte(s), &m)
		if nil == err {
			return m
		}
	}
	if b, _ := IsMap(val); b {
		return toMapOfString(val)
	}
	return nil
}

func ForceMap(val interface{}) map[string]interface{} {
	m := ToMap(val)
	if nil == m {
		return toMap(val)
	}
	return m
}

func ForceMapOfString(val interface{}) map[string]string {
	m := ToMapOfString(val)
	if nil == m {
		return toMapOfString(val)
	}
	return m
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

func Equals(val1, val2 interface{}) bool {
	if val1 == val2 {
		return true
	}

	if b, v := IsString(val1); b {
		return v == ToString(val2)
	}
	if b, v := IsInt(val1); b {
		return v == ToInt(val2)
	}
	if b, v := IsFloat32(val1); b {
		return v == ToFloat32(val2)
	}
	if b, v := IsFloat64(val1); b {
		return v == ToFloat64(val2)
	}
	if b, v := IsBool(val1); b {
		return v == ToBool(val2)
	}

	return false
}

func NotEquals(val1, val2 interface{}) bool {
	if val1 == val2 {
		return true
	}

	if b, v := IsString(val1); b {
		return v != ToString(val2)
	}
	if b, v := IsInt(val1); b {
		return v != ToInt(val2)
	}
	if b, v := IsFloat32(val1); b {
		return v != ToFloat32(val2)
	}
	if b, v := IsFloat64(val1); b {
		return v != ToFloat64(val2)
	}
	if b, v := IsBool(val1); b {
		return v != ToBool(val2)
	}

	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toArray(args ...interface{}) []interface{} {
	response := []interface{}{}
	for _, val := range args {
		aa, tt := IsArray(val)
		if aa {
			for i := 0; i < tt.Len(); i++ {
				v := tt.Index(i).Interface()
				response = append(response, v)
			}
		}else {
			response = append(response, ToString(val))
		}
	}
	return response
}

func toArrayOfString(args ...interface{}) []string {
	response := []string{}
	for _, val := range args {
		b, tt := IsArray(val)
		if b {
			for i := 0; i < tt.Len(); i++ {
				v := tt.Index(i).Interface()
				response = append(response, ToString(v))
			}
		} else {
			response = append(response, ToString(val))
		}
	}

	return response
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

func toMapOfString(val interface{}) map[string]string {
	if m, b := val.(map[string]string); b {
		return m
	}

	// warning: this change the pointer to original object
	data, err := json.Marshal(val)
	if nil == err {
		var m map[string]string
		err = json.Unmarshal(data, &m)
		if nil == err {
			return m
		}
	}

	return nil
}
