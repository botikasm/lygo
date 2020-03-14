package lygo_num2word

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimple(t *testing.T) {

	converter := NewNum2Word()

	exp := "uno"
	out := converter.Convert(1, "it")
	fmt.Println(out)
	assert.EqualValues(t, exp, out, "Expected %s, got %s", exp, out)

	exp = "mille cinquecento quarantadue"
	out = converter.Convert(1542, "it")
	fmt.Println(out)
	assert.EqualValues(t, exp, out, "Expected %s, got %s", exp, out)

	exp = "due mila cinquecento quarantadue"
	out = converter.Convert(2542, "it")
	fmt.Println(out)
	assert.EqualValues(t, exp, out, "Expected %s, got %s", exp, out)

	converter.Options.WordSeparator = ""
	exp = "duemilacinquecentoquarantadue"
	out = converter.Convert(2542, "it")
	fmt.Println(out)
	assert.EqualValues(t, exp, out, "Expected %s, got %s", exp, out)

}
