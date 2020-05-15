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
	text := "Mario Rossi "
	key:=[]byte("user_0001")
	enc, err := EncryptTextWithPrefix(text, key)
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
