// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.

package pipe

import (
	"container/ring"
	"sync"
	"time"

	"github.com/cheekybits/genny/generic"
)

// Thing is the generic type flowing thru the pipe network.
type Thing generic.Type

// ===========================================================================
// Beg of ThingMake creators

// ThingMakeChan returns a new open channel
// (simply a 'chan Thing' that is).
// Note: No 'Thing-producer' is launched here yet! (as is in all the other functions).
//  This is useful to easily create corresponding variables such as:
/*
var myThingPipelineStartsHere := ThingMakeChan()
// ... lot's of code to design and build Your favourite "myThingWorkflowPipeline"
   // ...
   // ... *before* You start pouring data into it, e.g. simply via:
   for drop := range water {
myThingPipelineStartsHere <- drop
   }
close(myThingPipelineStartsHere)
*/
//  Hint: especially helpful, if Your piping library operates on some hidden (non-exported) type
//  (or on a type imported from elsewhere - and You don't want/need or should(!) have to care.)
//
// Note: as always (except for ThingPipeBuffer) the channel is unbuffered.
//
func ThingMakeChan() (out chan Thing) {
	return make(chan Thing)
}

// End of ThingMake creators
// ===========================================================================

// ===========================================================================
// Beg of ThingChan producers

// ThingChan returns a channel to receive
// all inputs
// before close.
func ThingChan(inp ...Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go chanThing(cha, inp...)
	return cha
}

func chanThing(out chan<- Thing, inp ...Thing) {
	defer close(out)
	for i := range inp {
		out <- inp[i]
	}
}

// ThingChanSlice returns a channel to receive
// all inputs
// before close.
func ThingChanSlice(inp ...[]Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go chanThingSlice(cha, inp...)
	return cha
}

func chanThingSlice(out chan<- Thing, inp ...[]Thing) {
	defer close(out)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
}

// ThingChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func ThingChanFuncNok(gen func() (Thing, bool)) (out <-chan Thing) {
	cha := make(chan Thing)
	go chanThingFuncNok(cha, gen)
	return cha
}

func chanThingFuncNok(out chan<- Thing, gen func() (Thing, bool)) {
	defer close(out)
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out <- res
	}
}

// ThingChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func ThingChanFuncErr(gen func() (Thing, error)) (out <-chan Thing) {
	cha := make(chan Thing)
	go chanThingFuncErr(cha, gen)
	return cha
}

func chanThingFuncErr(out chan<- Thing, gen func() (Thing, error)) {
	defer close(out)
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out <- res
	}
}

// End of ThingChan producers
// ===========================================================================

// ===========================================================================
// Beg of ThingPipe functions

// ThingPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be PipeThingMap for functional people,
// but 'map' has a very different meaning in go lang.
func ThingPipeFunc(inp <-chan Thing, act func(a Thing) Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	if act == nil { // Make `nil` value useful
		act = func(a Thing) Thing { return a }
	}
	go pipeThingFunc(cha, inp, act)
	return cha
}

func pipeThingFunc(out chan<- Thing, inp <-chan Thing, act func(a Thing) Thing) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of ThingPipe functions
// ===========================================================================

// ===========================================================================
// Beg of ThingTube closures around ThingPipe

// ThingTubeFunc returns a closure around PipeThingFunc (_, act).
func ThingTubeFunc(act func(a Thing) Thing) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeFunc(inp, act)
	}
}

// End of ThingTube closures around ThingPipe
// ===========================================================================

// ===========================================================================
// Beg of ThingDone terminators

// ThingDone returns a channel to receive
// one signal before close after `inp` has been drained.
func ThingDone(inp <-chan Thing) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doneThing(sig, inp)
	return sig
}

func doneThing(done chan<- struct{}, inp <-chan Thing) {
	defer close(done)
	for i := range inp {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// ThingDoneSlice returns a channel to receive
// a slice with every Thing received on `inp`
// before close.
//
// Note: Unlike ThingDone, DoneThingSlice sends the fully accumulated slice, not just an event, once upon close of inp.
func ThingDoneSlice(inp <-chan Thing) (done <-chan []Thing) {
	sig := make(chan []Thing)
	go doneThingSlice(sig, inp)
	return sig
}

func doneThingSlice(done chan<- []Thing, inp <-chan Thing) {
	defer close(done)
	slice := []Thing{}
	for i := range inp {
		slice = append(slice, i)
	}
	done <- slice
}

// ThingDoneFunc returns a channel to receive
// one signal after `act` has been applied to every `inp`
// before close.
func ThingDoneFunc(inp <-chan Thing, act func(a Thing)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a Thing) { return }
	}
	go doneThingFunc(sig, inp, act)
	return sig
}

func doneThingFunc(done chan<- struct{}, inp <-chan Thing, act func(a Thing)) {
	defer close(done)
	for i := range inp {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of ThingDone terminators
// ===========================================================================

// ===========================================================================
// Beg of ThingFini closures

// ThingFini returns a closure around `ThingDone(_)`.
func ThingFini() func(inp <-chan Thing) (done <-chan struct{}) {

	return func(inp <-chan Thing) (done <-chan struct{}) {
		return ThingDone(inp)
	}
}

// ThingFiniSlice returns a closure around `ThingDoneSlice(_)`.
func ThingFiniSlice() func(inp <-chan Thing) (done <-chan []Thing) {

	return func(inp <-chan Thing) (done <-chan []Thing) {
		return ThingDoneSlice(inp)
	}
}

// ThingFiniFunc returns a closure around `ThingDoneFunc(_, act)`.
func ThingFiniFunc(act func(a Thing)) func(inp <-chan Thing) (done <-chan struct{}) {

	return func(inp <-chan Thing) (done <-chan struct{}) {
		return ThingDoneFunc(inp, act)
	}
}

// End of ThingFini closures
// ===========================================================================

// ===========================================================================
// Beg of ThingPair functions

// ThingPair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func ThingPair(inp <-chan Thing) (out1, out2 <-chan Thing) {
	cha1 := make(chan Thing)
	cha2 := make(chan Thing)
	go pairThing(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func pairThing(out1, out2 chan<- Thing, inp <-chan Thing) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func pairThing(out1, out2 chan<- Thing, inp <-chan Thing) {
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

// End of ThingPair functions
// ===========================================================================

// ===========================================================================
// Beg of ThingFork functions

// ThingFork returns two channels
// either of which is to receive
// every result of inp
// before close.
func ThingFork(inp <-chan Thing) (out1, out2 <-chan Thing) {
	cha1 := make(chan Thing)
	cha2 := make(chan Thing)
	go forkThing(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func forkThing(out1, out2 chan<- Thing, inp <-chan Thing) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func forkThing(out1, out2 chan<- Thing, inp <-chan Thing) {
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

// End of ThingFork functions
// ===========================================================================

// ===========================================================================
// Beg of ThingFanIn2 simple binary Fan-In

// ThingFanIn2 returns a channel to receive all to receive all from both `inp1` and `inp2` before close.
func ThingFanIn2(inp1, inp2 <-chan Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go fanIn2Thing(cha, inp1, inp2)
	return cha
}

/* not used - kept for reference only.
// fanin2Thing as seen in Go Concurrency Patterns
func fanin2Thing(out chan<- Thing, inp1, inp2 <-chan Thing) {
	for {
		select {
		case e := <-inp1:
			out <- e
		case e := <-inp2:
			out <- e
		}
	}
} */

func fanIn2Thing(out chan<- Thing, inp1, inp2 <-chan Thing) {
	defer close(out)

	var (
		closed bool  // we found a chan closed
		ok     bool  // did we read successfully?
		e      Thing // what we've read
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

// End of ThingFanIn2 simple binary Fan-In

// ===========================================================================
// Beg of ThingPipeBuffered - a buffered channel with capacity `cap` to receive

// ThingPipeBuffered returns a buffered channel with capacity `cap` to receive
// all `inp`
// before close.
func ThingPipeBuffered(inp <-chan Thing, cap int) (out <-chan Thing) {
	cha := make(chan Thing, cap)
	go pipeThingBuffered(cha, inp)
	return cha
}

func pipeThingBuffered(out chan<- Thing, inp <-chan Thing) {
	defer close(out)
	for i := range inp {
		out <- i
	}
}

// ThingTubeBuffered returns a closure around PipeThingBuffer (_, cap).
func ThingTubeBuffered(cap int) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeBuffered(inp, cap)
	}
}

// End of ThingPipeBuffered - a buffered channel with capacity `cap` to receive

// ===========================================================================
// Beg of ThingPipeEnter/Leave - Flapdoors observed by a Waiter

// ThingWaiter - as implemented by `*sync.WaitGroup` -
// attends Flapdoors and keeps counting
// who enters and who leaves.
//
// Use ThingDoneWait to learn about
// when the facilities are closed.
//
// Note: You may also use Your provided `*sync.WaitGroup.Wait()`
// to know when to close the facilities.
// Just: ThingDoneWait is more convenient
// as it also closes the primary channel for You.
//
// Just make sure to have _all_ entrances and exits attended,
// and `Wait()` only *after* You've started flooding the facilities.
type ThingWaiter interface {
	Add(delta int)
	Done()
	Wait()
}

// Note: The name is intentionally generic in order to avoid eventual multiple-declaration clashes.

// ThingPipeEnter returns a channel to receive
// all `inp`
// and registers throughput
// as arrival
// on the given `sync.WaitGroup`
// until close.
func ThingPipeEnter(inp <-chan Thing, wg ThingWaiter) (out <-chan Thing) {
	cha := make(chan Thing)
	go pipeThingEnter(cha, wg, inp)
	return cha
}

// ThingPipeLeave returns a channel to receive
// all `inp`
// and registers throughput
// as departure
// on the given `sync.WaitGroup`
// until close.
func ThingPipeLeave(inp <-chan Thing, wg ThingWaiter) (out <-chan Thing) {
	cha := make(chan Thing)
	go pipeThingLeave(cha, wg, inp)
	return cha
}

// ThingDoneLeave returns a channel to receive
// one signal after
// all throughput on `inp`
// has been registered
// as departure
// on the given `sync.WaitGroup`
// before close.
func ThingDoneLeave(inp <-chan Thing, wg ThingWaiter) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doneThingLeave(sig, wg, inp)
	return sig
}

func pipeThingEnter(out chan<- Thing, wg ThingWaiter, inp <-chan Thing) {
	defer close(out)
	for i := range inp {
		wg.Add(1)
		out <- i
	}
}

func pipeThingLeave(out chan<- Thing, wg ThingWaiter, inp <-chan Thing) {
	defer close(out)
	for i := range inp {
		out <- i
		wg.Done()
	}
}

func doneThingLeave(done chan<- struct{}, wg ThingWaiter, inp <-chan Thing) {
	defer close(done)
	for i := range inp {
		_ = i // discard
		wg.Done()
	}
	done <- struct{}{}
}

// ThingTubeEnter returns a closure around ThingPipeEnter (_, wg)
// registering throughput
// as arrival
// on the given `sync.WaitGroup`.
func ThingTubeEnter(wg ThingWaiter) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeEnter(inp, wg)
	}
}

// ThingTubeLeave returns a closure around ThingPipeLeave (_, wg)
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func ThingTubeLeave(wg ThingWaiter) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeLeave(inp, wg)
	}
}

// ThingFiniLeave returns a closure around `ThingDoneLeave(_, wg)`
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func ThingFiniLeave(wg ThingWaiter) func(inp <-chan Thing) (done <-chan struct{}) {

	return func(inp <-chan Thing) (done <-chan struct{}) {
		return ThingDoneLeave(inp, wg)
	}
}

// ThingDoneWait returns a channel to receive
// one signal
// after wg.Wait() has returned and inp has been closed
// before close.
//
// Note: Use only *after* You've started flooding the facilities.
func ThingDoneWait(inp chan<- Thing, wg ThingWaiter) (done <-chan struct{}) {
	cha := make(chan struct{})
	go doneThingWait(cha, inp, wg)
	return cha
}

func doneThingWait(done chan<- struct{}, inp chan<- Thing, wg ThingWaiter) {
	defer close(done)
	wg.Wait()
	close(inp)
	done <- struct{}{} // not really needed - but looks better
}

// ThingFiniWait returns a closure around `DoneThingWait(_, wg)`.
func ThingFiniWait(wg ThingWaiter) func(inp chan<- Thing) (done <-chan struct{}) {

	return func(inp chan<- Thing) (done <-chan struct{}) {
		return ThingDoneWait(inp, wg)
	}
}

// End of ThingPipeEnter/Leave - Flapdoors observed by a Waiter

// ===========================================================================
// Beg of ThingPipeDone

// ThingPipeDone returns a channel to receive every `inp` before close and a channel to signal this closing.
func ThingPipeDone(inp <-chan Thing) (out <-chan Thing, done <-chan struct{}) {
	cha := make(chan Thing)
	doit := make(chan struct{})
	go pipeThingDone(cha, doit, inp)
	return cha, doit
}

func pipeThingDone(out chan<- Thing, done chan<- struct{}, inp <-chan Thing) {
	defer close(out)
	defer close(done)
	for i := range inp {
		out <- i
	}
	done <- struct{}{}
}

// End of ThingPipeDone

// ===========================================================================
// Beg of ThingPlug - graceful terminator

// ThingPlug returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a stop signal,
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func ThingPlug(inp <-chan Thing, stop <-chan struct{}) (out <-chan Thing, done <-chan struct{}) {
	cha := make(chan Thing)
	doit := make(chan struct{})
	go plugThing(cha, doit, inp, stop)
	return cha, doit
}

func plugThing(out chan<- Thing, done chan<- struct{}, inp <-chan Thing, stop <-chan struct{}) {
	defer close(done)

	var end bool // shall we end?
	var ok bool  // did we read successfully?
	var e Thing  // what we've read

	for !end {
		select {
		case e, ok = <-inp:
			if ok {
				out <- e
			} else {
				end = true
			}
		case <-stop:
			end = true
		}
	}

	close(out)

	for range inp {
		// drain inp
	}

	done <- struct{}{}
}

// End of ThingPlug - graceful terminator

// ===========================================================================
// Beg of ThingPlugAfter - graceful terminator

// ThingPlugAfter returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a time signal
// (e.g. from `time.After(...)`),
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func ThingPlugAfter(inp <-chan Thing, after <-chan time.Time) (out <-chan Thing, done <-chan struct{}) {
	cha := make(chan Thing)
	doit := make(chan struct{})
	go plugThingAfter(cha, doit, inp, after)
	return cha, doit
}

func plugThingAfter(out chan<- Thing, done chan<- struct{}, inp <-chan Thing, after <-chan time.Time) {
	defer close(done)

	var end bool // shall we end?
	var ok bool  // did we read successfully?
	var e Thing  // what we've read

	for !end {
		select {
		case e, ok = <-inp:
			if ok {
				out <- e
			} else {
				end = true
			}
		case <-after:
			end = true
		}
	}

	close(out)

	for range inp {
		// drain inp
	}

	done <- struct{}{}
}

// End of ThingPlugAfter - graceful terminator

// Note: pipeThingAdjust imports "container/ring" for the expanding buffer.

// ===========================================================================
// Beg of ThingPipeAdjust

// ThingPipeAdjust returns a channel to receive
// all `inp`
// buffered by a ThingSendProxy process
// before close.
func ThingPipeAdjust(inp <-chan Thing, sizes ...int) (out <-chan Thing) {
	cap, que := sendThingProxySizes(sizes...)
	cha := make(chan Thing, cap)
	go pipeThingAdjust(cha, inp, que)
	return cha
}

// ThingTubeAdjust returns a closure around ThingPipeAdjust (_, sizes ...int).
func ThingTubeAdjust(sizes ...int) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeAdjust(inp, sizes...)
	}
}

// End of ThingPipeAdjust
// ===========================================================================

// ===========================================================================
// Beg of sendThingProxy

func sendThingProxySizes(sizes ...int) (cap, que int) {

	// CAP is the minimum capacity of the buffered proxy channel in `ThingSendProxy`
	const CAP = 10

	// QUE is the minimum initially allocated size of the circular queue in `ThingSendProxy`
	const QUE = 16

	cap = CAP
	que = QUE

	if len(sizes) > 0 && sizes[0] > CAP {
		que = sizes[0]
	}

	if len(sizes) > 1 && sizes[1] > QUE {
		que = sizes[1]
	}

	if len(sizes) > 2 {
		panic("ThingSendProxy: too many sizes")
	}

	return
}

// ThingSendProxy returns a channel to serve as a sending proxy to 'out'.
// Uses a goroutine to receive values from 'out' and store them
// in an expanding buffer, so that sending to 'out' never blocks.
//  Note: the expanding buffer is implemented via "container/ring"
//
// Note: ThingSendProxy is kept for the Sieve example
// and other dynamic use to be discovered
// even so it does not fit the pipe tube pattern as ThingPipeAdjust does.
func ThingSendProxy(out chan<- Thing, sizes ...int) chan<- Thing {
	cap, que := sendThingProxySizes(sizes...)
	cha := make(chan Thing, cap)
	go pipeThingAdjust(out, cha, que)
	return cha
}

// pipeThingAdjust uses an adjusting buffer to receive from 'inp'
// even so 'out' is not ready to receive yet. The buffer may grow
// until 'inp' is closed and then will shrink by every send to 'out'.
//  Note: the adjusting buffer is implemented via "container/ring"
func pipeThingAdjust(out chan<- Thing, inp <-chan Thing, QUE int) {
	defer close(out)
	n := QUE // the allocated size of the circular queue
	first := ring.New(n)
	last := first
	var c chan<- Thing
	var e Thing
	ok := true
	for ok {
		c = out
		if first == last {
			c = nil // buffer empty: disable output
		} else {
			e = first.Value.(Thing)
		}
		select {
		case e, ok = <-inp:
			if ok {
				last.Value = e
				if last.Next() == first {
					last.Link(ring.New(n)) // buffer full: expand it
					n *= 2
				}
				last = last.Next()
			}
		case c <- e:
			first = first.Next()
		}
	}

	for first != last {
		out <- first.Value.(Thing)
		first = first.Unlink(1) // first.Next()
	}
}

// End of sendThingProxy

// ===========================================================================
// Beg of ThingFanOut

// ThingFanOut returns a slice (of size = size) of channels
// each of which shall receive any inp before close.
func ThingFanOut(inp <-chan Thing, size int) (outS [](<-chan Thing)) {
	chaS := make([]chan Thing, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Thing)
	}

	go fanThingOut(inp, chaS...)

	outS = make([]<-chan Thing, size)
	for i := 0; i < size; i++ {
		outS[i] = (<-chan Thing)(chaS[i]) // convert `chan` to `<-chan`
	}

	return outS
}

// c fanThingOut(inp <-chan Thing, outs ...chan<- Thing) {
func fanThingOut(inp <-chan Thing, outs ...chan Thing) {

	for i := range inp {
		for o := range outs {
			outs[o] <- i
		}
	}

	for o := range outs {
		close(outs[o])
	}

}

// End of ThingFanOut

// ===========================================================================
// Beg of ThingStrew - scatter them

// ThingStrew returns a slice (of size = size) of channels
// one of which shall receive each inp before close.
func ThingStrew(inp <-chan Thing, size int) (outS [](<-chan Thing)) {
	chaS := make([]chan Thing, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Thing)
	}

	go strewThing(inp, chaS...)

	outS = make([]<-chan Thing, size)
	for i := 0; i < size; i++ {
		outS[i] = chaS[i] // convert `chan` to `<-chan`
	}

	return outS
}

// c strewThing(inp <-chan Thing, outS ...chan<- Thing) {
// Note: go does not convert the passed slice `[]chan Thing` to `[]chan<- Thing` automatically.
// So, we do neither here, as we are lazy (we just call an internal helper function).
func strewThing(inp <-chan Thing, outS ...chan Thing) {

	for i := range inp {
		for !trySendThing(i, outS...) {
			time.Sleep(time.Millisecond * 10) // wait a little before retry
		} // !sent
	} // inp

	for o := range outS {
		close(outS[o])
	}
}

func trySendThing(inp Thing, outS ...chan Thing) bool {

	for o := range outS {

		select { // try to send
		case outS[o] <- inp:
			return true
		default:
			// keep trying
		}

	} // outS
	return false
}

// End of ThingStrew - scatter them

// ===========================================================================
// Beg of ThingPipeSeen/ThingForkSeen - an "I've seen this Thing before" filter / forker

// ThingPipeSeen returns a channel to receive
// all `inp`
// not been seen before
// while silently dropping everything seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
// Note: ThingPipeFilterNotSeenYet might be a better name, but is fairly long.
func ThingPipeSeen(inp <-chan Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go pipeThingSeenAttr(cha, inp, nil)
	return cha
}

// ThingPipeSeenAttr returns a channel to receive
// all `inp`
// whose attribute `attr` has
// not been seen before
// while silently dropping everything seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
// Note: ThingPipeFilterAttrNotSeenYet might be a better name, but is fairly long.
func ThingPipeSeenAttr(inp <-chan Thing, attr func(a Thing) interface{}) (out <-chan Thing) {
	cha := make(chan Thing)
	go pipeThingSeenAttr(cha, inp, attr)
	return cha
}

// ThingForkSeen returns two channels, `new` and `old`,
// where `new` is to receive
// all `inp`
// not been seen before
// and `old`
// all `inp`
// seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
func ThingForkSeen(inp <-chan Thing) (new, old <-chan Thing) {
	cha1 := make(chan Thing)
	cha2 := make(chan Thing)
	go forkThingSeenAttr(cha1, cha2, inp, nil)
	return cha1, cha2
}

// ThingForkSeenAttr returns two channels, `new` and `old`,
// where `new` is to receive
// all `inp`
// whose attribute `attr` has
// not been seen before
// and `old`
// all `inp`
// seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
func ThingForkSeenAttr(inp <-chan Thing, attr func(a Thing) interface{}) (new, old <-chan Thing) {
	cha1 := make(chan Thing)
	cha2 := make(chan Thing)
	go forkThingSeenAttr(cha1, cha2, inp, attr)
	return cha1, cha2
}

func pipeThingSeenAttr(out chan<- Thing, inp <-chan Thing, attr func(a Thing) interface{}) {
	defer close(out)

	if attr == nil { // Make `nil` value useful
		attr = func(a Thing) interface{} { return a }
	}

	seen := sync.Map{}
	for i := range inp {
		if _, visited := seen.LoadOrStore(attr(i), struct{}{}); visited {
			// drop i silently
		} else {
			out <- i
		}
	}
}

func forkThingSeenAttr(new, old chan<- Thing, inp <-chan Thing, attr func(a Thing) interface{}) {
	defer close(new)
	defer close(old)

	if attr == nil { // Make `nil` value useful
		attr = func(a Thing) interface{} { return a }
	}

	seen := sync.Map{}
	for i := range inp {
		if _, visited := seen.LoadOrStore(attr(i), struct{}{}); visited {
			old <- i
		} else {
			new <- i
		}
	}
}

// ThingTubeSeen returns a closure around ThingPipeSeen()
// (silently dropping every Thing seen before).
func ThingTubeSeen() (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeSeen(inp)
	}
}

// ThingTubeSeenAttr returns a closure around ThingPipeSeenAttr()
// (silently dropping every Thing
// whose attribute `attr` was
// seen before).
func ThingTubeSeenAttr(attr func(a Thing) interface{}) (tube func(inp <-chan Thing) (out <-chan Thing)) {

	return func(inp <-chan Thing) (out <-chan Thing) {
		return ThingPipeSeenAttr(inp, attr)
	}
}

// End of ThingPipeSeen/ThingForkSeen - an "I've seen this Thing before" filter / forker

// ===========================================================================
// Beg of ThingFanIn

// ThingFanIn returns a channel to receive all inputs arriving
// on variadic inps
// before close.
//
//  Note: For each input one go routine is spawned to forward arrivals.
//
// See ThingFanIn1 in `fan-in1` for another implementation.
//
//  Ref: https://blog.golang.org/pipelines
//  Ref: https://github.com/QuentinPerez/go-stuff/channel/Fan-out-Fan-in/main.go
func ThingFanIn(inps ...<-chan Thing) (out <-chan Thing) {
	cha := make(chan Thing)

	wg := new(sync.WaitGroup)
	wg.Add(len(inps))

	go fanInThingWaitAndClose(cha, wg) // Spawn "close(out)" once all inps are done

	for i := range inps {
		go fanInThing(cha, inps[i], wg) // Spawn "output(c)"s
	}

	return cha
}

func fanInThing(out chan<- Thing, inp <-chan Thing, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range inp {
		out <- i
	}
}

func fanInThingWaitAndClose(out chan<- Thing, wg *sync.WaitGroup) {
	wg.Wait()
	close(out)
}

// End of ThingFanIn

// ===========================================================================
// Beg of ThingFanIn1 - fan-in using only one go routine

// ThingFanIn1 returns a channel to receive all inputs arriving
// on variadic inps
// before close.
//
//  Note: Only one go routine is used for all receives,
//  which keeps trying open inputs in round-robin fashion
//  until all inputs are closed.
//
// See ThingFanIn in `fan-in` for another implementation.
func ThingFanIn1(inpS ...<-chan Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go fanin1Thing(cha, inpS...)
	return cha
}

func fanin1Thing(out chan<- Thing, inpS ...<-chan Thing) {
	defer close(out)

	open := len(inpS)                 // assume: all are open
	closed := make([]bool, len(inpS)) // assume: each is not closed

	var item Thing // item received
	var ok bool    // receive channel is open?
	var sent bool  // some v has been sent?

	for open > 0 {
		sent = false
		for i := range inpS {
			if !closed[i] {
				select { // try to receive
				case item, ok = <-inpS[i]:
					if ok {
						out <- item
						sent = true
					} else {
						closed[i] = true
						open--
					}
				default: // keep going
				} // try
			} // not closed
		} // inpS
		if !sent && open > 0 {
			time.Sleep(time.Millisecond * 10) // wait a little before retry
		}
	} // open
}

// End of ThingFanIn1 - fan-in using only one go routine

// ===========================================================================
// Beg of ThingFan2 easy fan-in's

// ThingFan2 returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all inputs
// before close.
func ThingFan2(ori <-chan Thing, inp ...Thing) (out <-chan Thing) {
	return ThingFanIn2(ori, ThingChan(inp...))
}

// ThingFan2Slice returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all inputs
// before close.
func ThingFan2Slice(ori <-chan Thing, inp ...[]Thing) (out <-chan Thing) {
	return ThingFanIn2(ori, ThingChanSlice(inp...))
}

// ThingFan2Chan returns a channel to receive
// everything from the given original channel `ori`
// as well as
// from the the input channel `inp`
// before close.
// Note: ThingFan2Chan is nothing but ThingFanIn2
func ThingFan2Chan(ori <-chan Thing, inp <-chan Thing) (out <-chan Thing) {
	return ThingFanIn2(ori, inp)
}

// ThingFan2FuncNok returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all results of generator `gen`
// until `!ok`
// before close.
func ThingFan2FuncNok(ori <-chan Thing, gen func() (Thing, bool)) (out <-chan Thing) {
	return ThingFanIn2(ori, ThingChanFuncNok(gen))
}

// ThingFan2FuncErr returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all results of generator `gen`
// until `err != nil`
// before close.
func ThingFan2FuncErr(ori <-chan Thing, gen func() (Thing, error)) (out <-chan Thing) {
	return ThingFanIn2(ori, ThingChanFuncErr(gen))
}

// End of ThingFan2 easy fan-in's

// ===========================================================================
// Beg of ThingMerge

// ThingMerge returns a channel to receive all inputs sorted and free of duplicates.
// Each input channel needs to be sorted ascending and free of duplicates.
// The passed binary boolean function `less` defines the applicable order.
//  Note: If no inputs are given, a closed channel is returned.
func ThingMerge(less func(i, j Thing) bool, inps ...<-chan Thing) (out <-chan Thing) {

	if len(inps) < 1 { // none: return a closed channel
		cha := make(chan Thing)
		defer close(cha)
		return cha
	} else if len(inps) < 2 { // just one: return it
		return inps[0]
	} else { // tail recurse
		return mergeThing(less, inps[0], ThingMerge(less, inps[1:]...))
	}
}

// mergeThing takes two (eager) channels of comparable types,
// each of which needs to be sorted ascending and free of duplicates,
// and merges them into the returned channel, which will be sorted ascending and free of duplicates.
func mergeThing(less func(i, j Thing) bool, i1, i2 <-chan Thing) (out <-chan Thing) {
	cha := make(chan Thing)
	go func(out chan<- Thing, i1, i2 <-chan Thing) {
		defer close(out)
		var (
			clos1, clos2 bool  // we found the chan closed
			buff1, buff2 bool  // we've read 'from', but not sent (yet)
			ok           bool  // did we read successfully?
			from1, from2 Thing // what we've read
		)

		for !clos1 || !clos2 {

			if !clos1 && !buff1 {
				if from1, ok = <-i1; ok {
					buff1 = true
				} else {
					clos1 = true
				}
			}

			if !clos2 && !buff2 {
				if from2, ok = <-i2; ok {
					buff2 = true
				} else {
					clos2 = true
				}
			}

			if clos1 && !buff1 {
				from1 = from2
			}
			if clos2 && !buff2 {
				from2 = from1
			}

			if less(from1, from2) {
				out <- from1
				buff1 = false
			} else if less(from2, from1) {
				out <- from2
				buff2 = false
			} else {
				out <- from1 // == from2
				buff1 = false
				buff2 = false
			}
		}
	}(cha, i1, i2)
	return cha
}

// Note: mergeThing is not my own.
// Just: I forgot where found the original merge2 - please accept my apologies.
// I'd love to learn about it's origin/author, so I can give credit.
// Thus: Your hint, dear reader, is highly appreciated!

// End of ThingMerge

// ===========================================================================
// Beg of ThingSame comparator

// inspired by go/doc/play/tree.go

// ThingSame reads values from two channels in lockstep
// and iff they have the same contents then
// `true` is sent on the returned bool channel
// before close.
func ThingSame(same func(a, b Thing) bool, inp1, inp2 <-chan Thing) (out <-chan bool) {
	cha := make(chan bool)
	go sameThing(cha, same, inp1, inp2)
	return cha
}

func sameThing(out chan<- bool, same func(a, b Thing) bool, inp1, inp2 <-chan Thing) {
	defer close(out)
	for {
		v1, ok1 := <-inp1
		v2, ok2 := <-inp2

		if !ok1 || !ok2 {
			out <- ok1 == ok2
			return
		}
		if !same(v1, v2) {
			out <- false
			return
		}
	}
}

// End of ThingSame comparator

// ===========================================================================
// Beg of ThingJoin feedback back-feeders for circular networks

// ThingJoin sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func ThingJoin(out chan<- Thing, inp ...Thing) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinThing(sig, out, inp...)
	return sig
}

func joinThing(done chan<- struct{}, out chan<- Thing, inp ...Thing) {
	defer close(done)
	for i := range inp {
		out <- inp[i]
	}
	done <- struct{}{}
}

// ThingJoinSlice sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func ThingJoinSlice(out chan<- Thing, inp ...[]Thing) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinThingSlice(sig, out, inp...)
	return sig
}

func joinThingSlice(done chan<- struct{}, out chan<- Thing, inp ...[]Thing) {
	defer close(done)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
	done <- struct{}{}
}

// ThingJoinChan sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func ThingJoinChan(out chan<- Thing, inp <-chan Thing) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinThingChan(sig, out, inp)
	return sig
}

func joinThingChan(done chan<- struct{}, out chan<- Thing, inp <-chan Thing) {
	defer close(done)
	for i := range inp {
		out <- i
	}
	done <- struct{}{}
}

// End of ThingJoin feedback back-feeders for circular networks

// ===========================================================================
// Beg of ThingDaisyChain

// ThingProc is the signature of the inner process of any linear pipe-network
//  Example: the identity proc:
// samesame := func(into chan<- Thing, from <-chan Thing) { into <- <-from }
// Note: type ThingProc is provided for documentation purpose only.
// The implementation uses the explicit function signature
// in order to avoid some genny-related issue.
//  Note: In https://talks.golang.org/2012/waza.slide#40
// Rob Pike uses a ThingProc named `worker`.
type ThingProc func(into chan<- Thing, from <-chan Thing)

// Example: the identity proc - see `samesame` below
var _ ThingProc = func(out chan<- Thing, inp <-chan Thing) {
	// `out <- <-inp` or `into <- <-from`
	defer close(out)
	for i := range inp {
		out <- i
	}
}

// daisyThing returns a channel to receive all inp after having passed thru process `proc`.
func daisyThing(
	inp <-chan Thing, // a daisy to be chained
	proc func(into chan<- Thing, from <-chan Thing), // a process function
) (
	out chan Thing, // to receive all results
) { //  Body:

	cha := make(chan Thing)
	go proc(cha, inp)
	return cha
}

// ThingDaisyChain returns a channel to receive all inp
// after having passed
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making a null value useful.
func ThingDaisyChain(
	inp chan Thing, // a daisy to be chained
	procs ...func(out chan<- Thing, inp <-chan Thing), // a process function
) (
	out chan Thing, // to receive all results
) { //  Body:

	cha := inp

	if len(procs) < 1 {
		samesame := func(out chan<- Thing, inp <-chan Thing) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyThing(cha, samesame)
	} else {
		for _, proc := range procs {
			cha = daisyThing(cha, proc)
		}
	}
	return cha
}

// ThingDaisyChaiN returns a channel to receive all inp
// after having passed
// `somany` times
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If `somany` is less than 1 or no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making null values useful.
//
// Note: ThingDaisyChaiN(inp, 1, procs) <==> ThingDaisyChain(inp, procs)
func ThingDaisyChaiN(
	inp chan Thing, // a daisy to be chained
	somany int, // how many times? so many times
	procs ...func(out chan<- Thing, inp <-chan Thing), // a process function
) (
	out chan Thing, // to receive all results
) { //  Body:

	cha := inp

	if somany < 1 {
		samesame := func(out chan<- Thing, inp <-chan Thing) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyThing(cha, samesame)
	} else {
		for i := 0; i < somany; i++ {
			cha = ThingDaisyChain(cha, procs...)
		}
	}
	return cha
}

// End of ThingDaisyChain

// This file uses geanny to pull the type specific generic code
