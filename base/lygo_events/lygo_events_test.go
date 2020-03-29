package lygo_events

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

func TestTick(t *testing.T) {

	i := 0
	emitter := NewEmitter()
	emitter.Tick(3*time.Second, func(ticker *EventTicker) {
		i++
		switch i {
		case 3:
			ticker.Stop()
		}
		fmt.Println("emitter.Tick")
	}).Join()

}

func TestEvents(t *testing.T) {
	emitter := NewEmitter()
	emitter.On("my-event", func(event *Event) {
		fmt.Println("1)", event.Name, event.Arguments)
	})
	emitter.On("my-event", func(event *Event) {
		fmt.Println("2)", event.Name, event.Arguments)
	})
	emitter.On("my-event", listener3)
	emitter.Emit("no listener")
	emitter.Emit("my-event", "arg1", 2, 3, "arg4")

	emitter.Off("my-event", listener3)
	emitter.Emit("my-event", "SHOULD NOT BE HANDLED FROM LISTENER 3")
	emitter.Off("my-event")
	emitter.Emit("my-event", "SHOULD NOT HANDLE THIS")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func listener3(event *Event) {
	fmt.Println("3)", event.Name, event.Arguments)
}

func callback(w *EventTicker) {
	fmt.Println("callback")
	count++
	if count == 3 {
		// stop the ticker permanently
		w.Stop()
		w.Start() // do not start after was stopped
		fmt.Println("STOPPED")
		return
	} else if count == 1 {
		w.Pause()
		fmt.Println("PAUSED")
		time.Sleep(3 * time.Second)
		w.Resume()
		fmt.Println("RESUMED")
	}
}
