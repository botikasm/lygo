package lygo_http_server

import (
	"encoding/json"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServerConfig struct {
	// HOST
	Address    string `json:"addr"`
	AddressTLS string `json:"addr_tls"`

	// TLS
	SslCert string `json:"ssl_cert"`
	SslKey  string `json:"ssl_key"`

	// SETTINGS
	// Enables use of theSO_REUSEPORTsocket option. This will spawn multiple Go processes listening on the same port. learn more about socket sharding
	Prefork bool `json:"prefork"` // default: false
	// Enables the Server HTTP header with the given value.
	ServerHeader string `json:"server_header"` // default: ""
	// When enabled, the router treats /foo and /foo/ as different. Otherwise, the router treats /foo and /foo/ as the same.
	StrictRouting bool `json:"strict_routing"` // default: false
	// When enabled, /Foo and /foo are different routes. When disabled, /Fooand /foo are treated the same.
	CaseSensitive bool `json:"strict_routing"` // default: false
	// When enabled, all values returned by context methods are immutable. By default they are valid until you return from the handler, see issue #185.
	Immutable bool `json:"immutable"` // default: false
	// Sets the maximum allowed size for a request body, if the size exceeds the configured limit, it sends 413 - Request Entity Too Large response.
	BodyLimit int `json:"body_limit"` // default: 4 * 1024 * 1024
	// The amount of time allowed to read the full request including body. Default timeout is unlimited.
	ReadTimeout time.Duration `json:"read_timeout"` // default: 0
	// The maximum duration before timing out writes of the response. Default timeout is unlimited.
	WriteTimeout time.Duration `json:"write_timeout"` // default: 0
	// The maximum amount of time to wait for the next request when keep-alive is enabled. If IdleTimeout is zero, the value of ReadTimeout is used.
	IdleTimeout time.Duration `json:"idle_timeout"` // default: 0

	// STATIC
	Static []*HttpServerConfigStatic `json:"static"`

	// CORS
	CORS *HttpServerConfigCORS `json:"cors"`

	// Compression
	Compression *HttpServerConfigCompression `json:"compression"`

	// Limiter
	Limiter *HttpServerConfigLimiter `json:"limiter"`
}
type HttpServerConfigStatic struct {
	Enabled   bool   `json:"enabled"`
	Prefix    string `json:"prefix"`
	Root      string `json:"root"`
	Index     string `json:"index"`
	Compress  bool   `json:"compress"`
	ByteRange bool   `json:"byte_range"`
	Browse    bool   `json:"browse"`
}

type HttpServerConfigCORS struct {
	Enabled bool `json:"enabled"`
	// AllowOrigin defines a list of origins that may access the resource.
	AllowOrigins []string `json:"allow_origins"` // default: []string{"*"}
	// AllowMethods defines a list methods allowed when accessing the resource. This is used in response to a preflight request.
	AllowMethods []string `json:"allow_methods"` // default: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
	// AllowCredentials indicates whether or not the response to the request can be exposed when the credentials flag is true. When used as part of a response to a preflight request, this indicates whether or not the actual request can be made using credentials.
	AllowCredentials bool `json:"allow_credentials"` // default: false
	// ExposeHeaders defines a whitelist headers that clients are allowed to access.
	ExposeHeaders []string `json:"expose_headers"` // default: nil
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	MaxAge int `json:"max_age"` // default: 0
}

type HttpServerConfigCompression struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level"` // Level of compression, 0, 1, 2, 3, 4
}

type HttpServerConfigLimiter struct {
	Enabled bool `json:"enabled"`
	// Timeout in seconds on how long to keep records of requests in memory
	Timeout int `json:"timeout"` // default: 68
	// Max number of recent connections during Timeout seconds before sending a 429
	Max int `json:"timeout"` // default: 10
	// Response body
	Message string `json:"message"` // default: "Too many requests, please try again later."
	// Response status code
	StatusCode int `json:"status_code"` // default: 429

}

//----------------------------------------------------------------------------------------------------------------------
//	HttpServerConfig
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServerConfig() *HttpServerConfig {
	instance := new(HttpServerConfig)
	instance.init()
	return instance
}

func (instance *HttpServerConfig) Parse(text string) error {
	err := json.Unmarshal([]byte(text), &instance)
	instance.init()
	return err
}

func (instance *HttpServerConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *HttpServerConfig) init() {
	if nil == instance.Static {
		instance.Static = make([]*HttpServerConfigStatic, 0)
	}
	if nil == instance.CORS {
		instance.CORS = new(HttpServerConfigCORS)
	}
	if nil == instance.Compression {
		instance.Compression = new(HttpServerConfigCompression)
	}
	if nil == instance.Limiter {
		instance.Limiter = new(HttpServerConfigLimiter)
	}
}
