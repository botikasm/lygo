//+build ignore

package main

import "github.com/botikasm/lygo/base/lygo_resources"

func main(){
	var generator *lygo_resources.Generator = lygo_resources.NewGenerator()
	generator.Start()
}

