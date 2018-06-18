// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import (
	"github.com/cheekybits/genny/generic"
)

// Any is the generic type flowing thru the pipe network.
type Any generic.Type

// ===========================================================================
// Beg of PipeAnyEnter/Leave - Flapdoors observed by a Waiter

// AnyWaiter - as implemented by `*sync.WaitGroup` -
// attends Flapdoors and keeps track of
// how many enter and how many leave.
//
// Use Your provided `*sync.WaitGroup.Wait()`
// to know when to close the facilities.
//
// Just make sure to have _all_ entrances and exits attended,
// and don't `wg.Wait()` before You've flooded the facilities.
type AnyWaiter interface {
	Add(delta int)
	Done()
	// Wait() // here no need
}

// Note: Name is generic in order to avoid multiple-declaration clashes.

// PipeAnyEnter returns a channel to receive
// all `inp`
// and registers throughput
// as arrival
// on the given `sync.WaitGroup`
// until close.
func PipeAnyEnter(inp <-chan Any, wg AnyWaiter) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnyEnter(cha, wg, inp)
	return cha
}

// PipeAnyLeave returns a channel to receive
// all `inp`
// and registers throughput
// as departure
// on the given `sync.WaitGroup`
// until close.
func PipeAnyLeave(inp <-chan Any, wg AnyWaiter) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnyLeave(cha, wg, inp)
	return cha
}

func pipeAnyEnter(out chan<- Any, wg AnyWaiter, inp <-chan Any) {
	defer close(out)
	for i := range inp {
		wg.Add(1)
		out <- i
	}
}

func pipeAnyLeave(out chan<- Any, wg AnyWaiter, inp <-chan Any) {
	defer close(out)
	for i := range inp {
		out <- i
		wg.Done()
	}
}

// TubeAnyEnter returns a closure around PipeAnyEnter (_, wg)
// registering throughput
// on the given `sync.WaitGroup`
// as arrival.
func TubeAnyEnter(wg AnyWaiter) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return PipeAnyEnter(inp, wg)
	}
}

// TubeAnyLeave returns a closure around PipeAnyLeave (_, wg)
// registering throughput
// on the given `sync.WaitGroup`
// as departure.
func TubeAnyLeave(wg AnyWaiter) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return PipeAnyLeave(inp, wg)
	}
}

// End of PipeAnyEnter/Leave - Flapdoors observed by a Waiter
// ===========================================================================
