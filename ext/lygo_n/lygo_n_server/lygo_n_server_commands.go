package lygo_n_server

import "github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"

var (
	CmdVersion  = "n.sys_version"
	CmdAppToken = "n.sys_app_token"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func registerInternalCommands(controller *MessagingController) {

	controller.Register(CmdVersion, getVersion)
	controller.Register(CmdAppToken, getAppToken)

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
