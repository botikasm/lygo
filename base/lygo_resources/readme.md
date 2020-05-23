# lygo_resources
Embedd resource into executables using go:generate features.

This little package uses special comment with `//go:generate go run gen.go`.

That's the magic.

## How to Use

* Create a `./resources` directory and put here all your resource.
* Run `go generate ./...`

Sample:
```
func UseResource() {
	// get resource
	data, found := lygo_resources.Get("/my_resource.txt")
	if !found {
		fmt.Prinln("Resource not found")
	}
	fmt.Println(string(data))
}
```
