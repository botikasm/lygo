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
}

type client struct {
	//-- private --//
	rw *bufio.ReadWriter
}

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

func (instance *NioServer) handleConnection(conn net.Conn) {
	// new client
	instance.mux.Lock()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	// new client connection
	id := conn.RemoteAddr().String()
	c := new(client)
	c.rw = rw
	instance.clients++
	instance.clientsMap[id] = c
	instance.mux.Unlock()

	// connection loop
	for {
		var message NioMessage
		dec := gob.NewDecoder(rw)
		err := dec.Decode(&message)
		if nil != err {
			if err.Error() == "EOF" {
				// client closed connection
				instance.mux.Lock()
				_, ok := instance.clientsMap[id]
				if ok {
					instance.clients--
					delete(instance.clientsMap, id)
				}
				instance.mux.Unlock()
			}
			// exit
			return
		}

		if !instance.isHandshake(&message) {
			clientKey := message.PublicKey

			// TODO : do something with message
			fmt.Println("SERVER", c, clientKey, message)

		}

		// response OK
		response := new(NioMessage)
		response.PublicKey = instance.PublicKey
		response.Message = true
		enc := gob.NewEncoder(rw)
		err = enc.Encode(response)
		if err != nil {
			return
		}
		err = rw.Flush()
		if err != nil {
			return
		}
	}
}

func (instance *NioServer) isHandshake(message *NioMessage) bool {
	if v, b := message.Message.([]byte); b {
		return string(v) == string(HANDSHAKE.Message.([]byte))
	}
	return false
}
