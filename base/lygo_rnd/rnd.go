package lygo_rnd

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Guid() string {
	return uuid.New().String()
}

func Uuid() (string, error) {
	uuid := ""
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if nil == err {
		uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	}

	return uuid, err
}

func UuidTimestamp() (string, error) {
	uuid, err := Uuid()
	if nil == err {
		return time.Now().Format("20060102T150405") + "-" + uuid, nil
	}
	return "", err
}

func UuidDefault(defVal string) string {
	uuid, _ := Uuid()
	if len(uuid) == 0 {
		uuid = defVal
	}
	return uuid
}

func UuidTimestampDefault(defVal string) string {
	uuid, _ := UuidTimestamp()
	if len(uuid) == 0 {
		uuid = defVal
	}
	return uuid
}
