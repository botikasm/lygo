package lygo_scripting_tools

import "github.com/dop251/goja"

type ScriptingTool interface {
	SetContext(context interface{})
}

type ScriptingToolParams struct {
	Root    *string // pointer to external string (the engine Root)
	Name    *string // pointer to external string (the engine Name)
	Runtime *goja.Runtime
}
