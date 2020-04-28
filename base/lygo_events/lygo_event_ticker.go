package lygo_events

import (
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

// EventTicker ... use the GetKickChannel() to get notified when the watchdog barks
type EventTicker struct {
	timer    *time.Ticker
	timeout  time.Duration
	stopChan chan bool
	paused   bool
	callback EventTickerCallback
	mux sync.Mutex
	locked bool
}

type EventTickerCallback func(*EventTicker)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewEventTicker(timeout time.Duration, callback EventTickerCallback) *EventTicker {
	instance := &EventTicker{
		timer:    time.NewTicker(timeout),
		timeout:  timeout,
		stopChan: make(chan bool, 1),
		callback: callback,
	}

	go instance.loop()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Wait EventTicker is stopped
func (instance *EventTicker) Join() {
	// locks and wait for exit response
	<-instance.stopChan
}

// Start .... Start the timer
func (instance *EventTicker) Start() {
	if nil != instance.timer {
		instance.timer.Stop()
		instance.timer = nil
	}
	instance.timer = time.NewTicker(instance.timeout)
}

// Stop ... stops the timer
func (instance *EventTicker) Stop() {
	if nil != instance.timer {
		instance.timer.Stop()
		instance.stopChan <- true
		instance.timer = nil
	}
}

func (instance *EventTicker) Pause() {
	if nil != instance.timer && !instance.paused {
		instance.paused = true
	}
}

func (instance *EventTicker) Resume() {
	if nil != instance.timer && instance.paused {
		instance.paused = false
	}
}

func (instance *EventTicker) Lock() {
	if nil!= instance && !instance.locked {
		instance.locked = true
		instance.mux.Lock()
	}
}

func (instance *EventTicker) Unlock() {
	if nil!= instance && instance.locked {
		instance.mux.Unlock()
		instance.locked = false
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *EventTicker) loop() {
	if nil!= instance {
		for {
			if nil!= instance && nil != instance.timer {
				select {
				case <-instance.stopChan:
					return
				case <-instance.timer.C:
					// event
					if nil != instance.callback && !instance.paused {
						// thread safe call
						instance.Lock()
						instance.callback(instance)
						instance.Unlock()
					}
				}
			} else{
				return
			}
		}
	}

}
