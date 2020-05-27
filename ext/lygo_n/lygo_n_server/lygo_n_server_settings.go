package lygo_n_server

import (
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NServerSettings struct {
	Enabled bool                                      `json:"enabled"`
	Http    *lygo_http_server_config.HttpServerConfig `json:"http"`
	Nio     *lygo_nio.NioSettings                     `json:"nio"`
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NServerSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NServerSettings) String() string {
	return lygo_json.Stringify(instance)
}
