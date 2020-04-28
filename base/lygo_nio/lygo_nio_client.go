package lygo_nio

import (
	"bufio"
	"crypto/rsa"
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
	Timeout time.Duration
	Secure  bool

	//-- private --//
	conn     net.Conn
	host     string
	port     int
	mux      sync.Mutex
	stopChan chan bool
	// RSA
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	serverKey  *rsa.PublicKey // server signature (got on handshake)
	sessionKey []byte
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNioClient(host string, port int) *NioClient {
	instance := new(NioClient)
	instance.host = host
	instance.port = port
	instance.Timeout = 10 * time.Second
	instance.stopChan = make(chan bool, 1)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) Open() error {
	if nil != instance {
		err := instance.initRSA()
		if nil != err {
			return err
		}
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
		instance.stopChan <- true
	}
	return nil
}

// Wait is stopped
func (instance *NioClient) Join() {
	// locks and wait for exit response
	<-instance.stopChan
}

func (instance *NioClient) Send(data interface{}) (*NioMessage, error) {
	if nil != instance {

		// creates NIO message
		message := new(NioMessage)
		message.Body = data

		return instance.send(message, false)
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) initRSA() error {
	if nil != instance && instance.Secure && nil == instance.privateKey {
		// TODO: implement loading from file

		// auto-generates
		pri, pub, err := keysGenerate(KEY_LEN)
		if nil != err {
			return err
		}
		instance.privateKey = pri
		instance.publicKey = pub
	}
	return nil
}

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
			if nil == err {
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
		HANDSHAKE.PublicKey = instance.publicKey
		response, err := instance.send(HANDSHAKE, true)
		if nil != err {
			return err
		}
		instance.serverKey = response.PublicKey
	}
	return nil
}

func (instance *NioClient) send(message *NioMessage, handshake bool) (*NioMessage, error) {
	if nil != instance {
		conn, err := instance.connect()
		if nil != err {
			_ = instance.Close() // reset connection
			return nil, err
		}

		if handshake {
			message.PublicKey = instance.publicKey
		}

		// ENCRYPT BODY
		if !handshake && nil != instance.publicKey && len(instance.sessionKey) > 0 {
			v := serialize(message.Body)
			data, err := encrypt(v, instance.sessionKey)
			if nil == err {
				message.Body = data
			} else {
				return nil, errors.Wrap(err, "Client Encryption error")
			}
		} else {
			message.Body = serialize(message.Body)
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
			// RESPONSE FROM SERVER
			if !handshake {
				if len(instance.sessionKey) > 0 {
					// DECRYPT BODY
					if v, b := response.Body.([]byte); b {
						data, err := decrypt(v, instance.sessionKey)
						if nil == err {
							response.Body = data
						}
					}
				}
			} else {
				// handshake
				if len(response.SessionKey) > 0 {
					data, err := decryptKey(response.SessionKey, instance.privateKey)
					if nil == err {
						instance.sessionKey = data
					}
				}
			}
			return &response, nil
		}
	}
	return nil, nil
}
