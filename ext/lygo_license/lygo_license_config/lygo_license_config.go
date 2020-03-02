package lygo_license_config

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_paths"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type LicenseConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Path   string `json:"path"`
	UseSSL bool   `json:"use_ssl"`
}

//----------------------------------------------------------------------------------------------------------------------
//	LicenseConfig
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *LicenseConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

// ----
// Return file name if 'Path' is a static file, otherwise returns empty string.
// Path can be both a static file ("/licenses/lic.json")
// or url invoking a remote endpoint ("/endpoint?uid=123&act=get")
// ----
func (instance *LicenseConfig) GetRequestFileName() string {
	filename := lygo_paths.FileName(instance.Path, true)
	if len(filename) > 0 {
		if len(lygo_paths.ExtensionName(filename)) > 0 {
			return filename
		}
	}
	return ""
}
