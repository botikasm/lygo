package lygo_scripting_utils

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/dop251/goja"
)

func GetArgsString(context interface{}, args []goja.Value) string {
	arg1 := ""

	switch len(args) {
	case 1:
		arg1 = lygo_conv.ToString(args[0].Export())
	default:
		if nil != context {
			arg1 = lygo_conv.ToString(context)
		}
	}
	return arg1
}

func GetArgsStringString(context interface{}, args []goja.Value) (string, string) {
	arg1 := ""
	arg2 := ""

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
	default:
		if nil != context {
			arg1 = lygo_conv.ToString(context)
		}
	}

	return arg1, arg2
}

func GetArgsStringStringString(context interface{}, args []goja.Value) (string, string, string) {
	arg1 := ""
	arg2 := ""
	arg3 := ""

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
	default:
		if nil != context {
			arg1 = lygo_conv.ToString(context)
		}
	}

	return arg1, arg2, arg3
}
