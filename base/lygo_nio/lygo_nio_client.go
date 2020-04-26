package lygo_nio

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NioClient struct {
	Timeout   time.Duration
	PublicKey string

	//-- private --//
	conn      net.Conn
	host      string
	port      int
	mux       sync.Mutex
	serverKey string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNioClient(host string, port int) *NioClient {
	instance := new(NioClient)
	instance.host = host
	instance.port = port
	instance.Timeout = 10 * time.Second

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) Open() error {
	if nil != instance {
		return instance.handshake()
	}
	return nil
}

func (instance *NioClient) Close() error {
	if nil != instance {
		if nil != instance.conn {
			err := instance.conn.Close()
			instance.conn = nil
			return err
		}
	}
	return nil
}

func (instance *NioClient) Send(data interface{}) (*NioMessage, error) {
	if nil != instance {

		// creates NIO message
		message := new(NioMessage)
		message.PublicKey = instance.PublicKey
		message.Message = data

		return instance.send(message)
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) test() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", instance.host, instance.port), instance.Timeout)
	if nil != conn {
		defer conn.Close()
	}
	return err
}

func (instance *NioClient) connect() (net.Conn, error) {
	if nil != instance {
		if nil == instance.conn {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", instance.host, instance.port), instance.Timeout)
			if nil==err{
				instance.conn = conn
			}
			return conn, err
		}
		return instance.conn, nil
	}
	return nil, nil
}

func (instance *NioClient) handshake() error {
	if nil != instance {
		HANDSHAKE.PublicKey = instance.PublicKey
		response, err := instance.send(HANDSHAKE)
		if nil != err {
			return err
		}
		instance.serverKey = response.PublicKey
	}
	return nil
}

func (instance *NioClient) send(message *NioMessage) (*NioMessage, error) {
	if nil != instance {
		conn, err := instance.connect()
		if nil != err {
			_ = instance.Close() // reset connection
			return nil, err
		}

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		enc := gob.NewEncoder(rw)
		err = enc.Encode(message)
		if err != nil {
			return nil, errors.Wrapf(err, "Encode failed for message: %#v", message)
		}
		err = rw.Flush()
		if err != nil {
			return nil, errors.Wrap(err, "Flush failed.")
		}

		// read NIO response
		var response NioMessage
		dec := gob.NewDecoder(rw)
		err = dec.Decode(&response)
		if err != nil {
			// fmt.Println(errors.Wrap(err, "Client failed to read response"))
			return nil, errors.Wrap(err, "Client failed to read response")
		} else {
			//fmt.Println("MESSAGE RECEIVED FROM SERVER", response)
			return &response, nil
		}
	}
	return nil, nil
}
