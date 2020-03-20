package lygo_scripting

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/dop251/goja"
)

func GetArgsStringStringString(context interface{}, args []goja.Value) (string, string, string) {
	arg1 := ""
	arg2 := ""
	arg3 := ""

	if len(args) > 0 {
		switch len(args) {
		case 1:
			arg1 = lygo_conv.ToString(args[0].Export())
			// fallback on context for latest arg
			if nil != context {
				arg2 = lygo_conv.ToString(context)
			}
		case 2:
			arg1 = lygo_conv.ToString(args[0].Export())
			arg2 = lygo_conv.ToString(args[1].Export())
			// fallback on context for latest arg
			if nil != context {
				arg3 = lygo_conv.ToString(context)
			}
		case 3:
			arg1 = lygo_conv.ToString(args[0].Export())
			arg2 = lygo_conv.ToString(args[1].Export())
			arg3 = lygo_conv.ToString(args[2].Export())
		}

	}

	return arg1, arg2, arg3
}
