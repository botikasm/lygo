package lygo_nio

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NioServer struct {
	PublicKey string

	//-- private --//
	port       int
	listener   net.Listener
	clients    int
	clientsMap map[string]*client
	mux        sync.Mutex
	handler    NioMessageHandler
}

type client struct {
	Id string
}

type NioMessageHandler func(message *NioMessage) interface{}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNioServer(port int) *NioServer {
	instance := new(NioServer)
	instance.clients = 0
	instance.port = port
	instance.clientsMap = make(map[string]*client)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioServer) Open() error {
	if nil != instance {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%v", instance.port))
		if nil != err {
			return err
		}
		instance.listener = listener

		// main listener loop
		go instance.open()
	}
	return nil
}

func (instance *NioServer) Close() error {
	if nil != instance {
		if nil != instance.listener {
			return instance.listener.Close()
		}
	}
	return nil
}

func (instance *NioServer) ClientsCount() int {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		return instance.clients
	}
	return 0
}

func (instance *NioServer) ClientsId() []string {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		keys := make([]string, 0, len(instance.clientsMap))
		for k := range instance.clientsMap {
			keys = append(keys, k)
		}
		return keys
	}
	return []string{}
}

func (instance *NioServer) OnMessage(callback NioMessageHandler) {
	if nil != instance {
		instance.handler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioServer) open() {
	for {
		// accept connections
		conn, err := instance.listener.Accept()
		if err != nil {
			// error accepting connection
			continue
		}
		go instance.handleConnection(conn)
	}
}

func (instance *NioServer) incClients(conn net.Conn) *client {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		c := new(client)
		c.Id = conn.RemoteAddr().String()
		instance.clients++
		instance.clientsMap[c.Id] = c
		return c
	}
	return nil
}

func (instance *NioServer) decClients(id string) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		_, ok := instance.clientsMap[id]
		if ok {
			instance.clients--
			delete(instance.clientsMap, id)
		}
	}
}

func (instance *NioServer) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// new client connection
	c := instance.incClients(conn)

	// connection loop
	for {
		var message NioMessage
		dec := gob.NewDecoder(rw)
		err := dec.Decode(&message)
		if nil != err {
			if err.Error() == "EOF" {
				// client disconnected
			}
			// exit
			break
		}

		if !instance.isHandshake(&message) && nil != instance.handler {
			clientKey := message.PublicKey
			if len(clientKey) > 0 {
				// TODO: decode client message body

			}
			customResponse := instance.handler(&message)
			if nil == customResponse {
				customResponse = true
			}

			err := sendResponse(customResponse, rw, &instance.PublicKey)
			if err != nil {
				break
			}
		} else {
			// response OK (default)
			err := sendResponse(true, rw, &instance.PublicKey)
			if err != nil {
				break
			}
		}
	}

	// client removed
	instance.decClients(c.Id)
}

func (instance *NioServer) isHandshake(message *NioMessage) bool {
	if v, b := message.Message.([]byte); b {
		return string(v) == string(HANDSHAKE.Message.([]byte))
	}
	return false
}

func sendResponse(body interface{}, rw *bufio.ReadWriter, publicKey *string) error {
	response := new(NioMessage)
	response.PublicKey = *publicKey

	// TODO: encode server message body
	response.Message = body

	enc := gob.NewEncoder(rw)
	err := enc.Encode(response)
	if err != nil {
		return err
	}
	err = rw.Flush()
	return err
}
