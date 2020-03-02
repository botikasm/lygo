package lygo_http_server

import "encoding/json"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServerConfig struct {
	// HOST
	Address     string `json:"addr"`
	AddressTLS  string `json:"addr_tls"`
	VHost       bool   `json:"vhost"`
	StatEnabled bool   `json:"stat_enabled"`

	// TLS
	SslCert string `json:"ssl_cert"`
	SslKey  string `json:"ssl_key"`

	// FILE SERVER
	FileServerRoot    string   `json:"fileserver_root"`
	FileServerEnabled bool     `json:"fileserver_enabled"`
	IndexNames        []string `json:"index_names"`
	AcceptByteRange   bool     `json:"accept_byte_range"`
	Compress          bool     `json:"compress"`
}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerConfig
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServerConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *HttpServerConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}
