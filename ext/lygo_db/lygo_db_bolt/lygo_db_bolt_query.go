package lygo_db_bolt

import (
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
	Filter *BoltQueryFilter
}

type BoltQueryFilter struct {
	Conditions []*BoltQueryConditionGroup
}

type BoltQueryConditionGroup struct {
	Operator string
	Filters  []*BoltQueryCondition
}

type BoltQueryCondition struct {
	LeftField  string // absolute value or field "doc.name"
	Comparator string // ==, !=, > ...
	RightField string // absolute value or field "doc.surname", "Rossi"
}

//----------------------------------------------------------------------------------------------------------------------
//	BoltQuery
//----------------------------------------------------------------------------------------------------------------------

func (instance *BoltQuery) MatchFilter(entity interface{}) bool {
	response := false
	if nil != instance {
		if nil != instance.Filter && nil != entity {
			if nil != instance.Filter.Conditions {
				conditions := instance.Filter.Conditions
				if len(conditions) > 0 {
					// OPTIMISTIC CONDITION
					response = true
					for _, condition := range conditions {
						if nil != condition {
							operator := condition.Operator
							filters := condition.Filters
							for _, filter := range filters {
								f1 := getValue(entity, filter.LeftField)
								f2 := getValue(entity, filter.RightField)
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

	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getValue(entity interface{}, propertyOrValue string) interface{} {
	if strings.Index(propertyOrValue, ".") > -1 {
		// found a property
		field := ""
		tokens := lygo_strings.Split(propertyOrValue, ".")
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
	return propertyOrValue
}
