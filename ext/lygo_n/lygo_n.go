package lygo_n

import (
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type N struct {
	Settings *NSettings

	//-- private --//
	initialized bool
	client      *lygo_n_client.NClient
	server      *lygo_n_server.NServer
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewNode(settings *NSettings) *N {
	instance := new(N)
	instance.Settings = settings

	if nil == instance.Settings {
		instance.Settings = new(NSettings)
	}

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) Start() []error {
	if nil != instance {
		return instance.open()
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *N) Stop() []error {
	if nil != instance {
		return instance.close()
	}
	return []error{lygo_n_commons.PanicSystemError}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) open() []error {
	if !instance.initialized {
		instance.initialized = true

	}
	return nil
}

func (instance *N) close() []error {
	if instance.initialized {
		instance.initialized = false

	}
	return nil
}
