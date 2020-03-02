package _tests

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/ext/lygo_logs"
)

func InitContext() {

	var root = lygo_paths.Absolute("../")
	fmt.Println("WORKSPACE:", root)

	lygo_paths.SetWorkspaceParent(root)

	lygo_logs.SetLevel(lygo_logs.LEVEL_TRACE)
	lygo_logs.SetOutput(lygo_logs.OUTPUT_FILE)
	lygo_logs.Info(
		lygo_logs.LogContext{
			Caller: "initializer.InitContext()",
			Data:   "This is the INITIALIZATION context. LOG Level: " + lygo_conv.ToString(lygo_logs.GetLevel()),
		},
		"initialized")

}
