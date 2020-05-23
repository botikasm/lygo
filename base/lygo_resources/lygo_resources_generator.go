package lygo_resources

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//----------------------------------------------------------------------------------------------------------------------
//	v a r s
//----------------------------------------------------------------------------------------------------------------------

var packageTemplate = template.Must(template.New("").Funcs(map[string]interface{}{"conv": formatByteSlice}).Parse(`
// Code generated by go generate; DO NOT EDIT.
// generated using files from resources directory
// DO NOT COMMIT this file
package {{.Package}}

import "github.com/botikasm/lygo/base/lygo_resources"

func init(){
	{{- range $name, $file := .Resources }}
    	lygo_resources.Resources.Add("{{ $name }}", []byte{ {{ conv $file }} })
	{{- end }}
}
`))

const dirResources = "resources"
const outFileName = "blobResources.go"
const pkg = "lygo_resources"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Generator struct {
	Directory  string
	OutputFile string
	Package    string
	Exclude    []string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewGenerator() *Generator {
	instance := new(Generator)
	instance.Directory = dirResources
	instance.OutputFile = outFileName
	instance.Package = pkg
	instance.Exclude = make([]string, 0)
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Generator) Start() {
	fmt.Println("---------------------------------------")
	fmt.Println("Packing resources starting from directory '" + instance.Directory + "'")

	if _, err := os.Stat(instance.Directory); os.IsNotExist(err) {
		fmt.Println("Resources directory does not exists!")
		return
	}

	count := 0
	context := make(map[string]interface{})
	resources := make(map[string][]byte)
	context["Package"] = instance.Package
	context["Resources"] = resources
	err := filepath.Walk(instance.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error :", err)
			return err
		}
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, "resources"))
		if info.IsDir() {
			fmt.Println("[DIR] ", relativePath)
			return nil
		} else {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading %s: %s", path, err)
				return err
			}
			if !instance.isExcluded(relativePath) {
				fmt.Println("\t* INCLUDING: ", relativePath)
				count++
				resources[relativePath] = b
			} else {
				fmt.Println("\t! EXCLUDING: ", relativePath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through resources directory:", err)
		return
	}

	//lygo_paths.Mkdir(outputFileName)
	f, err := os.Create(instance.OutputFile)
	if err != nil {
		fmt.Println("Error creating blob file:", err)
		return
	}
	defer f.Close()

	builder := &bytes.Buffer{}

	// solve template
	err = packageTemplate.Execute(builder, context)
	if err != nil {
		fmt.Println("Error executing template", err)
		return
	}

	data, err := format.Source(builder.Bytes())
	if err != nil {
		fmt.Println("Error formatting generated code", err)
		return
	}
	err = ioutil.WriteFile(instance.OutputFile, data, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing blob file", err)
		return
	}

	fmt.Println("Packing resources done...")
	fmt.Println("TOTAL RESOURCES: ", count)
	fmt.Println("DO NOT COMMIT " + instance.OutputFile)
}

func (instance *Generator) isExcluded(path string) bool {
	for _, s := range instance.Exclude {
		if strings.Index(path, s) == 0 {
			return true
		}
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func formatByteSlice(sl []byte) string {
	builder := strings.Builder{}
	for _, v := range sl {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}
	return builder.String()
}