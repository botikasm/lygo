package lygo_scripting_tools

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/dop251/goja"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_ARRAYS = "$arrays"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolArrays struct {
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolArrays(params *ScriptingToolParams) *ScriptingToolArrays {
	result := new(ScriptingToolArrays)
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolArrays) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Return a value contained in an array
// @param index
// @param array
// @return nil or value found at index
func (tool *ScriptingToolArrays) GetAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		arg1, argCtx := tool.getArgsIntArray(args)
		if arg1 > -1 && nil != argCtx {
			if len(argCtx) > arg1 {
				v := argCtx[arg1]
				return tool.runtime.ToValue(v)
			}
		}
	}
	return tool.runtime.ToValue("")
}

// Return first value contained in an array. If array is nil or empty returns an empty string
// @param array
// @return nil or value
func (tool *ScriptingToolArrays) GetFirst(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		argCtx := tool.getArgsArray(args)
		if nil != argCtx && len(argCtx) > 0 {
			v := argCtx[0]
			return tool.runtime.ToValue(v)

		}
	}
	return tool.runtime.ToValue("")
}

// Return last value contained in an array. If array is nil or empty returns an empty string
// @param array
// @return nil or value
func (tool *ScriptingToolArrays) GetLast(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		argCtx := tool.getArgsArray(args)
		if nil != argCtx && len(argCtx) > 0 {
			v := argCtx[len(argCtx)-1]
			return tool.runtime.ToValue(v)

		}
	}
	return tool.runtime.ToValue("")
}

// Return range values contained in an array. If array is nil or empty returns an empty array
// @param array
// @return nil or sub-array
func (tool *ScriptingToolArrays) GetSub(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	start, end, argCtx := tool.getArgsIntIntArray(args)
	if nil != argCtx && len(argCtx) > 0 {
		if start < 0 {
			start = 0
		}
		if end < 1 || end > len(argCtx) {
			end = len(argCtx) - 1
		}
		v := argCtx[start : end+1]
		return tool.runtime.ToValue(v)

	}

	return tool.runtime.ToValue([]interface{}{})
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolArrays) getArgsArray(args []goja.Value) []interface{} {
	var argCtx []interface{}

	if len(args) == 1 {
		argCtx = lygo_conv.ToArray(args[0].Export())
	}

	// fallback on context for latest arg
	if nil == argCtx || len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToArray(tool.context)
		}
	}

	return argCtx
}

func (tool *ScriptingToolArrays) getArgsIntArray(args []goja.Value) (int, []interface{}) {
	var arg1 int
	var argCtx []interface{}

	arg1 = lygo_conv.ToInt(args[0])
	if arg1 > -1 {
		if len(args) == 2 {
			argCtx = lygo_conv.ToArray(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if nil == argCtx || len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToArray(tool.context)
		}
	}

	return arg1, argCtx
}

func (tool *ScriptingToolArrays) getArgsIntIntArray(args []goja.Value) (int, int, []interface{}) {
	var arg1 int
	var arg2 int
	var argCtx []interface{}

	if len(args) > 1 {
		arg1 = lygo_conv.ToInt(args[0])
		arg2 = lygo_conv.ToInt(args[1])

		if len(args) == 3 {
			argCtx = lygo_conv.ToArray(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if nil == argCtx || len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToArray(tool.context)
		}
	}

	return arg1, arg2, argCtx
}
