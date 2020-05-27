package lygo_n

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_fmt"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
	"time"
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
//		m e s s a g e    h a n d l e r s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) RegisterCommand(command string, handler lygo_n_server.CommandHandler) {
	if nil != instance {
		if nil != instance.server {
			instance.server.RegisterCommand(command, handler)
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		m e s s a g e    s e n d e r
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) Send(commandName string, params map[string]interface{}) ([]byte, error) {
	if nil != instance {
		if instance.client.IsOpen() {
			return instance.client.Send(commandName, params)
		}
		return nil, nil
	}
	return nil, lygo_n_commons.PanicSystemError
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) open() []error {
	if !instance.initialized {
		instance.initialized = true
		response := make([]error, 0)

		// workspace
		if len(instance.Settings.Workspace) > 0 {
			lygo_paths.SetWorkspacePath(instance.Settings.Workspace)
		} else {
			lygo_paths.SetWorkspacePath("./_workspace")
		}

		// log level
		lygo_logs.SetLevel(lygo_logs.LEVEL_WARN)
		if len(instance.Settings.LogLevel) > 0 {
			lygo_logs.SetLevelName(instance.Settings.LogLevel)
		}
		lygo_logs.SetOutput(lygo_logs.OUTPUT_FILE)

		// client
		instance.server = lygo_n_server.NewNServer(instance.Settings.Server)
		errs, warnings := instance.server.Start()
		response = append(response, errs...)
		logWarns(warnings)

		// client
		instance.client = lygo_n_client.NewNClient(instance.Settings.Client)
		errs, warnings = instance.client.Start()
		response = append(response, errs...)
		logWarns(warnings)

		// check configuration
		if len(warnings)==0{
			warnings = instance.checkConfiguration()
			logWarns(warnings)
		}

		return response
	}
	return nil
}

func (instance *N) close() []error {
	if instance.initialized {
		instance.initialized = false
		response := make([]error, 0)
		if nil != instance.client {
			response = append(response, instance.client.Stop()...)
		}
		if nil != instance.server {
			response = append(response, instance.server.Stop()...)
		}
		return response
	}
	return nil
}

func (instance *N) checkConfiguration() []string {
	response := make([]string, 0)
	clientHost := instance.client.Settings.Nio.Host()
	clientPort := instance.client.Settings.Nio.Port()
	if clientHost == "localhost" || clientHost == "127.0.0.1" {
		serverPort := instance.server.Settings.Nio.Port()
		if serverPort == clientPort {
			// client is connecting to itself
			response = append(response, "Client auto-connect to itself at address "+instance.client.Settings.Nio.Address+
				"\n\tEnsure this is only for testing purpose.")
		}
	}
	return response
}

// ---------------------------------------------------------------------------------------------------------------------
//		S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func logWarns(warnings []string) {
	if len(warnings) > 0 {
		for _, w := range warnings {
			fmt.Println("WARN", lygo_fmt.FormatDate(time.Now(), "yyyy/MM/dd HH:mm"), w)
			lygo_logs.Warn("", w)
		}
	}
}
