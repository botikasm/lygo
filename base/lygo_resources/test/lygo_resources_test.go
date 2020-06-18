package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_exec"
	"github.com/botikasm/lygo/base/lygo_resources"
	"github.com/botikasm/lygo/base/lygo_resources/test/resources"
	"os"
	"testing"
)

const packageName = "resources"
const startDirectory = "src_resources"

func TestGenerator(t *testing.T) {
	// run generator
	generator := lygo_resources.NewGenerator()
	generator.Package = packageName
	generator.StartDirectory = startDirectory
	// generator.OutputFile = "./" + packageName + "/blob_{{ .count }}.go"
	generator.Exclude = []string{"/excluded/"}
	generator.ForceSingleResourceFile = true // creates a single file ignoring custom "OutputFile" param
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
	resName := startDirectory + "/my_resource.txt"
	// get resource from packageName
	data, found := resources.Get(resName)
	if !found {
		t.Error("Resource not found")
		t.FailNow()
	}
	fmt.Println(resName, ":\n", string(data))
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
