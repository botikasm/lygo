# lygo_resources
Embedd resource into executables using go:generate features.

This little package uses special comment with `//go:generate go run gen.go`.

That's the magic.

## Credits

Thank you at the great [post](https://levelup.gitconnected.com/how-i-embedded-resources-in-go-514b72f6ef0a) of Kasun Vithanage
 

## How to Use

* Create a `./resources` directory and put here all your resource.
* Run `go generate ./...` or..
* ... use `lygo_resources.Generator` library as below

Sample:
```
func Generate(){
    var generator *lygo_resources.Generator = lygo_resources.NewGenerator()
    generator.Package = "resources"
    generator.Directory = "./test/src_resources"
    generator.OutputFile = "./test/resources/blob{{ .count }}.go"
    generator.Start()
}
func UseResource() {
	// get resource: uses "resources.go" generated file
	data, found := resources.Get("/my_resource.txt")
	if !found {
		fmt.Prinln("Resource not found")
	}
	fmt.Println(string(data))
}
```
