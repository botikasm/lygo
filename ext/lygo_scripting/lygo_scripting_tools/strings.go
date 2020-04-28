package lygo_scripting_tools

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/dop251/goja"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_STRINGS = "$strings"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolStrings struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolStrings(params *ScriptingToolParams) *ScriptingToolStrings {
	result := new(ScriptingToolStrings)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolStrings) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolStrings) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Split a string by sep.
// Support multiple separators
// @param sep string
// @param text string (Optional) CONTEXT is used if not found
// @return []string
func (tool *ScriptingToolStrings) Split(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		sep, text := tool.getArgsStringString(args)
		if len(sep) > 0 && len(text) > 0 {
			tokens := lygo_strings.Split(text, sep)
			return tool.runtime.ToValue(tokens)
		}
	}

	return tool.runtime.ToValue("")
}

// Get a substring
// @param start int Start index
// @param end int End index
// @param text string (Optional) CONTEXT is used if not found
// @return []string
func (tool *ScriptingToolStrings) Sub(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		start, end, text := tool.getArgsIntIntString(args)
		if len(text) > 0 {
			value := lygo_strings.Sub(text, start, end)
			return tool.runtime.ToValue(value)
		}
	}

	return tool.runtime.ToValue("")
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m p o u n d
//----------------------------------------------------------------------------------------------------------------------

// Split a string by spaces AND get a word at index
// @param index int
// @param text string (Optional) CONTEXT is used if not found
// @return bool
func (tool *ScriptingToolStrings) SplitBySpaceWordAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		index, text := tool.getArgsIntString(args)
		if index > -1 && len(text) > 0 {
			tokens := lygo_strings.Split(text, " \n")
			if len(tokens) > index {
				result := tokens[index]
				return tool.runtime.ToValue(result)
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Split a string by "sep" AND get a word at index
// @param index int
// @param sep string
// @param text string (Optional) CONTEXT is used if not found
// @return string
func (tool *ScriptingToolStrings) SplitWordAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		index, sep, text := tool.getArgsIntStringString(args)
		if index > -1 && len(text) > 0 {
			tokens := strings.Split(text, sep)
			if len(tokens) > index {
				result := tokens[index]
				return tool.runtime.ToValue(result)
			}
		}
	}

	return tool.runtime.ToValue("")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolStrings) getArgsStringString(args []goja.Value) (string, string) {
	var arg1, argCtx string
	arg1 = lygo_conv.ToString(args[0].Export())
	if len(arg1) > 0 {
		if len(args) == 2 {
			argCtx = lygo_conv.ToString(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToString(tool.context)
		}
	}

	return arg1, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntString(args []goja.Value) (int, string) {
	var arg1 int
	var argCtx string

	arg1 = lygo_conv.ToInt(args[0].Export())
	if arg1 > -1 {
		if len(args) == 2 {
			argCtx = lygo_conv.ToString(args[1].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToString(tool.context)
		}
	}

	return arg1, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntStringString(args []goja.Value) (int, string, string) {
	var arg1 int
	var arg2 string
	var argCtx string

	arg1 = lygo_conv.ToInt(args[0].Export())
	if arg1 > -1 {

		if len(args) > 1 {
			arg2 = lygo_conv.ToString(args[1].Export())
		}

		if len(args) == 3 {
			argCtx = lygo_conv.ToString(args[2].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToString(tool.context)
		}
	}

	return arg1, arg2, argCtx
}

func (tool *ScriptingToolStrings) getArgsIntIntString(args []goja.Value) (int, int, string) {
	var arg1 int
	var arg2 int
	var argCtx string

	arg1 = lygo_conv.ToInt(args[0].Export())
	if arg1 > -1 {

		if len(args) > 1 {
			arg2 = lygo_conv.ToInt(args[1].Export())
		}

		if len(args) == 3 {
			argCtx = lygo_conv.ToString(args[2].Export())
		}
	}

	// fallback on context for latest arg
	if len(argCtx) == 0 {
		if nil != tool.context {
			argCtx = lygo_conv.ToString(tool.context)
		}
	}

	return arg1, arg2, argCtx
}
