package lygo_http_server

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"testing"
)

func TestBasic(t *testing.T) {

	// load configuration
	config := config()
	if nil == config {
		t.Errorf("Configuration is not valid")
	}

	server := NewHttpServer(config)
	server.Config.FileServerEnabled = true
	server.CallbackError = onError

	// start server
	err := server.Start()
	if nil != err {
		t.Error(err)
	}

	// Wait forever.
	server.Join()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func config() *HttpServerConfig {
	text_cfg, _ := lygo_io.ReadTextFromFile("./lygo_http_server_config.json")
	config := new(HttpServerConfig)
	config.Parse(text_cfg)

	return config
}

func onError(errCtx *HttpServerError) {
	fmt.Println(errCtx.Message, errCtx.Error.Error())
}
