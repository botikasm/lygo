//+build ignore

package main

import "github.com/botikasm/lygo/base/lygo_resources"

// THIS FILE IS USED FROM generate.sh

func main(){
	var generator *lygo_resources.Generator = lygo_resources.NewGenerator()
	generator.Package = "resources"
	generator.StartDirectory = "./test/src_resources"
	generator.OutputFile = "./test/resources/blob_{{ .count }}.go"
	generator.Start()
}

