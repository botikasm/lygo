package lygo_db_bolt

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/base/lygo_strings"
	"strings"
)

const (
	ComparatorEqual        = "=="
	ComparatorNotEqual     = "!="
	ComparatorGreater      = ">"
	ComparatorLower        = "<"
	ComparatorLowerEqual   = "<="
	ComparatorGreaterEqual = ">="

	OperatorAnd = "&&"
	OperatorOr  = "||"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltQuery struct {
	Conditions []*BoltQueryConditionGroup `json:"conditions"`
}

type BoltQueryConditionGroup struct {
	Operator string                `json:"operator"`
	Filters  []*BoltQueryCondition `json:"filters"`
}

type BoltQueryCondition struct {
	Field      interface{} `json:"field"`      // absolute value or field "doc.name"
	Comparator string      `json:"comparator"` // ==, !=, > ...
	Value      interface{} `json:"value"`      // absolute value or field "doc.surname", "Rossi"
}

//----------------------------------------------------------------------------------------------------------------------
//	BoltQuery
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltQuery) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *BoltQuery) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *BoltQuery) MatchFilter(entity interface{}) bool {
	response := false
	if nil != instance {
		if nil != instance.Conditions && nil != entity {
			conditions := instance.Conditions
			if len(conditions) > 0 {
				// OPTIMISTIC CONDITION
				response = true
				for _, condition := range conditions {
					if nil != condition {
						operator := condition.Operator
						filters := condition.Filters
						for _, filter := range filters {
							f1 := getValue(entity, filter.Field)
							f2 := getValue(entity, filter.Value)
							cp := filter.Comparator

							match := false
							switch cp {
							case ComparatorEqual:
								match = lygo_conv.Equals(f1, f2)
							case ComparatorNotEqual:
								match = lygo_conv.NotEquals(f1, f2)
							default:
								match = false
							}

							switch operator {
							case OperatorAnd:
								// AND
								response = response && match
							case OperatorOr:
								// OR
								response = response || match
							default:
								// invalid operator
								response = false
							}
							if !response {
								break
							}
						} // filters
						if !response {
							break
						}
					}
				} // conditions
			}
		}

	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getValue(entity interface{}, propertyOrValue interface{}) interface{} {
	if b, property := lygo_conv.IsString(propertyOrValue); b {
		if strings.Index(property, ".") > -1 {
			// found a property
			field := ""
			tokens := lygo_strings.Split(property, ".")
			switch len(tokens) {
			case 1:
				field = tokens[0]
			default:
				field = tokens[1]
			}
			if len(field) > 0 {
				return lygo_reflect.Get(entity, field)
			}
		}
	}

	return propertyOrValue
}
