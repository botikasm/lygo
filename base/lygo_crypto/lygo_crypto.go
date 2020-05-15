package lygo_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/botikasm/lygo/base/lygo_io"
	"io"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func GenerateSessionKey() [32]byte {
	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader
	key := make([]byte, 32)
	if _, err := io.ReadFull(rng, key); err != nil {
		panic("RNG failure")
	}
	return sha256.Sum256(key)
}

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func EncryptTextWithPrefix(text string, key []byte) ([]byte, error) {
	if strings.Index(text, "enc-") == -1 {
		return EncryptBytesAES([]byte(text), key)
	}
	return []byte(text), nil
}

func DecryptTextWithPrefix(text string, key []byte) ([]byte, error) {
	if strings.Index(text, "enc-") != -1 {
		return DecryptBytesAES([]byte(text), key)
	}
	return []byte(text), nil
}

func EncryptTextAES(text string, key []byte) ([]byte, error) {
	return EncryptBytesAES([]byte(text), key)
}

func EncryptFileAES(fileName string, key []byte, optOutFileName string) ([]byte, error) {
	data, err := lygo_io.ReadBytesFromFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	encoded, err := EncryptBytesAES(data, key)
	if err != nil {
		return []byte{}, err
	}

	// write file
	if len(optOutFileName) > 0 {
		_, err := lygo_io.WriteBytesToFile(encoded, optOutFileName)
		if err != nil {
			return []byte{}, err
		}
	} else {
		_, err := lygo_io.WriteBytesToFile(encoded, fileName)
		if err != nil {
			return []byte{}, err
		}
	}
	return encoded, nil
}

func EncryptBytesAES(data []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key) // key must be 32 bytes
	// if there are any errors, handle them
	if err != nil {
		return []byte{}, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return []byte{}, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptTextAES(text string, key []byte) ([]byte, error) {
	return DecryptBytesAES([]byte(text), key)
}

func DecryptFileAES(fileName string, key []byte, optOutFileName string) ([]byte, error) {
	data, err := lygo_io.ReadBytesFromFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	encoded, err := DecryptBytesAES(data, key)
	if err != nil {
		return []byte{}, err
	}

	// write file
	if len(optOutFileName) > 0 {
		_, err := lygo_io.WriteBytesToFile(encoded, optOutFileName)
		if err != nil {
			return []byte{}, err
		}
	} else {
		_, err := lygo_io.WriteBytesToFile(encoded, fileName)
		if err != nil {
			return []byte{}, err
		}
	}
	return encoded, nil
}

func DecryptBytesAES(data []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return []byte{}, err
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return []byte{}, err
	}

	return plain, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
