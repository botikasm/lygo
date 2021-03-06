package lygo_n_commons

import (
	"errors"
	"github.com/botikasm/lygo/base/lygo_rnd"
)

const Version = "1.0.3"

// ---------------------------------------------------------------------------------------------------------------------
//		t o k e n s
// ---------------------------------------------------------------------------------------------------------------------

var AppToken = ""

// ---------------------------------------------------------------------------------------------------------------------
//		e r r o r s
// ---------------------------------------------------------------------------------------------------------------------

var (
	PanicSystemError            = errors.New("panic_system_error")
	CommandNotFoundError        = errors.New("command_not_found")
	UnsupportedMessageTypeError = errors.New("unsupported_message_type")
	InvalidTokenError           = errors.New("invalid_token_error")
	EmptyResponseError           = errors.New("empty_response_error")

	// warns
	ServerNotEnabledWarning       = errors.New("server_not_enabled_warning")
	ClientNotEnabledWarning     = errors.New("client_not_enabled_warning")
)

// ---------------------------------------------------------------------------------------------------------------------
//		e r r o r    c o n t e x t
// ---------------------------------------------------------------------------------------------------------------------

const (
	ContextDatabase  = "db"
	ContextWebsocket = "ws"
)

// ---------------------------------------------------------------------------------------------------------------------
//		e v e n t s
// ---------------------------------------------------------------------------------------------------------------------

const (
	EventError = "on_error"

	// discovery events
	EventQuitDiscovery = "on_quit_discovery"
	EventNewNode = "on_new_node"
	EventNewPublisher = "on_new_publisher"
	EventRemovedNode = "on_removed_node"
	EventRemovedPublisher = "on_removed_publisher"
)

// ---------------------------------------------------------------------------------------------------------------------
//		c o m m a n d s
// ---------------------------------------------------------------------------------------------------------------------

const (
	CMD_GET_NODE_LIST = "n.sys_get_node_list"
	CMD_REGISTER_NODE = "n.sys_register_node"
)

// ---------------------------------------------------------------------------------------------------------------------
//		i n i t i a l i z a t i o n
// ---------------------------------------------------------------------------------------------------------------------

func init() {
	AppToken = lygo_rnd.Uuid()
}
