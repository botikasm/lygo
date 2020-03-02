package lygo_stopwatch

import (
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type StopWatch struct {
	mutex     *sync.Mutex
	startTime time.Time
	endTime   time.Time
	Elapsed   time.Duration // int64 nanoseconds
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func New() *StopWatch {
	response := new(StopWatch)
	response.startTime = time.Now()
	response.endTime = response.startTime
	response.Elapsed = 0
	response.mutex = new(sync.Mutex)

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (watch *StopWatch) Start() {
	watch.mutex.Lock()
	defer watch.mutex.Unlock()
	watch.startTime = time.Now()
	watch.endTime = watch.startTime
}

func (watch *StopWatch) Stop() {
	watch.mutex.Lock()
	defer watch.mutex.Unlock()
	watch.endTime = time.Now()
	watch.recalculate()
}

func (watch *StopWatch) Milliseconds() int {
	return int(watch.Elapsed / 1000000)
}

func (watch *StopWatch) Seconds() float32 {
	return float32(watch.Milliseconds()) / 1000.0
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (watch *StopWatch) recalculate() {
	watch.Elapsed = watch.endTime.Sub(watch.startTime)
}
