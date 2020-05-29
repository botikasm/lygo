package lygo_n

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_fmt"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_sys"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_client"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_server"
	"io"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type N struct {
	Settings *NSettings

	//-- private --//
	uuid        string
	statusBuff  bytes.Buffer
	initialized bool
	client      *lygo_n_client.NClient
	server      *lygo_n_server.NServer
	discovery   *NDiscovery
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewNode(settings *NSettings) *N {
	instance := new(N)
	instance.Settings = settings

	if nil == instance.Settings {
		instance.Settings = new(NSettings)
		instance.Settings.Discovery = new(NDiscoverySettings)
		instance.Settings.Discovery.Publish = new(NDiscoveryPublishSettings)
		instance.Settings.Discovery.Publish.Enabled = false
		instance.Settings.Discovery.Publisher = new(NDiscoveryPublisherSettings)
		instance.Settings.Discovery.Publisher.Enabled = false
		instance.Settings.Discovery.Network = new(NDiscoveryNetworkSettings)
		instance.Settings.Discovery.Network.Enabled = false
	}

	sysid, err := lygo_sys.ID()
	if nil != err {
		sysid = lygo_rnd.Uuid()
	}
	instance.uuid = fmt.Sprintf("%v", sysid)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) GetUUID() string {
	if nil != instance {
		return instance.uuid
	}
	return ""
}

func (instance *N) GetStatus() string {
	if nil != instance {
		var buff bytes.Buffer

		if nil != instance.server {
			buff.WriteString("------------------------\n")
			buff.WriteString("\tSERVER\t")
			buff.WriteString(instance.server.GetUUID() + "\n")
			buff.WriteString(instance.indentLines(instance.server.GetStatus()))
		}

		if nil != instance.client && instance.client.IsOpen() {
			buff.WriteString("------------------------\n")
			buff.WriteString("\tCLIENT\t")
			buff.WriteString(instance.client.GetUUID() + "\n")
			buff.WriteString(instance.indentLines(instance.client.GetStatus()))
			buff.WriteString(instance.indentLines(instance.statusBuff.String()))
		}

		return buff.String()
	}
	return ""
}

func (instance *N) WriteStatus(w io.Writer) (int64, error) {
	if nil != instance {
		var buff bytes.Buffer
		var count int64 = 0

		buff.WriteString(instance.GetStatus())

		c, e := buff.WriteTo(w)
		if nil != e {
			return 0, e
		}
		count += c

		return c, e
	}
	return 0, nil
}

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

		// discovery
		instance.discovery = NewNodeDiscovery(instance.uuid, instance.Settings.Discovery)
		err := instance.discovery.Start()
		if nil != err {
			response = append(response, err)
		}

		// server
		instance.server = lygo_n_server.NewNServer(instance.Settings.Server)
		errs, warnings := instance.server.Start()
		response = append(response, errs...)
		instance.logWarns(warnings)

		// client
		instance.client = lygo_n_client.NewNClient(instance.Settings.Client)
		errs, warnings = instance.client.Start()
		response = append(response, errs...)
		instance.logWarns(warnings)

		// system message handler
		instance.handleSystemMessages()

		// check configuration
		if len(warnings) == 0 {
			warnings = instance.checkConfiguration()
			instance.logWarns(warnings)
		}

		return response
	}
	return nil
}

func (instance *N) close() []error {
	if instance.initialized {
		instance.initialized = false
		response := make([]error, 0)

		// discovery
		if nil != instance.discovery {
			response = append(response, instance.discovery.Stop())
		}

		// client
		if nil != instance.client {
			response = append(response, instance.client.Stop()...)
		}

		// server
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

func (instance *N) logWarns(warnings []string) {
	if len(warnings) > 0 {
		for _, w := range warnings {
			instance.statusBuff.WriteString(fmt.Sprintf("WARN", lygo_fmt.FormatDate(time.Now(), "yyyy/MM/dd HH:mm\n"), w))
			lygo_logs.Warn("", w)
		}
	}
}

func (instance *N) indentLines(text string) string {
	var buff bytes.Buffer
	text = strings.ReplaceAll(text, "\r", "")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			buff.WriteString("\t\t" + line + "\n")
		}
	}
	return buff.String()
}

func (instance *N) handleSystemMessages()  {
	if nil!=instance.server && instance.server.IsOpen(){

		// discovery messages
		if nil!=instance.discovery && instance.discovery.IsEnabled(){
			instance.server.RegisterCommand(CMD_GET_NODE_LIST, instance.discovery.getNodeList)
		}

	}
}
