// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.

package result

// ===========================================================================
// Beg of ResultMake creators

// ResultMakeChan returns a new open channel
// (simply a 'chan Result' that is).
// Note: No 'Result-producer' is launched here yet! (as is in all the other functions).
//  This is useful to easily create corresponding variables such as:
/*
var myResultPipelineStartsHere := ResultMakeChan()
// ... lot's of code to design and build Your favourite "myResultWorkflowPipeline"
   // ...
   // ... *before* You start pouring data into it, e.g. simply via:
   for drop := range water {
myResultPipelineStartsHere <- drop
   }
close(myResultPipelineStartsHere)
*/
//  Hint: especially helpful, if Your piping library operates on some hidden (non-exported) type
//  (or on a type imported from elsewhere - and You don't want/need or should(!) have to care.)
//
// Note: as always (except for ResultPipeBuffer) the channel is unbuffered.
//
func ResultMakeChan() (out chan Result) {
	return make(chan Result)
}

// End of ResultMake creators
// ===========================================================================

// ===========================================================================
// Beg of ResultChan producers

// ResultChan returns a channel to receive
// all inputs
// before close.
func ResultChan(inp ...Result) (out <-chan Result) {
	cha := make(chan Result)
	go chanResult(cha, inp...)
	return cha
}

func chanResult(out chan<- Result, inp ...Result) {
	defer close(out)
	for i := range inp {
		out <- inp[i]
	}
}

// ResultChanSlice returns a channel to receive
// all inputs
// before close.
func ResultChanSlice(inp ...[]Result) (out <-chan Result) {
	cha := make(chan Result)
	go chanResultSlice(cha, inp...)
	return cha
}

func chanResultSlice(out chan<- Result, inp ...[]Result) {
	defer close(out)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
}

// ResultChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func ResultChanFuncNok(gen func() (Result, bool)) (out <-chan Result) {
	cha := make(chan Result)
	go chanResultFuncNok(cha, gen)
	return cha
}

func chanResultFuncNok(out chan<- Result, gen func() (Result, bool)) {
	defer close(out)
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out <- res
	}
}

// ResultChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func ResultChanFuncErr(gen func() (Result, error)) (out <-chan Result) {
	cha := make(chan Result)
	go chanResultFuncErr(cha, gen)
	return cha
}

func chanResultFuncErr(out chan<- Result, gen func() (Result, error)) {
	defer close(out)
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out <- res
	}
}

// End of ResultChan producers
// ===========================================================================

// ===========================================================================
// Beg of ResultPipe functions

// ResultPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be ResultPipeMap for functional people,
// but 'map' has a very different meaning in go lang.
func ResultPipeFunc(inp <-chan Result, act func(a Result) Result) (out <-chan Result) {
	cha := make(chan Result)
	if act == nil { // Make `nil` value useful
		act = func(a Result) Result { return a }
	}
	go pipeResultFunc(cha, inp, act)
	return cha
}

func pipeResultFunc(out chan<- Result, inp <-chan Result, act func(a Result) Result) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of ResultPipe functions
// ===========================================================================

// ===========================================================================
// Beg of ResultTube closures around ResultPipe

// ResultTubeFunc returns a closure around PipeResultFunc (_, act).
func ResultTubeFunc(act func(a Result) Result) (tube func(inp <-chan Result) (out <-chan Result)) {

	return func(inp <-chan Result) (out <-chan Result) {
		return ResultPipeFunc(inp, act)
	}
}

// End of ResultTube closures around ResultPipe
// ===========================================================================

// ===========================================================================
// Beg of ResultDone terminators

// ResultDone returns a channel to receive
// one signal
// upon close
// and after `inp` has been drained.
func ResultDone(inp <-chan Result) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doneResult(sig, inp)
	return sig
}

func doneResult(done chan<- struct{}, inp <-chan Result) {
	defer close(done)
	for i := range inp {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// ResultDoneSlice returns a channel to receive
// a slice with every Result received on `inp`
// upon close.
//
// Note: Unlike ResultDone, ResultDoneSlice sends the fully accumulated slice, not just an event, once upon close of inp.
func ResultDoneSlice(inp <-chan Result) (done <-chan []Result) {
	sig := make(chan []Result)
	go doneResultSlice(sig, inp)
	return sig
}

func doneResultSlice(done chan<- []Result, inp <-chan Result) {
	defer close(done)
	slice := []Result{}
	for i := range inp {
		slice = append(slice, i)
	}
	done <- slice
}

// ResultDoneFunc
// will apply `act` to every `inp` and
// returns a channel to receive
// one signal
// upon close.
func ResultDoneFunc(inp <-chan Result, act func(a Result)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a Result) { return }
	}
	go doneResultFunc(sig, inp, act)
	return sig
}

func doneResultFunc(done chan<- struct{}, inp <-chan Result, act func(a Result)) {
	defer close(done)
	for i := range inp {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of ResultDone terminators
// ===========================================================================

// ===========================================================================
// Beg of ResultFini closures

// ResultFini returns a closure around `ResultDone(_)`.
func ResultFini() func(inp <-chan Result) (done <-chan struct{}) {

	return func(inp <-chan Result) (done <-chan struct{}) {
		return ResultDone(inp)
	}
}

// ResultFiniSlice returns a closure around `ResultDoneSlice(_)`.
func ResultFiniSlice() func(inp <-chan Result) (done <-chan []Result) {

	return func(inp <-chan Result) (done <-chan []Result) {
		return ResultDoneSlice(inp)
	}
}

// ResultFiniFunc returns a closure around `ResultDoneFunc(_, act)`.
func ResultFiniFunc(act func(a Result)) func(inp <-chan Result) (done <-chan struct{}) {

	return func(inp <-chan Result) (done <-chan struct{}) {
		return ResultDoneFunc(inp, act)
	}
}

// End of ResultFini closures
// ===========================================================================

// ===========================================================================
// Beg of ResultPair functions

// ResultPair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func ResultPair(inp <-chan Result) (out1, out2 <-chan Result) {
	cha1 := make(chan Result)
	cha2 := make(chan Result)
	go pairResult(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func pairResult(out1, out2 chan<- Result, inp <-chan Result) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func pairResult(out1, out2 chan<- Result, inp <-chan Result) {
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

// End of ResultPair functions
// ===========================================================================

// ===========================================================================
// Beg of ResultFork functions

// ResultFork returns two channels
// either of which is to receive
// every result of inp
// before close.
func ResultFork(inp <-chan Result) (out1, out2 <-chan Result) {
	cha1 := make(chan Result)
	cha2 := make(chan Result)
	go forkResult(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func forkResult(out1, out2 chan<- Result, inp <-chan Result) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func forkResult(out1, out2 chan<- Result, inp <-chan Result) {
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

// End of ResultFork functions
// ===========================================================================

// ===========================================================================
// Beg of ResultFanIn2 simple binary Fan-In

// ResultFanIn2 returns a channel to receive
// all from both `inp1` and `inp2`
// before close.
func ResultFanIn2(inp1, inp2 <-chan Result) (out <-chan Result) {
	cha := make(chan Result)
	go fanIn2Result(cha, inp1, inp2)
	return cha
}

/* not used - kept for reference only.
// fanin2Result as seen in Go Concurrency Patterns
func fanin2Result(out chan<- Result, inp1, inp2 <-chan Result) {
	for {
		select {
		case e := <-inp1:
			out <- e
		case e := <-inp2:
			out <- e
		}
	}
} */

func fanIn2Result(out chan<- Result, inp1, inp2 <-chan Result) {
	defer close(out)

	var (
		closed bool   // we found a chan closed
		ok     bool   // did we read successfully?
		e      Result // what we've read
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

// End of ResultFanIn2 simple binary Fan-In
// ===========================================================================
