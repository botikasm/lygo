package lygo_scripting_tools

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/dop251/goja"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_PATHS = "$paths"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolPaths struct {
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolPaths(params *ScriptingToolParams) *ScriptingToolPaths {
	result := new(ScriptingToolPaths)
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolPaths) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Return absolute path
// @param [string] path
// @return string
func (tool *ScriptingToolPaths) GetAbsolutePath(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		path := tool.getArgsString(args)
		if len(path) > 0 {
			absolutePath := lygo_paths.Absolute(path)
			return tool.runtime.ToValue(absolutePath)
		}
	}
	return tool.runtime.ToValue([]map[string]string{})
}

func (tool *ScriptingToolPaths) GetWorkspacePath(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		path := tool.getArgsString(args)
		if len(path) > 0 {
			absolutePath := lygo_paths.WorkspacePath(path)
			return tool.runtime.ToValue(absolutePath)
		}
	}
	return tool.runtime.ToValue([]map[string]string{})
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolPaths) getArgsString(args []goja.Value) string {
	arg1 := ""

	if len(args) > 0 {
		switch len(args) {
		case 1:
			arg1 = lygo_conv.ToString(args[0])
		default:
			arg1 = lygo_conv.ToString(args[0])
		}
	}

	return arg1
}
