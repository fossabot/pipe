// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

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

// Any is the generic type flowing thru the pipe network.
type Any generic.Type

// ===========================================================================
// Beg of AnyMake creators

// AnyMakeChan returns a new open channel
// (simply a 'chan Any' that is).
// Note: No 'Any-producer' is launched here yet! (as is in all the other functions).
//  This is useful to easily create corresponding variables such as:
/*
var myAnyPipelineStartsHere := AnyMakeChan()
// ... lot's of code to design and build Your favourite "myAnyWorkflowPipeline"
   // ...
   // ... *before* You start pouring data into it, e.g. simply via:
   for drop := range water {
myAnyPipelineStartsHere <- drop
   }
close(myAnyPipelineStartsHere)
*/
//  Hint: especially helpful, if Your piping library operates on some hidden (non-exported) type
//  (or on a type imported from elsewhere - and You don't want/need or should(!) have to care.)
//
// Note: as always (except for AnyPipeBuffer) the channel is unbuffered.
//
func AnyMakeChan() (out chan Any) {
	return make(chan Any)
}

// End of AnyMake creators
// ===========================================================================

// ===========================================================================
// Beg of AnyChan producers

// AnyChan returns a channel to receive
// all inputs
// before close.
func AnyChan(inp ...Any) (out <-chan Any) {
	cha := make(chan Any)
	go chanAny(cha, inp...)
	return cha
}

func chanAny(out chan<- Any, inp ...Any) {
	defer close(out)
	for i := range inp {
		out <- inp[i]
	}
}

// AnyChanSlice returns a channel to receive
// all inputs
// before close.
func AnyChanSlice(inp ...[]Any) (out <-chan Any) {
	cha := make(chan Any)
	go chanAnySlice(cha, inp...)
	return cha
}

func chanAnySlice(out chan<- Any, inp ...[]Any) {
	defer close(out)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
}

// AnyChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func AnyChanFuncNok(gen func() (Any, bool)) (out <-chan Any) {
	cha := make(chan Any)
	go chanAnyFuncNok(cha, gen)
	return cha
}

func chanAnyFuncNok(out chan<- Any, gen func() (Any, bool)) {
	defer close(out)
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out <- res
	}
}

// AnyChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func AnyChanFuncErr(gen func() (Any, error)) (out <-chan Any) {
	cha := make(chan Any)
	go chanAnyFuncErr(cha, gen)
	return cha
}

func chanAnyFuncErr(out chan<- Any, gen func() (Any, error)) {
	defer close(out)
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out <- res
	}
}

// End of AnyChan producers
// ===========================================================================

// ===========================================================================
// Beg of AnyPipe functions

// AnyPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be PipeAnyMap for functional people,
// but 'map' has a very different meaning in go lang.
func AnyPipeFunc(inp <-chan Any, act func(a Any) Any) (out <-chan Any) {
	cha := make(chan Any)
	if act == nil { // Make `nil` value useful
		act = func(a Any) Any { return a }
	}
	go pipeAnyFunc(cha, inp, act)
	return cha
}

func pipeAnyFunc(out chan<- Any, inp <-chan Any, act func(a Any) Any) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of AnyPipe functions
// ===========================================================================

// ===========================================================================
// Beg of AnyTube closures around AnyPipe

// AnyTubeFunc returns a closure around PipeAnyFunc (_, act).
func AnyTubeFunc(act func(a Any) Any) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeFunc(inp, act)
	}
}

// End of AnyTube closures around AnyPipe
// ===========================================================================

// ===========================================================================
// Beg of AnyDone terminators

// AnyDone returns a channel to receive
// one signal before close after `inp` has been drained.
func AnyDone(inp <-chan Any) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doneAny(sig, inp)
	return sig
}

func doneAny(done chan<- struct{}, inp <-chan Any) {
	defer close(done)
	for i := range inp {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// AnyDoneSlice returns a channel to receive
// a slice with every Any received on `inp`
// before close.
//
// Note: Unlike AnyDone, DoneAnySlice sends the fully accumulated slice, not just an event, once upon close of inp.
func AnyDoneSlice(inp <-chan Any) (done <-chan []Any) {
	sig := make(chan []Any)
	go doneAnySlice(sig, inp)
	return sig
}

func doneAnySlice(done chan<- []Any, inp <-chan Any) {
	defer close(done)
	slice := []Any{}
	for i := range inp {
		slice = append(slice, i)
	}
	done <- slice
}

// AnyDoneFunc returns a channel to receive
// one signal after `act` has been applied to every `inp`
// before close.
func AnyDoneFunc(inp <-chan Any, act func(a Any)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a Any) { return }
	}
	go doneAnyFunc(sig, inp, act)
	return sig
}

func doneAnyFunc(done chan<- struct{}, inp <-chan Any, act func(a Any)) {
	defer close(done)
	for i := range inp {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of AnyDone terminators
// ===========================================================================

// ===========================================================================
// Beg of AnyFini closures

// AnyFini returns a closure around `AnyDone(_)`.
func AnyFini() func(inp <-chan Any) (done <-chan struct{}) {

	return func(inp <-chan Any) (done <-chan struct{}) {
		return AnyDone(inp)
	}
}

// AnyFiniSlice returns a closure around `AnyDoneSlice(_)`.
func AnyFiniSlice() func(inp <-chan Any) (done <-chan []Any) {

	return func(inp <-chan Any) (done <-chan []Any) {
		return AnyDoneSlice(inp)
	}
}

// AnyFiniFunc returns a closure around `AnyDoneFunc(_, act)`.
func AnyFiniFunc(act func(a Any)) func(inp <-chan Any) (done <-chan struct{}) {

	return func(inp <-chan Any) (done <-chan struct{}) {
		return AnyDoneFunc(inp, act)
	}
}

// End of AnyFini closures
// ===========================================================================

// ===========================================================================
// Beg of AnyPair functions

// AnyPair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func AnyPair(inp <-chan Any) (out1, out2 <-chan Any) {
	cha1 := make(chan Any)
	cha2 := make(chan Any)
	go pairAny(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func pairAny(out1, out2 chan<- Any, inp <-chan Any) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func pairAny(out1, out2 chan<- Any, inp <-chan Any) {
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

// End of AnyPair functions
// ===========================================================================

// ===========================================================================
// Beg of AnyFork functions

// AnyFork returns two channels
// either of which is to receive
// every result of inp
// before close.
func AnyFork(inp <-chan Any) (out1, out2 <-chan Any) {
	cha1 := make(chan Any)
	cha2 := make(chan Any)
	go forkAny(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func forkAny(out1, out2 chan<- Any, inp <-chan Any) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func forkAny(out1, out2 chan<- Any, inp <-chan Any) {
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

// End of AnyFork functions
// ===========================================================================

// ===========================================================================
// Beg of AnyFanIn2 simple binary Fan-In

// AnyFanIn2 returns a channel to receive all to receive all from both `inp1` and `inp2` before close.
func AnyFanIn2(inp1, inp2 <-chan Any) (out <-chan Any) {
	cha := make(chan Any)
	go fanIn2Any(cha, inp1, inp2)
	return cha
}

/* not used - kept for reference only.
// fanin2Any as seen in Go Concurrency Patterns
func fanin2Any(out chan<- Any, inp1, inp2 <-chan Any) {
	for {
		select {
		case e := <-inp1:
			out <- e
		case e := <-inp2:
			out <- e
		}
	}
} */

func fanIn2Any(out chan<- Any, inp1, inp2 <-chan Any) {
	defer close(out)

	var (
		closed bool // we found a chan closed
		ok     bool // did we read successfully?
		e      Any  // what we've read
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

// End of AnyFanIn2 simple binary Fan-In

// ===========================================================================
// Beg of AnyPipeBuffered - a buffered channel with capacity `cap` to receive

// AnyPipeBuffered returns a buffered channel with capacity `cap` to receive
// all `inp`
// before close.
func AnyPipeBuffered(inp <-chan Any, cap int) (out <-chan Any) {
	cha := make(chan Any, cap)
	go pipeAnyBuffered(cha, inp)
	return cha
}

func pipeAnyBuffered(out chan<- Any, inp <-chan Any) {
	defer close(out)
	for i := range inp {
		out <- i
	}
}

// AnyTubeBuffered returns a closure around PipeAnyBuffer (_, cap).
func AnyTubeBuffered(cap int) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeBuffered(inp, cap)
	}
}

// End of AnyPipeBuffered - a buffered channel with capacity `cap` to receive

// ===========================================================================
// Beg of AnyPipeEnter/Leave - Flapdoors observed by a Waiter

// AnyWaiter - as implemented by `*sync.WaitGroup` -
// attends Flapdoors and keeps counting
// who enters and who leaves.
//
// Use AnyDoneWait to learn about
// when the facilities are closed.
//
// Note: You may also use Your provided `*sync.WaitGroup.Wait()`
// to know when to close the facilities.
// Just: AnyDoneWait is more convenient
// as it also closes the primary channel for You.
//
// Just make sure to have _all_ entrances and exits attended,
// and `Wait()` only *after* You've started flooding the facilities.
type AnyWaiter interface {
	Add(delta int)
	Done()
	Wait()
}

// Note: The name is intentionally generic in order to avoid eventual multiple-declaration clashes.

// AnyPipeEnter returns a channel to receive
// all `inp`
// and registers throughput
// as arrival
// on the given `sync.WaitGroup`
// until close.
func AnyPipeEnter(inp <-chan Any, wg AnyWaiter) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnyEnter(cha, wg, inp)
	return cha
}

// AnyPipeLeave returns a channel to receive
// all `inp`
// and registers throughput
// as departure
// on the given `sync.WaitGroup`
// until close.
func AnyPipeLeave(inp <-chan Any, wg AnyWaiter) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnyLeave(cha, wg, inp)
	return cha
}

// AnyDoneLeave returns a channel to receive
// one signal after
// all throughput on `inp`
// has been registered
// as departure
// on the given `sync.WaitGroup`
// before close.
func AnyDoneLeave(inp <-chan Any, wg AnyWaiter) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doneAnyLeave(sig, wg, inp)
	return sig
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

func doneAnyLeave(done chan<- struct{}, wg AnyWaiter, inp <-chan Any) {
	defer close(done)
	for i := range inp {
		_ = i // discard
		wg.Done()
	}
	done <- struct{}{}
}

// AnyTubeEnter returns a closure around AnyPipeEnter (_, wg)
// registering throughput
// as arrival
// on the given `sync.WaitGroup`.
func AnyTubeEnter(wg AnyWaiter) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeEnter(inp, wg)
	}
}

// AnyTubeLeave returns a closure around AnyPipeLeave (_, wg)
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func AnyTubeLeave(wg AnyWaiter) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeLeave(inp, wg)
	}
}

// AnyFiniLeave returns a closure around `AnyDoneLeave(_, wg)`
// registering throughput
// as departure
// on the given `sync.WaitGroup`.
func AnyFiniLeave(wg AnyWaiter) func(inp <-chan Any) (done <-chan struct{}) {

	return func(inp <-chan Any) (done <-chan struct{}) {
		return AnyDoneLeave(inp, wg)
	}
}

// AnyDoneWait returns a channel to receive
// one signal
// after wg.Wait() has returned and inp has been closed
// before close.
//
// Note: Use only *after* You've started flooding the facilities.
func AnyDoneWait(inp chan<- Any, wg AnyWaiter) (done <-chan struct{}) {
	cha := make(chan struct{})
	go doneAnyWait(cha, inp, wg)
	return cha
}

func doneAnyWait(done chan<- struct{}, inp chan<- Any, wg AnyWaiter) {
	defer close(done)
	wg.Wait()
	close(inp)
	done <- struct{}{} // not really needed - but looks better
}

// AnyFiniWait returns a closure around `DoneAnyWait(_, wg)`.
func AnyFiniWait(wg AnyWaiter) func(inp chan<- Any) (done <-chan struct{}) {

	return func(inp chan<- Any) (done <-chan struct{}) {
		return AnyDoneWait(inp, wg)
	}
}

// End of AnyPipeEnter/Leave - Flapdoors observed by a Waiter

// ===========================================================================
// Beg of AnyPipeDone

// AnyPipeDone returns a channel to receive every `inp` before close and a channel to signal this closing.
func AnyPipeDone(inp <-chan Any) (out <-chan Any, done <-chan struct{}) {
	cha := make(chan Any)
	doit := make(chan struct{})
	go pipeAnyDone(cha, doit, inp)
	return cha, doit
}

func pipeAnyDone(out chan<- Any, done chan<- struct{}, inp <-chan Any) {
	defer close(out)
	defer close(done)
	for i := range inp {
		out <- i
	}
	done <- struct{}{}
}

// End of AnyPipeDone

// ===========================================================================
// Beg of AnyPlug - graceful terminator

// AnyPlug returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a stop signal,
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func AnyPlug(inp <-chan Any, stop <-chan struct{}) (out <-chan Any, done <-chan struct{}) {
	cha := make(chan Any)
	doit := make(chan struct{})
	go plugAny(cha, doit, inp, stop)
	return cha, doit
}

func plugAny(out chan<- Any, done chan<- struct{}, inp <-chan Any, stop <-chan struct{}) {
	defer close(done)

	var end bool // shall we end?
	var ok bool  // did we read successfully?
	var e Any    // what we've read

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

// End of AnyPlug - graceful terminator

// ===========================================================================
// Beg of AnyPlugAfter - graceful terminator

// AnyPlugAfter returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a time signal
// (e.g. from `time.After(...)`),
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func AnyPlugAfter(inp <-chan Any, after <-chan time.Time) (out <-chan Any, done <-chan struct{}) {
	cha := make(chan Any)
	doit := make(chan struct{})
	go plugAnyAfter(cha, doit, inp, after)
	return cha, doit
}

func plugAnyAfter(out chan<- Any, done chan<- struct{}, inp <-chan Any, after <-chan time.Time) {
	defer close(done)

	var end bool // shall we end?
	var ok bool  // did we read successfully?
	var e Any    // what we've read

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

// End of AnyPlugAfter - graceful terminator

// Note: pipeAnyAdjust imports "container/ring" for the expanding buffer.

// ===========================================================================
// Beg of AnyPipeAdjust

// AnyPipeAdjust returns a channel to receive
// all `inp`
// buffered by a AnySendProxy process
// before close.
func AnyPipeAdjust(inp <-chan Any, sizes ...int) (out <-chan Any) {
	cap, que := sendAnyProxySizes(sizes...)
	cha := make(chan Any, cap)
	go pipeAnyAdjust(cha, inp, que)
	return cha
}

// AnyTubeAdjust returns a closure around AnyPipeAdjust (_, sizes ...int).
func AnyTubeAdjust(sizes ...int) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeAdjust(inp, sizes...)
	}
}

// End of AnyPipeAdjust
// ===========================================================================

// ===========================================================================
// Beg of sendAnyProxy

func sendAnyProxySizes(sizes ...int) (cap, que int) {

	// CAP is the minimum capacity of the buffered proxy channel in `AnySendProxy`
	const CAP = 10

	// QUE is the minimum initially allocated size of the circular queue in `AnySendProxy`
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
		panic("AnySendProxy: too many sizes")
	}

	return
}

// AnySendProxy returns a channel to serve as a sending proxy to 'out'.
// Uses a goroutine to receive values from 'out' and store them
// in an expanding buffer, so that sending to 'out' never blocks.
//  Note: the expanding buffer is implemented via "container/ring"
//
// Note: AnySendProxy is kept for the Sieve example
// and other dynamic use to be discovered
// even so it does not fit the pipe tube pattern as AnyPipeAdjust does.
func AnySendProxy(out chan<- Any, sizes ...int) chan<- Any {
	cap, que := sendAnyProxySizes(sizes...)
	cha := make(chan Any, cap)
	go pipeAnyAdjust(out, cha, que)
	return cha
}

// pipeAnyAdjust uses an adjusting buffer to receive from 'inp'
// even so 'out' is not ready to receive yet. The buffer may grow
// until 'inp' is closed and then will shrink by every send to 'out'.
//  Note: the adjusting buffer is implemented via "container/ring"
func pipeAnyAdjust(out chan<- Any, inp <-chan Any, QUE int) {
	defer close(out)
	n := QUE // the allocated size of the circular queue
	first := ring.New(n)
	last := first
	var c chan<- Any
	var e Any
	ok := true
	for ok {
		c = out
		if first == last {
			c = nil // buffer empty: disable output
		} else {
			e = first.Value.(Any)
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
		out <- first.Value.(Any)
		first = first.Unlink(1) // first.Next()
	}
}

// End of sendAnyProxy

// ===========================================================================
// Beg of AnyFanOut

// AnyFanOut returns a slice (of size = size) of channels
// each of which shall receive any inp before close.
func AnyFanOut(inp <-chan Any, size int) (outS [](<-chan Any)) {
	chaS := make([]chan Any, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Any)
	}

	go fanAnyOut(inp, chaS...)

	outS = make([]<-chan Any, size)
	for i := 0; i < size; i++ {
		outS[i] = (<-chan Any)(chaS[i]) // convert `chan` to `<-chan`
	}

	return outS
}

// c fanAnyOut(inp <-chan Any, outs ...chan<- Any) {
func fanAnyOut(inp <-chan Any, outs ...chan Any) {

	for i := range inp {
		for o := range outs {
			outs[o] <- i
		}
	}

	for o := range outs {
		close(outs[o])
	}

}

// End of AnyFanOut

// ===========================================================================
// Beg of AnyStrew - scatter them

// AnyStrew returns a slice (of size = size) of channels
// one of which shall receive each inp before close.
func AnyStrew(inp <-chan Any, size int) (outS [](<-chan Any)) {
	chaS := make([]chan Any, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Any)
	}

	go strewAny(inp, chaS...)

	outS = make([]<-chan Any, size)
	for i := 0; i < size; i++ {
		outS[i] = chaS[i] // convert `chan` to `<-chan`
	}

	return outS
}

// c strewAny(inp <-chan Any, outS ...chan<- Any) {
// Note: go does not convert the passed slice `[]chan Any` to `[]chan<- Any` automatically.
// So, we do neither here, as we are lazy (we just call an internal helper function).
func strewAny(inp <-chan Any, outS ...chan Any) {

	for i := range inp {
		for !trySendAny(i, outS...) {
			time.Sleep(time.Millisecond * 10) // wait a little before retry
		} // !sent
	} // inp

	for o := range outS {
		close(outS[o])
	}
}

func trySendAny(inp Any, outS ...chan Any) bool {

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

// End of AnyStrew - scatter them

// ===========================================================================
// Beg of AnyPipeSeen/AnyForkSeen - an "I've seen this Any before" filter / forker

// AnyPipeSeen returns a channel to receive
// all `inp`
// not been seen before
// while silently dropping everything seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
// Note: AnyPipeFilterNotSeenYet might be a better name, but is fairly long.
func AnyPipeSeen(inp <-chan Any) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnySeenAttr(cha, inp, nil)
	return cha
}

// AnyPipeSeenAttr returns a channel to receive
// all `inp`
// whose attribute `attr` has
// not been seen before
// while silently dropping everything seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
// Note: AnyPipeFilterAttrNotSeenYet might be a better name, but is fairly long.
func AnyPipeSeenAttr(inp <-chan Any, attr func(a Any) interface{}) (out <-chan Any) {
	cha := make(chan Any)
	go pipeAnySeenAttr(cha, inp, attr)
	return cha
}

// AnyForkSeen returns two channels, `new` and `old`,
// where `new` is to receive
// all `inp`
// not been seen before
// and `old`
// all `inp`
// seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
func AnyForkSeen(inp <-chan Any) (new, old <-chan Any) {
	cha1 := make(chan Any)
	cha2 := make(chan Any)
	go forkAnySeenAttr(cha1, cha2, inp, nil)
	return cha1, cha2
}

// AnyForkSeenAttr returns two channels, `new` and `old`,
// where `new` is to receive
// all `inp`
// whose attribute `attr` has
// not been seen before
// and `old`
// all `inp`
// seen before
// (internally growing a `sync.Map` to discriminate)
// until close.
func AnyForkSeenAttr(inp <-chan Any, attr func(a Any) interface{}) (new, old <-chan Any) {
	cha1 := make(chan Any)
	cha2 := make(chan Any)
	go forkAnySeenAttr(cha1, cha2, inp, attr)
	return cha1, cha2
}

func pipeAnySeenAttr(out chan<- Any, inp <-chan Any, attr func(a Any) interface{}) {
	defer close(out)

	if attr == nil { // Make `nil` value useful
		attr = func(a Any) interface{} { return a }
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

func forkAnySeenAttr(new, old chan<- Any, inp <-chan Any, attr func(a Any) interface{}) {
	defer close(new)
	defer close(old)

	if attr == nil { // Make `nil` value useful
		attr = func(a Any) interface{} { return a }
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

// AnyTubeSeen returns a closure around AnyPipeSeen()
// (silently dropping every Any seen before).
func AnyTubeSeen() (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeSeen(inp)
	}
}

// AnyTubeSeenAttr returns a closure around AnyPipeSeenAttr()
// (silently dropping every Any
// whose attribute `attr` was
// seen before).
func AnyTubeSeenAttr(attr func(a Any) interface{}) (tube func(inp <-chan Any) (out <-chan Any)) {

	return func(inp <-chan Any) (out <-chan Any) {
		return AnyPipeSeenAttr(inp, attr)
	}
}

// End of AnyPipeSeen/AnyForkSeen - an "I've seen this Any before" filter / forker

// ===========================================================================
// Beg of AnyFanIn

// AnyFanIn returns a channel to receive all inputs arriving
// on variadic inps
// before close.
//
//  Note: For each input one go routine is spawned to forward arrivals.
//
// See AnyFanIn1 in `fan-in1` for another implementation.
//
//  Ref: https://blog.golang.org/pipelines
//  Ref: https://github.com/QuentinPerez/go-stuff/channel/Fan-out-Fan-in/main.go
func AnyFanIn(inps ...<-chan Any) (out <-chan Any) {
	cha := make(chan Any)

	wg := new(sync.WaitGroup)
	wg.Add(len(inps))

	go fanInAnyWaitAndClose(cha, wg) // Spawn "close(out)" once all inps are done

	for i := range inps {
		go fanInAny(cha, inps[i], wg) // Spawn "output(c)"s
	}

	return cha
}

func fanInAny(out chan<- Any, inp <-chan Any, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range inp {
		out <- i
	}
}

func fanInAnyWaitAndClose(out chan<- Any, wg *sync.WaitGroup) {
	wg.Wait()
	close(out)
}

// End of AnyFanIn

// ===========================================================================
// Beg of AnyFanIn1 - fan-in using only one go routine

// AnyFanIn1 returns a channel to receive all inputs arriving
// on variadic inps
// before close.
//
//  Note: Only one go routine is used for all receives,
//  which keeps trying open inputs in round-robin fashion
//  until all inputs are closed.
//
// See AnyFanIn in `fan-in` for another implementation.
func AnyFanIn1(inpS ...<-chan Any) (out <-chan Any) {
	cha := make(chan Any)
	go fanin1Any(cha, inpS...)
	return cha
}

func fanin1Any(out chan<- Any, inpS ...<-chan Any) {
	defer close(out)

	open := len(inpS)                 // assume: all are open
	closed := make([]bool, len(inpS)) // assume: each is not closed

	var item Any  // item received
	var ok bool   // receive channel is open?
	var sent bool // some v has been sent?

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

// End of AnyFanIn1 - fan-in using only one go routine

// ===========================================================================
// Beg of AnyFan2 easy fan-in's

// AnyFan2 returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all inputs
// before close.
func AnyFan2(ori <-chan Any, inp ...Any) (out <-chan Any) {
	return AnyFanIn2(ori, AnyChan(inp...))
}

// AnyFan2Slice returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all inputs
// before close.
func AnyFan2Slice(ori <-chan Any, inp ...[]Any) (out <-chan Any) {
	return AnyFanIn2(ori, AnyChanSlice(inp...))
}

// AnyFan2Chan returns a channel to receive
// everything from the given original channel `ori`
// as well as
// from the the input channel `inp`
// before close.
// Note: AnyFan2Chan is nothing but AnyFanIn2
func AnyFan2Chan(ori <-chan Any, inp <-chan Any) (out <-chan Any) {
	return AnyFanIn2(ori, inp)
}

// AnyFan2FuncNok returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all results of generator `gen`
// until `!ok`
// before close.
func AnyFan2FuncNok(ori <-chan Any, gen func() (Any, bool)) (out <-chan Any) {
	return AnyFanIn2(ori, AnyChanFuncNok(gen))
}

// AnyFan2FuncErr returns a channel to receive
// everything from the given original channel `ori`
// as well as
// all results of generator `gen`
// until `err != nil`
// before close.
func AnyFan2FuncErr(ori <-chan Any, gen func() (Any, error)) (out <-chan Any) {
	return AnyFanIn2(ori, AnyChanFuncErr(gen))
}

// End of AnyFan2 easy fan-in's

// ===========================================================================
// Beg of AnyMerge

// AnyMerge returns a channel to receive all inputs sorted and free of duplicates.
// Each input channel needs to be sorted ascending and free of duplicates.
// The passed binary boolean function `less` defines the applicable order.
//  Note: If no inputs are given, a closed channel is returned.
func AnyMerge(less func(i, j Any) bool, inps ...<-chan Any) (out <-chan Any) {

	if len(inps) < 1 { // none: return a closed channel
		cha := make(chan Any)
		defer close(cha)
		return cha
	} else if len(inps) < 2 { // just one: return it
		return inps[0]
	} else { // tail recurse
		return mergeAny(less, inps[0], AnyMerge(less, inps[1:]...))
	}
}

// mergeAny takes two (eager) channels of comparable types,
// each of which needs to be sorted ascending and free of duplicates,
// and merges them into the returned channel, which will be sorted ascending and free of duplicates.
func mergeAny(less func(i, j Any) bool, i1, i2 <-chan Any) (out <-chan Any) {
	cha := make(chan Any)
	go func(out chan<- Any, i1, i2 <-chan Any) {
		defer close(out)
		var (
			clos1, clos2 bool // we found the chan closed
			buff1, buff2 bool // we've read 'from', but not sent (yet)
			ok           bool // did we read successfully?
			from1, from2 Any  // what we've read
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

// Note: mergeAny is not my own.
// Just: I forgot where found the original merge2 - please accept my apologies.
// I'd love to learn about it's origin/author, so I can give credit.
// Thus: Your hint, dear reader, is highly appreciated!

// End of AnyMerge

// ===========================================================================
// Beg of AnySame comparator

// inspired by go/doc/play/tree.go

// AnySame reads values from two channels in lockstep
// and iff they have the same contents then
// `true` is sent on the returned bool channel
// before close.
func AnySame(same func(a, b Any) bool, inp1, inp2 <-chan Any) (out <-chan bool) {
	cha := make(chan bool)
	go sameAny(cha, same, inp1, inp2)
	return cha
}

func sameAny(out chan<- bool, same func(a, b Any) bool, inp1, inp2 <-chan Any) {
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

// End of AnySame comparator

// ===========================================================================
// Beg of AnyJoin feedback back-feeders for circular networks

// AnyJoin sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func AnyJoin(out chan<- Any, inp ...Any) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinAny(sig, out, inp...)
	return sig
}

func joinAny(done chan<- struct{}, out chan<- Any, inp ...Any) {
	defer close(done)
	for i := range inp {
		out <- inp[i]
	}
	done <- struct{}{}
}

// AnyJoinSlice sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func AnyJoinSlice(out chan<- Any, inp ...[]Any) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinAnySlice(sig, out, inp...)
	return sig
}

func joinAnySlice(done chan<- struct{}, out chan<- Any, inp ...[]Any) {
	defer close(done)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
	done <- struct{}{}
}

// AnyJoinChan sends inputs on the given out channel and returns a done channel to receive one signal when inp has been drained
func AnyJoinChan(out chan<- Any, inp <-chan Any) (done <-chan struct{}) {
	sig := make(chan struct{})
	go joinAnyChan(sig, out, inp)
	return sig
}

func joinAnyChan(done chan<- struct{}, out chan<- Any, inp <-chan Any) {
	defer close(done)
	for i := range inp {
		out <- i
	}
	done <- struct{}{}
}

// End of AnyJoin feedback back-feeders for circular networks

// ===========================================================================
// Beg of AnyDaisyChain

// AnyProc is the signature of the inner process of any linear pipe-network
//  Example: the identity proc:
// samesame := func(into chan<- Any, from <-chan Any) { into <- <-from }
// Note: type AnyProc is provided for documentation purpose only.
// The implementation uses the explicit function signature
// in order to avoid some genny-related issue.
//  Note: In https://talks.golang.org/2012/waza.slide#40
// Rob Pike uses a AnyProc named `worker`.
type AnyProc func(into chan<- Any, from <-chan Any)

// Example: the identity proc - see `samesame` below
var _ AnyProc = func(out chan<- Any, inp <-chan Any) {
	// `out <- <-inp` or `into <- <-from`
	defer close(out)
	for i := range inp {
		out <- i
	}
}

// daisyAny returns a channel to receive all inp after having passed thru process `proc`.
func daisyAny(
	inp <-chan Any, // a daisy to be chained
	proc func(into chan<- Any, from <-chan Any), // a process function
) (
	out chan Any, // to receive all results
) { //  Body:

	cha := make(chan Any)
	go proc(cha, inp)
	return cha
}

// AnyDaisyChain returns a channel to receive all inp
// after having passed
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making a null value useful.
func AnyDaisyChain(
	inp chan Any, // a daisy to be chained
	procs ...func(out chan<- Any, inp <-chan Any), // a process function
) (
	out chan Any, // to receive all results
) { //  Body:

	cha := inp

	if len(procs) < 1 {
		samesame := func(out chan<- Any, inp <-chan Any) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyAny(cha, samesame)
	} else {
		for _, proc := range procs {
			cha = daisyAny(cha, proc)
		}
	}
	return cha
}

// AnyDaisyChaiN returns a channel to receive all inp
// after having passed
// `somany` times
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If `somany` is less than 1 or no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making null values useful.
//
// Note: AnyDaisyChaiN(inp, 1, procs) <==> AnyDaisyChain(inp, procs)
func AnyDaisyChaiN(
	inp chan Any, // a daisy to be chained
	somany int, // how many times? so many times
	procs ...func(out chan<- Any, inp <-chan Any), // a process function
) (
	out chan Any, // to receive all results
) { //  Body:

	cha := inp

	if somany < 1 {
		samesame := func(out chan<- Any, inp <-chan Any) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyAny(cha, samesame)
	} else {
		for i := 0; i < somany; i++ {
			cha = AnyDaisyChain(cha, procs...)
		}
	}
	return cha
}

// End of AnyDaisyChain
// ===========================================================================
