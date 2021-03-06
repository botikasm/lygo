package lygo_n_commons

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_nio"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server/lygo_http_server_config"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	NSettings
//----------------------------------------------------------------------------------------------------------------------

type NSettings struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
	LogLevel  string `json:"log_level"` // warn, info, error, debug

	Discovery *NDiscoverySettings `json:"discovery"`
	Server    *NHostSettings      `json:"server"`
}

func (instance *NSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NSettings) String() string {
	return lygo_json.Stringify(instance)
}

//----------------------------------------------------------------------------------------------------------------------
//	NAddress
//----------------------------------------------------------------------------------------------------------------------

type NAddress string

func NewNAddress(text string) NAddress {
	return NAddress(text)
}

func (instance *NAddress) String() string {
	return string(*instance)
}

func (instance *NAddress) Host() string {
	return lygo_array.GetAt(instance.tokens(), 0, "").(string)
}

func (instance *NAddress) Port() int {
	return lygo_conv.ToIntDef(lygo_array.GetAt(instance.tokens(), 1, "").(string), 0)
}

func (instance *NAddress) tokens() []string {
	return strings.Split(instance.String(), ":")
}

//----------------------------------------------------------------------------------------------------------------------
//	NDiscoverySettings
//----------------------------------------------------------------------------------------------------------------------

type NDiscoverySettings struct {
	Publishers []NAddress                   `json:"publishers"` // target publishers
	NetworkId  string                       `json:"network_id"`
	Publisher  *NDiscoveryPublisherSettings `json:"publisher"`
	Publish    *NDiscoveryPublishSettings   `json:"publish"`
}

func (instance *NDiscoverySettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NDiscoverySettings) String() string {
	return lygo_json.Stringify(instance)
}

//----------------------------------------------------------------------------------------------------------------------
//	NDiscoveryPublisherSettings
//----------------------------------------------------------------------------------------------------------------------

type NDiscoveryPublisherSettings struct {
	Enabled bool `json:"enabled"`
}

//----------------------------------------------------------------------------------------------------------------------
//	NDiscoveryPublishSettings
//----------------------------------------------------------------------------------------------------------------------

type NDiscoveryPublishSettings struct {
	Enabled bool     `json:"enabled"`
	Address NAddress `json:"address"`
}

func (instance *NDiscoveryPublishSettings) IsAddress(address string) bool {
	if nil != instance {
		return len(address) > 0 && address == instance.Address.String()
	}
	return false
}

func (instance *NDiscoveryPublishSettings) HasAddress() bool {
	if nil != instance {
		return len(instance.Address) > 0
	}
	return false
}

func (instance *NDiscoveryPublishSettings) IsEnabled() bool {
	if nil != instance {
		return len(instance.Address) > 0 && instance.Enabled
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	NDiscoveryBroadcastSettings
//----------------------------------------------------------------------------------------------------------------------

type NDiscoveryNetworkSettings struct {
	Enabled bool `json:"enabled"`
}

//----------------------------------------------------------------------------------------------------------------------
//	NHostSettings
//----------------------------------------------------------------------------------------------------------------------

type NHostSettings struct {
	Enabled bool                                      `json:"enabled"`
	Nio     *lygo_nio.NioSettings                     `json:"nio"`
	Http    *lygo_http_server_config.HttpServerConfig `json:"http"`
}

func (instance *NHostSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NHostSettings) String() string {
	return lygo_json.Stringify(instance)
}

//----------------------------------------------------------------------------------------------------------------------
//	NConnSettings
//----------------------------------------------------------------------------------------------------------------------

type NConnSettings struct {
	Enabled bool                  `json:"enabled"`
	Nio     *lygo_nio.NioSettings `json:"nio"`
}

func (instance *NConnSettings) Parse(text string) error {
	return lygo_json.Read(text, &instance)
}

func (instance *NConnSettings) String() string {
	return lygo_json.Stringify(instance)
}