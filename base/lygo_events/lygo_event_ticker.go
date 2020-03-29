package lygo_events

import "time"

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
}

type EventTickerCallback func(*EventTicker)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewEventTicker(timeout time.Duration, callback EventTickerCallback) *EventTicker {
	w := &EventTicker{
		timer:    time.NewTicker(timeout),
		timeout:  timeout,
		stopChan: make(chan bool, 1),
		callback: callback,
	}

	go w.loop()

	return w
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Wait EventTicker is stopped
func (w *EventTicker) Join() {
	// locks and wait for exit response
	<-w.stopChan
}

// Start .... Start the timer
func (w *EventTicker) Start() {
	if nil != w.timer {
		w.timer.Stop()
		w.timer = nil
	}
	w.timer = time.NewTicker(w.timeout)
}

// Stop ... stops the timer
func (w *EventTicker) Stop() {
	if nil != w.timer {
		w.timer.Stop()
		w.stopChan <- true
		w.timer = nil
	}
}

func (w *EventTicker) Pause() {
	if nil != w.timer && !w.paused {
		w.paused = true
	}
}

func (w *EventTicker) Resume() {
	if nil != w.timer && w.paused {
		w.paused = false
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (w *EventTicker) loop() {
	for {
		if nil != w.timer {
			select {
			case <-w.stopChan:
				return
			case <-w.timer.C:
				// event
				if nil != w.callback && !w.paused {
					w.callback(w)
				}
			}
		}
	}
}
