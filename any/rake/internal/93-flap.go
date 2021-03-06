// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rake

// ===========================================================================
// Beg of itemPipeEnter/Leave - Flapdoors observed by a Waiter

// itemWaiter - as implemented by `*sync.WaitGroup` -
// attends Flapdoors and keeps counting
// who enters and who leaves.
//
// Use itemDoneWait to learn about
// when the facilities are closed.
//
// Note: You may also use Your provided `*sync.WaitGroup.Wait()`
// to know when to close the facilities.
// Just: itemDoneWait is more convenient
// as it also closes the primary channel for You.
//
// Just make sure to have _all_ entrances and exits attended,
// and `Wait()` only *after* You've started flooding the facilities.
type itemWaiter interface {
	Add(delta int)
	Done()
	Wait()
}

// Note: The name is intentionally generic in order to avoid eventual multiple-declaration clashes.

// itemPipeEnter returns a channel to receive
// all `inp`
// and registers throughput
// as arrival
// on the given `sync.WaitGroup`
// until close.
func (inp itemFrom) itemPipeEnter(wg itemWaiter) (out itemFrom) {
	cha := make(chan item)
	go inp.pipeitemEnter(cha, wg)
	return cha
}

// itemPipeLeave returns a channel to receive
// all `inp`
// and registers throughput
// as departure
// on the given `sync.WaitGroup`
// until close.
func (inp itemFrom) itemPipeLeave(wg itemWaiter) (out itemFrom) {
	cha := make(chan item)
	go inp.pipeitemLeave(cha, wg)
	return cha
}

// itemDoneLeave returns a channel to receive
// one signal after
// all throughput on `inp`
// has been registered
// as departure
// on the given `sync.WaitGroup`
// before close.
func (inp itemFrom) itemDoneLeave(wg itemWaiter) (done <-chan struct{}) {
	sig := make(chan struct{})
	go inp.doneitemLeave(sig, wg)
	return sig
}

func (inp itemFrom) pipeitemEnter(out itemInto, wg itemWaiter) {
	defer close(out)
	for i := range inp {
		wg.Add(1)
		out <- i
	}
}

func (inp itemFrom) pipeitemLeave(out itemInto, wg itemWaiter) {
	defer close(out)
	for i := range inp {
		out <- i
		wg.Done()
	}
}

func (inp itemFrom) doneitemLeave(done chan<- struct{}, wg itemWaiter) {
	defer close(done)
	for i := range inp {
		_ = i // discard
		wg.Done()
	}
	done <- struct{}{}
}

// itemTubeEnter returns a closure around itemPipeEnter (wg)
// registering throughput
// as arrival
// on the given `sync.WaitGroup`.
func (inp itemFrom) itemTubeEnter(wg itemWaiter) (tube func(inp itemFrom) (out itemFrom)) {

	return func(inp itemFrom) (out itemFrom) {
		return inp.itemPipeEnter(wg)
	}
}

// itemTubeLeave returns a closure around itemPipeLeave (wg)
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func (inp itemFrom) itemTubeLeave(wg itemWaiter) (tube func(inp itemFrom) (out itemFrom)) {

	return func(inp itemFrom) (out itemFrom) {
		return inp.itemPipeLeave(wg)
	}
}

// itemFiniLeave returns a closure around `itemDoneLeave(wg)`
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func (inp itemFrom) itemFiniLeave(wg itemWaiter) func(inp itemFrom) (done <-chan struct{}) {

	return func(inp itemFrom) (done <-chan struct{}) {
		return inp.itemDoneLeave(wg)
	}
}

// itemDoneWait returns a channel to receive
// one signal
// after wg.Wait() has returned and out has been closed
// before close.
//
// Note: Use only *after* You've started flooding the facilities.
func (out itemInto) itemDoneWait(wg itemWaiter) (done <-chan struct{}) {
	cha := make(chan struct{})
	go out.doneitemWait(cha, wg)
	return cha
}

func (out itemInto) doneitemWait(done chan<- struct{}, wg itemWaiter) {
	defer close(done)
	wg.Wait()
	close(out)
	done <- struct{}{} // not really needed - but looks better
}

// itemFiniWait returns a closure around `itemDoneWait(wg)`.
func (out itemInto) itemFiniWait(wg itemWaiter) func(out itemInto) (done <-chan struct{}) {

	return func(out itemInto) (done <-chan struct{}) {
		return out.itemDoneWait(wg)
	}
}

// End of itemPipeEnter/Leave - Flapdoors observed by a Waiter
// ===========================================================================
