//go:generate go run lygo_resources_generator_launcher.go

package lygo_resources

import (
	"bufio"
	"os"
	"path"
)

//----------------------------------------------------------------------------------------------------------------------
//	v a r s
//----------------------------------------------------------------------------------------------------------------------


//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type resourceBox struct {
	storage map[string][]byte
}

func newResourceBox() *resourceBox {
	return &resourceBox{storage: make(map[string][]byte)}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Find a file
func (r *resourceBox) Has(file string) bool {
	if _, ok := r.storage[file]; ok {
		return true
	}
	return false
}

// Get file's content
func (r *resourceBox) Get(file string) ([]byte, bool) {
	if f, ok := r.storage[file]; ok {
		return f, ok
	}
	return nil, false
}

// Add a file to box
func (r *resourceBox) Add(file string, content []byte) {
	r.storage[file] = content
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// private resource container
var resources = newResourceBox()

// Get a file from box
func Get(resource string) ([]byte, bool) {
	return resources.Get(resource)
}

// Add a file content to box
func Add(resource string, content []byte) {
	resources.Add(resource, content)
}

// Has a file in box
func Has(resource string) bool {
	return resources.Has(resource)
}

func SaveToDir(dir, resource string) (string, bool) {
	p := path.Join(dir, resource)
	return SaveTo(p, resource)
}

func SaveTo(outFileName, resource string) (string, bool) {
	data, found := resources.Get(resource)
	if found {
		f, err := os.Create(outFileName)
		if err == nil {
			defer f.Close()
			w := bufio.NewWriter(f)
			_, err = w.Write(data)
			_ = w.Flush()
		}
		return f.Name(), true
	}
	return "", found
}
