// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.

package sites

// ===========================================================================

// ===========================================================================
// Beg of SiteMake creators

// SiteMakeChan returns a new open channel
// (simply a 'chan Site' that is).
// Note: No 'Site-producer' is launched here yet! (as is in all the other functions).
//  This is useful to easily create corresponding variables such as:
/*
var mySitePipelineStartsHere := SiteMakeChan()
// ... lot's of code to design and build Your favourite "mySiteWorkflowPipeline"
   // ...
   // ... *before* You start pouring data into it, e.g. simply via:
   for drop := range water {
mySitePipelineStartsHere <- drop
   }
close(mySitePipelineStartsHere)
*/
//  Hint: especially helpful, if Your piping library operates on some hidden (non-exported) type
//  (or on a type imported from elsewhere - and You don't want/need or should(!) have to care.)
//
// Note: as always (except for SitePipeBuffer) the channel is unbuffered.
//
func (my *Traffic) SiteMakeChan() (out chan Site) {
	return make(chan Site)
}

// End of SiteMake creators
// ===========================================================================

// ===========================================================================
// Beg of SiteChan producers

// SiteChan returns a channel to receive
// all inputs
// before close.
func (my *Traffic) SiteChan(inp ...Site) (out <-chan Site) {
	cha := make(chan Site)
	go my.chanSite(cha, inp...)
	return cha
}

func (my *Traffic) chanSite(out chan<- Site, inp ...Site) {
	defer close(out)
	for i := range inp {
		out <- inp[i]
	}
}

// SiteChanSlice returns a channel to receive
// all inputs
// before close.
func (my *Traffic) SiteChanSlice(inp ...[]Site) (out <-chan Site) {
	cha := make(chan Site)
	go my.chanSiteSlice(cha, inp...)
	return cha
}

func (my *Traffic) chanSiteSlice(out chan<- Site, inp ...[]Site) {
	defer close(out)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
}

// SiteChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func (my *Traffic) SiteChanFuncNok(gen func() (Site, bool)) (out <-chan Site) {
	cha := make(chan Site)
	go my.chanSiteFuncNok(cha, gen)
	return cha
}

func (my *Traffic) chanSiteFuncNok(out chan<- Site, gen func() (Site, bool)) {
	defer close(out)
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out <- res
	}
}

// SiteChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func (my *Traffic) SiteChanFuncErr(gen func() (Site, error)) (out <-chan Site) {
	cha := make(chan Site)
	go my.chanSiteFuncErr(cha, gen)
	return cha
}

func (my *Traffic) chanSiteFuncErr(out chan<- Site, gen func() (Site, error)) {
	defer close(out)
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out <- res
	}
}

// End of SiteChan producers
// ===========================================================================

// ===========================================================================
// Beg of SitePipe functions

// SitePipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be PipeSiteMap for functional people,
// but 'map' has a very different meaning in go lang.
func (my *Traffic) SitePipeFunc(inp <-chan Site, act func(a Site) Site) (out <-chan Site) {
	cha := make(chan Site)
	if act == nil { // Make `nil` value useful
		act = func(a Site) Site { return a }
	}
	go my.pipeSiteFunc(cha, inp, act)
	return cha
}

func (my *Traffic) pipeSiteFunc(out chan<- Site, inp <-chan Site, act func(a Site) Site) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of SitePipe functions
// ===========================================================================

// ===========================================================================
// Beg of SiteTube closures around SitePipe

// SiteTubeFunc returns a closure around PipeSiteFunc (_, act).
func (my *Traffic) SiteTubeFunc(act func(a Site) Site) (tube func(inp <-chan Site) (out <-chan Site)) {

	return func(inp <-chan Site) (out <-chan Site) {
		return my.SitePipeFunc(inp, act)
	}
}

// End of SiteTube closures around SitePipe
// ===========================================================================

// ===========================================================================
// Beg of SiteDone terminators

// SiteDone returns a channel to receive
// one signal before close after `inp` has been drained.
func (my *Traffic) SiteDone(inp <-chan Site) (done <-chan struct{}) {
	sig := make(chan struct{})
	go my.doneSite(sig, inp)
	return sig
}

func (my *Traffic) doneSite(done chan<- struct{}, inp <-chan Site) {
	defer close(done)
	for i := range inp {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// SiteDoneSlice returns a channel to receive
// a slice with every Site received on `inp`
// before close.
//
// Note: Unlike SiteDone, DoneSiteSlice sends the fully accumulated slice, not just an event, once upon close of inp.
func (my *Traffic) SiteDoneSlice(inp <-chan Site) (done <-chan []Site) {
	sig := make(chan []Site)
	go my.doneSiteSlice(sig, inp)
	return sig
}

func (my *Traffic) doneSiteSlice(done chan<- []Site, inp <-chan Site) {
	defer close(done)
	slice := []Site{}
	for i := range inp {
		slice = append(slice, i)
	}
	done <- slice
}

// SiteDoneFunc returns a channel to receive
// one signal after `act` has been applied to every `inp`
// before close.
func (my *Traffic) SiteDoneFunc(inp <-chan Site, act func(a Site)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a Site) { return }
	}
	go my.doneSiteFunc(sig, inp, act)
	return sig
}

func (my *Traffic) doneSiteFunc(done chan<- struct{}, inp <-chan Site, act func(a Site)) {
	defer close(done)
	for i := range inp {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of SiteDone terminators
// ===========================================================================

// ===========================================================================
// Beg of SiteFini closures

// SiteFini returns a closure around `SiteDone(_)`.
func (my *Traffic) SiteFini() func(inp <-chan Site) (done <-chan struct{}) {

	return func(inp <-chan Site) (done <-chan struct{}) {
		return my.SiteDone(inp)
	}
}

// SiteFiniSlice returns a closure around `SiteDoneSlice(_)`.
func (my *Traffic) SiteFiniSlice() func(inp <-chan Site) (done <-chan []Site) {

	return func(inp <-chan Site) (done <-chan []Site) {
		return my.SiteDoneSlice(inp)
	}
}

// SiteFiniFunc returns a closure around `SiteDoneFunc(_, act)`.
func (my *Traffic) SiteFiniFunc(act func(a Site)) func(inp <-chan Site) (done <-chan struct{}) {

	return func(inp <-chan Site) (done <-chan struct{}) {
		return my.SiteDoneFunc(inp, act)
	}
}

// End of SiteFini closures
// ===========================================================================

// ===========================================================================
// Beg of SitePair functions

// SitePair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func (my *Traffic) SitePair(inp <-chan Site) (out1, out2 <-chan Site) {
	cha1 := make(chan Site)
	cha2 := make(chan Site)
	go my.pairSite(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func (my *Traffic) pairSite(out1, out2 chan<- Site, inp <-chan Site) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func (my *Traffic) pairSite(out1, out2 chan<- Site, inp <-chan Site) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		select { // send first to whomever is ready to receive
		case out1 <- i:
			out2 <- i
		case out2 <- i:
			out1 <- i
		}
	}
}

// End of SitePair functions
// ===========================================================================

// ===========================================================================
// Beg of SiteFork functions

// SiteFork returns two channels
// either of which is to receive
// every result of inp
// before close.
func (my *Traffic) SiteFork(inp <-chan Site) (out1, out2 <-chan Site) {
	cha1 := make(chan Site)
	cha2 := make(chan Site)
	go my.forkSite(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func (my *Traffic) forkSite(out1, out2 chan<- Site, inp <-chan Site) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func (my *Traffic) forkSite(out1, out2 chan<- Site, inp <-chan Site) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		select { // send first to whomever is ready to receive
		case out1 <- i:
			out2 <- i
		case out2 <- i:
			out1 <- i
		}
	}
}

// End of SiteFork functions
// ===========================================================================

// ===========================================================================
// Beg of SiteFanIn2 simple binary Fan-In

// SiteFanIn2 returns a channel to receive all to receive all from both `inp1` and `inp2` before close.
func (my *Traffic) SiteFanIn2(inp1, inp2 <-chan Site) (out <-chan Site) {
	cha := make(chan Site)
	go my.fanIn2Site(cha, inp1, inp2)
	return cha
}

/* not used - kept for reference only.
// (my *Traffic) fanin2Site as seen in Go Concurrency Patterns
func fanin2Site(out chan<- Site, inp1, inp2 <-chan Site) {
	for {
		select {
		case e := <-inp1:
			out <- e
		case e := <-inp2:
			out <- e
		}
	}
} */

func (my *Traffic) fanIn2Site(out chan<- Site, inp1, inp2 <-chan Site) {
	defer close(out)

	var (
		closed bool // we found a chan closed
		ok     bool // did we read successfully?
		e      Site // what we've read
	)

	for !closed {
		select {
		case e, ok = <-inp1:
			if ok {
				out <- e
			} else {
				inp1 = inp2   // swap inp2 into inp1
				closed = true // break out of the loop
			}
		case e, ok = <-inp2:
			if ok {
				out <- e
			} else {
				closed = true // break out of the loop				}
			}
		}
	}

	// inp1 might not be closed yet. Drain it.
	for e = range inp1 {
		out <- e
	}
}

// End of SiteFanIn2 simple binary Fan-In
