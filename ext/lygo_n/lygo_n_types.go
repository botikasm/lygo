package lygo_n

import (
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

const CMD_GET_NODE_LIST = "n.sys_get_node_list"

//----------------------------------------------------------------------------------------------------------------------
//	NSettings
//----------------------------------------------------------------------------------------------------------------------

type NSettings struct {
	Workspace string `json:"workspace"`
	LogLevel  string `json:"log_level"` // warn, info, error, debug

	Discovery *NDiscoverySettings `json:"discovery"`

	Client *lygo_n_client.NClientSettings `json:"client"`
	Server *lygo_n_server.NServerSettings `json:"server"`
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
	Network    *NDiscoveryNetworkSettings   `json:"network"`
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

//----------------------------------------------------------------------------------------------------------------------
//	NDiscoveryBroadcastSettings
//----------------------------------------------------------------------------------------------------------------------

type NDiscoveryNetworkSettings struct {
	Enabled bool `json:"enabled"`
}
