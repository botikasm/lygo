package lygo_scripting

import (
	"github.com/botikasm/lygo/ext/lygo_scripting/lygo_scripting_tools"
	"github.com/dop251/goja"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScripEngine struct {
	Root string
	Name string

	//-- private --//
	runtime *goja.Runtime
	tools   []lygo_scripting_tools.ScriptingTool
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func New() *ScripEngine {
	response := new(ScripEngine)

	// init runtime
	response.runtime = goja.New()
	response.initTools()

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (engine *ScripEngine) RunString(program string) (goja.Value, error) {
	v, err := engine.runtime.RunString(program)
	return v, err
}

func (engine *ScripEngine) Set(name string, value interface{}) {
	engine.runtime.Set(name, value)
}

func (engine *ScripEngine) Get(name string) goja.Value {
	return engine.runtime.Get(name)
}

func (engine *ScripEngine) ToValue(value interface{}) goja.Value {
	return engine.runtime.ToValue(value)
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n t e x t
//----------------------------------------------------------------------------------------------------------------------

func (engine *ScripEngine) SetToolsContext(value interface{}) {
	for _, tool := range engine.tools {
		tool.SetContext(value)
	}
}

func (engine *ScripEngine) SetToolContext(toolName string, value interface{}) {
	v := engine.runtime.Get(toolName)
	if nil != v {
		tool := v.Export().(lygo_scripting_tools.ScriptingTool)
		if nil != tool {
			tool.SetContext(value)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (engine *ScripEngine) initTools() {

	params := new(lygo_scripting_tools.ScriptingToolParams)
	params.Runtime = engine.runtime
	params.Root = &engine.Root
	params.Name = &engine.Name

	Tconsole := lygo_scripting_tools.NewToolConsole(params)
	engine.runtime.Set(lygo_scripting_tools.TOOL_CONSOLE, Tconsole)
	engine.tools = append(engine.tools, Tconsole)

	Tstrings := lygo_scripting_tools.NewToolStrings(params)
	engine.runtime.Set(lygo_scripting_tools.TOOL_STRINGS, Tstrings)
	engine.tools = append(engine.tools, Tstrings)

	Tarrays := lygo_scripting_tools.NewToolArrays(params)
	engine.runtime.Set(lygo_scripting_tools.TOOL_ARRAYS, Tarrays)
	engine.tools = append(engine.tools, Tarrays)

	Tregexps := lygo_scripting_tools.NewToolRegExps(params)
	engine.runtime.Set(lygo_scripting_tools.TOOL_REGEXPS, Tregexps)
	engine.tools = append(engine.tools, Tregexps)

	Tconvert := lygo_scripting_tools.NewToolConvert(params)
	engine.runtime.Set(lygo_scripting_tools.TOOL_CONVERT, Tconvert)
	engine.tools = append(engine.tools, Tconvert)

}
