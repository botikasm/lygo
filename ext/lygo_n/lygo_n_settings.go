package lygo_n

import (
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NSettings struct {
	Client *lygo_n_client.NClientSettings
	Server lygo_n_server.NServerSettings
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NSettings) String() string {
	return lygo_json.Stringify(instance)
}
