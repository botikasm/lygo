package lygo_mu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFmt(t *testing.T) {

	exp := "1.00 KB"
	got := FmtBytes(1024)
	assert.Equal(t, exp, got, "Expected %s but got %s", exp, got)
	fmt.Println(exp)

	exp = "1.33 TB"
	got = FmtBytes(1457688937635)
	assert.Equal(t, exp, got, "Expected %s but got %s", exp, got)
	fmt.Println(exp)

	exp = "500 B"
	got = FmtBytes(500)
	assert.Equal(t, exp, got, "Expected %s but got %s", exp, got)
	fmt.Println(exp)
}
