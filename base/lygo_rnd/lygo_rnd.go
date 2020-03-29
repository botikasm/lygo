package lygo_rnd

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// New creates a new random UUID or panics.
func Uuid() string {
	return uuid.New().String()
}

func RndId() string {
	length := 32
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func UuidTimestamp() string {
	return time.Now().Format("20060102T150405") + "-" + Uuid()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
