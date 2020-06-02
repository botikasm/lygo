package lygo_n_host

import (
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NHostSettings struct {
	Enabled bool                                      `json:"enabled"`
	Nio     *lygo_nio.NioSettings                     `json:"nio"`
	Http    *lygo_http_server_config.HttpServerConfig `json:"http"`
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NHostSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NHostSettings) String() string {
	return lygo_json.Stringify(instance)
}
