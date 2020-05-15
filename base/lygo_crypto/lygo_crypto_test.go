package lygo_crypto

import (
	"fmt"
	"testing"
)

func TestSessionKey(t *testing.T) {
	key := GenerateSessionKey()
	if len(key) != 64 {
		t.Error("Bad key length")
		t.FailNow()
	}
	fmt.Println(key)
	fmt.Println(fmt.Printf("% x", key))
}

func TestAESWithPrefix(t *testing.T) {
	key:=[]byte("1234")
	enc, err := EncryptTextWithPrefix("hola", key)
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(enc)
	enc, err = EncryptTextWithPrefix(enc, key)
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(enc)

	dec, err := DecryptTextWithPrefix(enc, key)
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(dec)
}
