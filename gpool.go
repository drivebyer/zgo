package zgo

import (
	"sync/atomic"

)

// This file support goroutine pool.
//
// Problem:
// When we create a groutine, wo use it for our service function.
// But once the goroutine is work done, it will back to runtime.
// Runtime would cache it or free it, it's out of our control.
// In our system, we make one goroutine to read from TCP connection,
// another goroutine to write to the connection. There would be many
// 'make' and 'destory'. It would lead our system slower.
//
// Solve:
// When goroutine done is's work, We can cache the goroutine on
// the channel's sudug queue instead of exiting it. When we need goroutine,
// we can ready from the queue.

type gpool struct {
	q     chan struct{} // waiting channel
	count int32         // waiting goroutine count
	size  int32         // pool size

	work chan interface{} // When goroutine waked up, there must be a new work.
}

func MakeGPool(size int32) *gpool {
	gp := &gpool{}
	gp.q = make(chan struct{})
	gp.work = make(chan interface{}, 1)
	gp.size = size
	return gp
}

// Put make the getg() block on the channel's sudog send queue.
func (gp *gpool) Put() (interface{}, bool) {
	if atomic.LoadInt32(&gp.count) == gp.size {
		// If the gpool is full, do nothing.
		return false, false
	}
	atomic.AddInt32(&gp.count, 1)
	//log.Printf("Put a goroutine in gpool, count:%d", atomic.LoadInt32(&gp.count))
	gp.q <- struct{}{}
	work := <-gp.work // When the goroutine wake up, it must be given a new work.
	return work, true
}

// Get wake a goroutine which wait on the channel's sudog send queue.
// When we wake a goroutine, we must give it a new work.
func (gp *gpool) Get(work interface{}) bool {
	if atomic.LoadInt32(&gp.count) == 0 {
		return false
	}
	atomic.AddInt32(&gp.count, -1)
	//log.Printf("Get a goroutine in gpool, count:%d", atomic.LoadInt32(&gp.count))
	<-gp.q
	gp.work <- work
	return true
}
