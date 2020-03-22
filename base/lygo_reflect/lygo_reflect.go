package lygo_reflect

import (
	"database/sql"
	"errors"
	"github.com/botikasm/lygo/base/lygo_conv"
	"reflect"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Get(object interface{}, name string) interface{} {
	if m, b := object.(map[string]interface{}); b {
		if nil != m {
			return m[name]
		}
	} else if b, _ := lygo_conv.IsMap(object); b {
		m := lygo_conv.ToMap(object)
		if nil != m {
			return m[name]
		}
	} else {
		v := reflect.ValueOf(object)
		if v.IsValid() {
			return getFieldValue(v, name)
		}
	}
	return nil
}

func Set(object interface{}, name string, value interface{}) bool {
	if b, _ := lygo_conv.IsMap(object); b {
		b := setMapField(object, name, value)
		if b {
			return b
		}

		// fallback ( WARN: unmarshal changes the object referenced )
		mp := lygo_conv.ToMap(object)
		if nil != mp {
			mp[name] = value
			return true
		}
	} else {
		return setFieldValue(reflect.ValueOf(object), name, value)
	}
	return false
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

// Copy an object into another
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// check source
		if source.IsValid() {
			fromTypeFields := deepFields(fromType)
			//fmt.Printf("%#v", fromTypeFields)
			// Copy from field to field or method
			for _, field := range fromTypeFields {
				name := field.Name

				if fromField := source.FieldByName(name); fromField.IsValid() {
					// has field
					if toField := dest.FieldByName(name); toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField) {
								if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
									return err
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from method to field
			for _, field := range deepFields(toType) {
				name := field.Name

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0])
						}
					}
				}
			}
		}
		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getFieldValue(e reflect.Value, name string) interface{} {
	switch e.Kind() {
	case reflect.Ptr:
		elem := e.Elem()
		return getFieldValue(elem, name)
	case reflect.Struct:
		f := e.FieldByName(strings.Title(name))
		if f.IsValid() {
			return f.Interface()
		}
	case reflect.Map:
		m, _ := e.Interface().(map[string]interface{})
		return m[name]
	}
	return nil
}

func setFieldValue(e reflect.Value, name string, value interface{}) bool {
	switch e.Kind() {
	case reflect.Ptr:
		return setFieldValue(e.Elem(), name, value)
	case reflect.Struct:
		f := e.FieldByName(strings.Title(name))
		if f.IsValid() {
			if f.CanSet() {
				f.Set(reflect.ValueOf(value))
				return true
			}
		}
	case reflect.Map:
		return setMapField(e.Interface(), name, value)
	}
	return false
}

func setMapField(object interface{}, name string, value interface{}) bool {
	if m, b := object.(map[string]interface{}); b {
		m[name] = value
		return true
	}
	if m, b := object.(map[string]string); b {
		m[name] = lygo_conv.ToString(value)
		return true
	}
	if m, b := object.(map[string]int); b {
		m[name] = lygo_conv.ToInt(value)
		return true
	}
	if m, b := object.(map[string][]interface{}); b {
		m[name] = lygo_conv.ToArray(value)
		return true
	}
	return false
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
