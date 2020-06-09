package lygo_n

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_fmt"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_sys"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_server"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_commons"
	"github.com/botikasm/lygo/ext/lygo_n/lygo_n_net"
	"io"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e s
// ---------------------------------------------------------------------------------------------------------------------

type N struct {
	Settings *lygo_n_commons.NSettings

	//-- private --//
	uuid        string
	statusBuff  bytes.Buffer
	initialized bool
	server      *lygo_n_net.NHost // nio server
	http        *NHttp            // http interface
	discovery   *NDiscovery
	events      *lygo_events.Emitter
	stopChan    chan bool
	status      int
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewNode(settings *lygo_n_commons.NSettings) *N {
	instance := new(N)
	instance.stopChan = make(chan bool, 1)
	instance.status = -1
	instance.Settings = settings

	if nil == instance.Settings {
		instance.Settings = new(lygo_n_commons.NSettings)
		instance.Settings.Discovery = new(lygo_n_commons.NDiscoverySettings)
		instance.Settings.Discovery.Publish = new(lygo_n_commons.NDiscoveryPublishSettings)
		instance.Settings.Discovery.Publish.Enabled = false
		instance.Settings.Discovery.Publisher = new(lygo_n_commons.NDiscoveryPublisherSettings)
		instance.Settings.Discovery.Publisher.Enabled = false
	}

	sysid, err := lygo_sys.ID()
	if nil != err {
		sysid = lygo_rnd.Uuid()
	}
	instance.uuid = fmt.Sprintf("%v", sysid)

	instance.events = lygo_events.NewEmitter()
	if len(instance.Settings.Name) == 0 {
		instance.Settings.Name = sysid
	}

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

func (instance *N) Name() string {
	if nil != instance {
		if len(instance.Settings.Name) > 0 {
			return instance.Settings.Name
		}
		return instance.GetUUID()
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

func (instance *N) Events() *lygo_events.Emitter {
	if nil != instance {
		return instance.events
	}
	return nil
}

func (instance *N) Start() []error {
	if nil != instance {
		if instance.status < 1 {
			instance.status = 1
			return instance.open()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *N) Stop() []error {
	if nil != instance {
		if instance.status == 1 {
			instance.status = 0
			return instance.close()
		}
	}
	return []error{lygo_n_commons.PanicSystemError}
}

func (instance *N) Join() {
	if nil != instance {
		if instance.status == 1 {
			<-instance.stopChan
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *N) Http() *lygo_http_server.HttpServer {
	if nil != instance {
		return instance.getHttp().Http()
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
//		m e s s a g e    h a n d l e r s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *N) RegisterCommand(command string, handler lygo_n_net.CommandHandler) {
	if nil != instance {
		server := instance.getServer()
		if nil != server {
			server.RegisterCommand(command, handler)
		}
	}
}

// Execute run local command
func (instance *N) Execute(commandName string, params map[string]interface{}) *lygo_n_commons.Response {
	if nil != instance {
		server := instance.getServer()
		if nil != server {
			return server.Execute(commandName, params)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
//		m e s s a g e    s e n d e r
// ---------------------------------------------------------------------------------------------------------------------

// Send try to run remote command, if not possible run a local command but always uses a connection as bridge if Network of Nodes is enabled
func (instance *N) Send(commandName string, params map[string]interface{}) *lygo_n_commons.Response {
	if nil != instance {
		if instance.initialized && nil != instance.discovery {
			if nil != instance.discovery && instance.discovery.IsNetworkOfNodesEnabled() {
				// try registered node
				conn := instance.discovery.AcquireNode() // always nil if network_id is empty
				if nil != conn {
					defer instance.discovery.ReleaseNode(conn)

					return conn.Send(commandName, params)
				} else {
					// try internal host if enabled
					conn = instance.getSelfHostedConn()
					if nil != conn {
						return conn.Send(commandName, params)
					}
				}
			}
		}
		// fallback
		return instance.Execute(commandName, params)
	}
	return &lygo_n_commons.Response{
		Error: lygo_n_commons.PanicSystemError.Error(),
		Data:  nil,
	}
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

		// [01] discovery: enable network discovery
		instance.discovery = NewNodeDiscovery(instance.events, instance.uuid, instance.Settings.Discovery)
		err := instance.discovery.Start()
		if nil != err {
			response = append(response, err)
		}

		// [02] server: enable client connections
		server := instance.getServer()
		if nil != instance.server {
			server.Info.Name = instance.Name()
			server.SetEventManager(instance.events)
			errs, warnings := server.Start()
			response = append(response, errs...)
			instance.logWarns(warnings)
			// system message handler
			instance.handleSystemMessages(server)
		}

		//[03] http interface
		instance.http = instance.getHttp()
		if nil != instance.http {
			instance.http.SetEventManager(instance.events)
			instance.http.SendCommandHandler = instance.Send
			errs, warnings := instance.http.Start()
			response = append(response, errs...)
			instance.logWarns(warnings)
		}

		return response
	}
	return nil
}

func (instance *N) getHttp() *NHttp {
	if nil != instance {
		if nil == instance.http {
			if nil != instance.Settings.Server && instance.Settings.Server.Enabled && instance.Settings.Server.Http.Enabled {
				instance.http = NewNHttp(instance.Settings.Server)
			}
		}
		return instance.http
	}
	return nil
}

func (instance *N) getServer() *lygo_n_net.NHost {
	if nil != instance {
		if nil == instance.server {
			if nil != instance.Settings.Server && instance.Settings.Server.Enabled {
				instance.server = lygo_n_net.NewNHost(instance.Settings.Server)
			}
		}
		return instance.server
	}
	return nil
}

func (instance *N) close() []error {
	if instance.initialized {
		instance.initialized = false
		response := make([]error, 0)

		// discovery
		if nil != instance.discovery {
			err := instance.discovery.Stop()
			if nil != err {
				response = append(response, err)
			}
		}

		// server
		if nil != instance.server {
			response = append(response, instance.server.Stop()...)
		}
		return response
	}
	return nil
}

func (instance *N) logWarns(warnings []string) {
	if len(warnings) > 0 {
		for _, w := range warnings {
			instance.statusBuff.WriteString(fmt.Sprintf("WARN", lygo_fmt.FormatDate(time.Now(), "yyyy/MM/dd HH:mm\n"), w))
			lygo_logs.Warn("", w)
		}
	}
}
func (instance *N) getSelfHostedConn() *lygo_n_net.NConn {
	if nil != instance.server && instance.server.IsOpen() {
		host := "localhost"
		port := instance.Settings.Server.Nio.Port()
		conn := lygo_n_net.NewNConn(host, port)
		errs, _ := conn.Start()
		if len(errs) == 0 {
			return conn
		}
	}
	return nil
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

func (instance *N) handleSystemMessages(server *lygo_n_net.NHost) {
	if nil != server && server.IsOpen() {

		// discovery messages
		if nil != instance.discovery && instance.discovery.IsEnabled() {
			server.RegisterCommand(lygo_n_commons.CMD_GET_NODE_LIST, instance.discovery.getNodeList)
			server.RegisterCommand(lygo_n_commons.CMD_REGISTER_NODE, instance.discovery.registerNode)
		}

	}
}
