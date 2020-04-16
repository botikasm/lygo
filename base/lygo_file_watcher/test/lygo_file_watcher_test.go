package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_file_watcher"
	"github.com/botikasm/lygo/base/lygo_paths"
	"testing"
	"time"
)

func TestFileWatcher(t *testing.T) {

	path := lygo_paths.Absolute("./")

	watcher := lygo_file_watcher.New()
	watcher.FilterOps(lygo_file_watcher.Move, lygo_file_watcher.Create, lygo_file_watcher.Remove)
	watcher.IgnoreHiddenFiles(true)
	watcher.AddRecursive(path)

	fmt.Println("Watching at: ", path)

	go (func() {
		for {
			select {
			case event := <-watcher.Event:
				go handle(&event)
			case err := <-watcher.Error:
				fmt.Println(err)
				return
			case <-watcher.Closed:
				return
			}
		}
		watcher.Close()
	})()

	watcher.Start(1 * time.Second)

	fmt.Println("Exiting.")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func handle(event *lygo_file_watcher.Event) {
	fmt.Println(event.Op, event.Path, event.FileInfo) // Print the event's info.
}
