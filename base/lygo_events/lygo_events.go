package lygo_events

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_strings"
	"reflect"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

/**
 *  THREAD SAFE EVENT EMITTER
 */
type Emitter struct {
	listeners map[string][]func(event *Event)
	mux       sync.Mutex
}

type Event struct {
	Name      string
	Arguments []interface{}
}

func (instance *Event) ArgumentsInterface() interface{} {
	return interface{}(instance.Arguments)
}

func (instance *Event) Argument(index int) interface{} {
	if len(instance.Arguments) > index {
		return instance.Arguments[index]
	}
	return nil
}

func (instance *Event) ArgumentAsError(index int) error {
	v := instance.Argument(index)
	if nil != v {
		if e, b := v.(error); b {
			return e
		}
	}
	return nil
}

func (instance *Event) ArgumentAsString(index int) string {
	v := instance.Argument(index)
	return lygo_conv.ToString(v)
}

func (instance *Event) ArgumentAsInt(index int) int {
	v := instance.Argument(index)
	if nil != v {
		return lygo_conv.ToInt(v)
	}
	return -1
}

func (instance *Event) ArgumentAsBytes(index int) []byte {
	v := instance.Argument(index)
	if nil != v {
		if e, b := v.([]byte); b {
			return e
		}
	}
	return nil
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

func (instance *Emitter) Has(eventName string) bool {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		if _, b := instance.listeners[eventName]; b {
			return len(instance.listeners[eventName]) > 0
		}
	}
	return false
}

func (instance *Emitter) On(eventName string, callback func(event *Event)) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		instance.listeners[eventName] = append(instance.listeners[eventName], callback)
	}
}

func (instance *Emitter) Off(eventName string, callback ...func(event *Event)) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
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
}

func (instance *Emitter) Clear() {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		instance.listeners = make(map[string][]func(event *Event), 0)
	}
}

func (instance *Emitter) Emit(eventName string, args ...interface{}) {
	if nil != instance {
		instance.emit(eventName, args...)
	}
}

func (instance *Emitter) EmitAsync(eventName string, args ...interface{}) {
	if nil != instance {
		go instance.emit(eventName, args...)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func removeIndex(a []func(event *Event), index int) []func(event *Event) {
	return append(a[:index], a[index+1:]...)
}

func (instance *Emitter) emit(eventName string, args ...interface{}) {
	if nil != instance {
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic
				message := lygo_strings.Format("Emit '%s' ERROR: %s", eventName, r)
				fmt.Println(message)
			}
		}()
		instance.mux.Lock()
		defer instance.mux.Unlock()
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
}
