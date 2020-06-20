package lygo_license_client

import (
	"errors"
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_license"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_config"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_struct"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type LicenseClient struct {
	Config *lygo_license_config.LicenseConfig

	mux sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLicenseClient(config *lygo_license_config.LicenseConfig) *LicenseClient {
	instance := new(LicenseClient)
	instance.Config = config

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseClient) GetUrl() string {
	protocol := "http"
	if instance.Config.UseSSL {
		protocol = "https"
	}
	host := instance.Config.Host
	port := instance.Config.Port
	return lygo_strings.Format("%s://%s:%s/", protocol, host, port)
}

func (instance *LicenseClient) RequestLicense(path string) (license *lygo_license_struct.License, err error) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	license = new(lygo_license_struct.License)
	bytes := make([]byte, 0)
	if len(path) > 0 {
		bytes, err = instance.download(path)
	} else if len(instance.Config.Path) > 0 {
		bytes, err = instance.download(instance.Config.Path)
	}

	if nil == err && len(bytes) > 0 {
		if string(bytes[0]) != "{" {
			bytes, err = lygo_crypto.DecryptBytesAES(bytes, []byte(lygo_license.KEY))
		}
		data := string(bytes)
		if len(data) > 0 {
			license.Parse(data)
		}
	}

	return license, err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseClient) download(path string) ([]byte, error) {
	if len(path) > 0 {
		url := instance.GetUrl() + path

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    15 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(url)
		if nil == err {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if nil == err {
				return body, nil
			} else {
				return []byte{}, err
			}
		} else {
			return []byte{}, err
		}

		/*
			http := lygo_http_client.NewHttpClient()
			code, data, err := http.GetTimeout(url, 15*time.Second)
			if code != 200 {
				return []byte{}, errors.New(lygo_strings.Format("http_error: %s", code))
			}
			if nil == err {
				return data, nil
			} else {
				return []byte{}, err
			}*/
	}
	return []byte{}, errors.New("missing_path: path parameter is empty string")
}
