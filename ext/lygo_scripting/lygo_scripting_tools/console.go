package lygo_scripting_tools

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/dop251/goja"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_CONSOLE = "$console"

const (
	LEVEL_INFO  = "info"
	LEVEL_ERROR = "error"
	LEVEL_DEBUG = "debug"
	LEVEL_WARN  = "warning"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolConsole struct {
	context interface{}
	params  *ScriptingToolParams
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolConsole(params *ScriptingToolParams) *ScriptingToolConsole {
	result := new(ScriptingToolConsole)
	result.params = params

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolConsole) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

//
func (tool *ScriptingToolConsole) Log(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		array := tool.getArgsArray(LEVEL_INFO, args)
		tool.write(tool.createLogRow(array))
	}
	return tool.params.Runtime.ToValue(nil)
}

func (tool *ScriptingToolConsole) Info(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		array := tool.getArgsArray(LEVEL_INFO, args)
		tool.write(tool.createLogRow(array))
	}
	return tool.params.Runtime.ToValue(nil)
}

func (tool *ScriptingToolConsole) Error(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		array := tool.getArgsArray(LEVEL_ERROR, args)
		tool.write(tool.createLogRow(array))
	}
	return tool.params.Runtime.ToValue(nil)
}

func (tool *ScriptingToolConsole) Debug(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		array := tool.getArgsArray(LEVEL_DEBUG, args)
		tool.write(tool.createLogRow(array))
	}
	return tool.params.Runtime.ToValue(nil)
}

func (tool *ScriptingToolConsole) Warn(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		array := tool.getArgsArray(LEVEL_WARN, args)
		tool.write(tool.createLogRow(array))
	}
	return tool.params.Runtime.ToValue(nil)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolConsole) getArgsArray(level string, args []goja.Value) []interface{} {
	var response []interface{}

	if len(args) > 0 {
		response = append(response, level)
		for _, val := range args {
			response = append(response, val.Export())
		}
	}

	return response
}

func (tool *ScriptingToolConsole) getFilename() string {
	params := tool.params
	root := "./console/"
	if len(*params.Root) > 0 {
		root = *params.Root
	}
	filename := "./scrip_log_" + lygo_rnd.Uuid() + ".log"
	if len(*params.Name) > 0 {
		filename = *params.Name + ".log"
	}

	// set output filename
	return lygo_paths.Concat(lygo_paths.WorkspacePath(root), filename)
}

func (tool *ScriptingToolConsole) createLogRow(args []interface{}) string {
	response := lygo_strings.Format("[%s]", time.Now())
	for i, val := range args {
		s := lygo_conv.ToString(val)
		if i == 0 {
			response += " - "
			s = strings.ToUpper(s)
		} else {
			response += ", "
		}
		response += s
	}

	return response + "\n"
}

func (tool *ScriptingToolConsole) write(row string) {
	filename := tool.getFilename()
	if len(filename) > 0 {
		lygo_paths.Mkdir(filename)
		lygo_io.AppendTextToFile(row, filename)
	}
}
