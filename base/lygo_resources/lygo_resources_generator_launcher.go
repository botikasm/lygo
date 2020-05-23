//+build ignore

package main

import "github.com/botikasm/lygo/base/lygo_resources"

func main(){
	var generator *lygo_resources.Generator = lygo_resources.NewGenerator()
	generator.Package = "test"
	generator.Directory = "./test/resources"
	generator.OutputFile = "./test/blobResources.go"
	generator.Start()
}

