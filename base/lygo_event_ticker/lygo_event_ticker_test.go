package lygo_event_ticker

import (
	"fmt"
	"testing"
	"time"
)

var count int

func TestBasic(t *testing.T) {

	count = 0

	wd := NewEventTicker(time.Second*3, callback)
	wd.Start()

	wd.Join()
}

func callback(w *EventTicker) {
	fmt.Println("callback")
	count++
	if count == 3 {
		// stop the watchdog permanently
		w.Stop()
		w.Start()
		fmt.Println("STOPPED")
		return
	}
}
