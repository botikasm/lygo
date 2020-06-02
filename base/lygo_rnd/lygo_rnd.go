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

func Between(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max-min) + min
}

func BetweenDuration(max, min int64) time.Duration {
	return time.Duration(Between(max, min))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
