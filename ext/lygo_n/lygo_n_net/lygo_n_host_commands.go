package lygo_n_net

import "github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"

var (
	// no token required
	CmdAppToken = "n.sys_app_token"

	// token required
	CmdVersion  = "n.sys_version"
	CmdPing = "n.sys_ping"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func registerInternalCommands(controller *MessagingController) {

	controller.Register(CmdVersion, getVersion)
	controller.Register(CmdAppToken, getAppToken)

	controller.Register(CmdPing, doPing)

}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getVersion(message *lygo_n_commons.Command) interface{} {
	return lygo_n_commons.Version
}

func getAppToken(message *lygo_n_commons.Command) interface{} {
	return lygo_n_commons.AppToken
}

func doPing(message *lygo_n_commons.Command) interface{} {
	return true
}
