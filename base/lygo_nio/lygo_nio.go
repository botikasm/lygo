package lygo_nio

import (
	"crypto/rsa"
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/botikasm/lygo/base/lygo_json"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const KEY_LEN = 1024 * 3

var (
	HANDSHAKE = &NioMessage{
		PublicKey: nil,
		Body:      []byte("ACK"),
	}
)

type NioMessage struct {
	PublicKey  *rsa.PublicKey // public key for response
	SessionKey []byte         // session key
	Body       interface{}    // message object
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func serialize(data interface{}) []byte {
	if nil != data {
		if v, b := data.([]byte); b {
			return v
		} else if v, b := data.(string); b {
			return []byte(v)
		}
		return lygo_json.Bytes(data)
	}
	return []byte{}
}

func keysGenerate(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return lygo_crypto.GenerateKeyPair(bits)
}

func newSessionKey() [32]byte {
	return lygo_crypto.GenerateSessionKey()
}

func encryptKey(data []byte, key *rsa.PublicKey) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := lygo_crypto.EncryptWithPublicKey(data, key)
		return response, err
	}
	return []byte{}, nil
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := lygo_crypto.EncryptBytesAES(data, key)
		return response, err
	}
	return []byte{}, nil
}

func decryptKey(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := lygo_crypto.DecryptWithPrivateKey(data, privateKey)
		return response, err
	}
	return []byte{}, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := lygo_crypto.DecryptBytesAES(data, key)
		return response, err
	}
	return []byte{}, nil
}
