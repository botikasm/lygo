package lygo_n_client

import (
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NClientSettings struct {
	Enabled bool                  `json:"enabled"`
	Nio     *lygo_nio.NioSettings `json:"nio"`
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NClientSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NClientSettings) String() string {
	return lygo_json.Stringify(instance)
}
