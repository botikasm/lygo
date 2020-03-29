package lygo_rnd

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {

	guid := Uuid()
	fmt.Println(guid)

	uuid := RndId()
	fmt.Println(uuid)

	uuid_t := UuidTimestamp()
	fmt.Println(uuid_t)

}
