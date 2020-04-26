package lygo_crypto

import (
	"fmt"
	"testing"
)

func TestSessionKey(t *testing.T) {
	key := GenerateSessionKey(16)
	if len(key) != 64 {
		t.Error("Bad key length")
		t.FailNow()
	}
	fmt.Println(key)
	fmt.Println(fmt.Printf("% x", key))
}
