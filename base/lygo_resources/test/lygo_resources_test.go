package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_exec"
	"github.com/botikasm/lygo/base/lygo_resources"
	"os"
	"testing"
)

func TestGenerator(t *testing.T) {
	// run generator
	generator := lygo_resources.NewGenerator()
	generator.Package = "test"
	generator.Exclude = []string{"/excluded/"}
	generator.Start()
}


func TestGeneratorSh(t *testing.T) {
	// run generator
	_, err := lygo_exec.Run("go", "generate", "./...")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func TestResource(t *testing.T) {
	// get resource
	data, found := lygo_resources.Get("/my_resource.txt")
	if !found {
		t.Error("Resource not found")
		t.FailNow()
	}
	fmt.Println(string(data))
}

func TestSaveToPath(t *testing.T) {
	// get resource
	fileName, found := lygo_resources.SaveTo("./image.png", "/more/image.png")
	if !found {
		t.Error("Resource not found")
		t.FailNow()
	}
	os.Remove(fileName)
}
