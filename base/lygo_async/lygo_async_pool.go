package lygo_async

import (
	"sync/atomic"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const (
	// DefaultMaxConcurrent is max number of concurrent routines
	DefaultMaxConcurrent = 100
)

type ConcurrentPool struct {
	limit                int
	tickets              chan int
	numOfRunningRoutines int32
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewConcurrentPool(limit int) *ConcurrentPool {
	if limit < 1 {
		limit = DefaultMaxConcurrent
	}

	// allocate a limiter instance
	instance := new(ConcurrentPool)
	instance.limit = limit
	instance.tickets = make(chan int, limit) // buffered channel with a limit
	instance.numOfRunningRoutines = 0

	// allocate the tickets:
	for i := 0; i < instance.limit; i++ {
		instance.tickets <- i
	}

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConcurrentPool) Limit() int {
	return instance.limit
}

func (instance *ConcurrentPool) Count() int32 {
	return atomic.LoadInt32(&instance.numOfRunningRoutines)
}

// Execute adds a function to the execution queue.
// If num of go routines allocated by this instance is < limit
// launch a new go routine to execute job
// else wait until a go routine becomes available
func (instance *ConcurrentPool) Run(job func()) int {
	// pop a ticket
	ticket := <-instance.tickets
	atomic.AddInt32(&instance.numOfRunningRoutines, 1)
	go func() {
		defer func() {
			// push a ticket
			instance.tickets <- ticket
			atomic.AddInt32(&instance.numOfRunningRoutines, -1)
		}()

		// run the job
		job()
	}()
	return ticket
}

// Wait all jobs are executed, if any in queue
func (instance *ConcurrentPool) Join() {
	for {
		time.Sleep(100 * time.Millisecond)
		if instance.numOfRunningRoutines == 0 {
			break
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
