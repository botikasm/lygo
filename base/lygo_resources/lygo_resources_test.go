package lygo_resources

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_exec"
	"testing"
)

func TestGenerator(t *testing.T) {
	// run generator
	_, err := lygo_exec.Run("go", "generate", "./...")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func TestResource(t *testing.T) {
	// get resource
	data, found := Get("/my_resource.txt")
	if !found {
		t.Error("Resource not found")
		t.FailNow()
	}
	fmt.Println(string(data))
}
