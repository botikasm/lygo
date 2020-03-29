package lygo_events

import (
	"reflect"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Emitter struct {
	listeners map[string][]func(event *Event)
}

type Event struct {
	Name      string
	Arguments []interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func NewEmitter() *Emitter {
	instance := new(Emitter)
	instance.listeners = make(map[string][]func(event *Event))
	return instance
}

// Simple ticker loop
func (instance *Emitter) Tick(timeout time.Duration, callback EventTickerCallback) *EventTicker {
	et := NewEventTicker(timeout, callback)
	et.Start()

	return et
}

func (instance *Emitter) On(eventName string, callback func(event *Event)) {
	instance.listeners[eventName] = append(instance.listeners[eventName], callback)
}

func (instance *Emitter) Off(eventName string, callback ...func(event *Event)) {
	if _, ok := instance.listeners[eventName]; ok {
		if len(callback) == 0 {
			instance.listeners[eventName] = make([]func(event *Event), 0)
		} else {
			handlers := instance.listeners[eventName]
			// loop starting from end
			for i := len(handlers) - 1; i > -1; i-- {
				f := handlers[i]
				for _, h := range callback {
					v1 := reflect.ValueOf(f)
					v2 := reflect.ValueOf(h)
					if v1 == v2 {
						handlers = removeIndex(handlers, i)
						break
					}
				}

			}
			instance.listeners[eventName] = handlers
		}
	}
}

func (instance *Emitter) Emit(eventName string, args ...interface{}) {
	for k, handlers := range instance.listeners {
		if k == eventName {
			for _, handler := range handlers {
				if nil != handler {
					handler(&Event{
						Name:      k,
						Arguments: args,
					})
				}
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func removeIndex(a []func(event *Event), index int) []func(event *Event) {
	return append(a[:index], a[index+1:]...)
}
