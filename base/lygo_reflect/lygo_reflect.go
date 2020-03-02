package lygo_reflect

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"reflect"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Get(object interface{}, name string) interface{} {
	if b, _ := lygo_conv.IsMap(object); b {
		m := lygo_conv.ToMap(object)
		if nil != m {
			return m[name]
		}
	} else {
		v := reflect.ValueOf(object)
		e := v.Elem()
		if e.Kind() == reflect.Struct {
			f := e.FieldByName(name)
			if f.IsValid() {
				return f.Interface()
			}
		}
	}
	return nil
}

func Set(object interface{}, name string, value interface{}) (interface{}, bool) {
	if b, _ := lygo_conv.IsMap(object); b {
		// TEST DIFFERENT MAPS
		if m, b := object.(map[string]interface{}); b {
			m[name] = value
			return m, true
		}
		if m, b := object.(map[string]string); b {
			m[name] = lygo_conv.ToString(value)
			return m, true
		}
		if m, b := object.(map[string]int); b {
			m[name] = lygo_conv.ToInt(value)
			return m, true
		}
		if m, b := object.(map[string][]interface{}); b {
			m[name] = lygo_conv.ToArray(value)
			return m, true
		}

		// fallback ( WARN: unmarshal changes the object referenced )
		m := lygo_conv.ToMap(object)
		if nil != m {
			m[name] = value
			return m, true
		}
	} else {
		e := reflect.ValueOf(object).Elem()
		if e.Kind() == reflect.Struct {
			f := e.FieldByName(name)
			if f.IsValid() {
				if f.CanSet() {
					f.Set(reflect.ValueOf(value))
					return object, true
				}
			}
		}
	}
	return nil, false
}

func GetString(object interface{}, name string) string {
	v := Get(object, name)
	if nil != v {
		return lygo_conv.ToString(v)
	}
	return ""
}

func GetInt(object interface{}, name string) int {
	v := Get(object, name)
	if nil != v {
		return lygo_conv.ToInt(v)
	}
	return 0
}

func GetArray(object interface{}, name string) []interface{} {
	v := Get(object, name)
	if nil != v {
		return lygo_conv.ToArray(v)
	}
	return []interface{}{}
}

func GetArrayOfString(object interface{}, name string) []string {
	v := Get(object, name)
	if nil != v {
		return lygo_conv.ToArrayOfString(v)
	}
	return []string{}
}
