# lygo_resources
Embedd resource into executables using go:generate features.

This little package uses special comment with `//go:generate go run gen.go`.

That's the magic.

## How to Use

* Create a `./resources` directory and put here all your resource.
* Run `go generate ./...`

Sample:
```
func TestResource(t *testing.T) {
	// get resource
	data, found := Get("/my_resource.txt")
	if !found {
		t.Error("Resource not found")
		t.FailNow()
	}
	fmt.Println(string(data))
}
```
