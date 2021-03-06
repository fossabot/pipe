// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.

package httpsyet

// ===========================================================================
// Beg of stringMake creators

// stringMakeChan returns a new open channel
// (simply a 'chan string' that is).
// Note: No 'string-producer' is launched here yet! (as is in all the other functions).
//  This is useful to easily create corresponding variables such as:
/*
var mystringPipelineStartsHere := stringMakeChan()
// ... lot's of code to design and build Your favourite "mystringWorkflowPipeline"
   // ...
   // ... *before* You start pouring data into it, e.g. simply via:
   for drop := range water {
mystringPipelineStartsHere <- drop
   }
close(mystringPipelineStartsHere)
*/
//  Hint: especially helpful, if Your piping library operates on some hidden (non-exported) type
//  (or on a type imported from elsewhere - and You don't want/need or should(!) have to care.)
//
// Note: as always (except for stringPipeBuffer) the channel is unbuffered.
//
func stringMakeChan() (out chan string) {
	return make(chan string)
}

// End of stringMake creators
// ===========================================================================

// ===========================================================================
// Beg of stringChan producers

// stringChan returns a channel to receive
// all inputs
// before close.
func stringChan(inp ...string) (out <-chan string) {
	cha := make(chan string)
	go chanstring(cha, inp...)
	return cha
}

func chanstring(out chan<- string, inp ...string) {
	defer close(out)
	for i := range inp {
		out <- inp[i]
	}
}

// stringChanSlice returns a channel to receive
// all inputs
// before close.
func stringChanSlice(inp ...[]string) (out <-chan string) {
	cha := make(chan string)
	go chanstringSlice(cha, inp...)
	return cha
}

func chanstringSlice(out chan<- string, inp ...[]string) {
	defer close(out)
	for i := range inp {
		for j := range inp[i] {
			out <- inp[i][j]
		}
	}
}

// stringChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func stringChanFuncNok(gen func() (string, bool)) (out <-chan string) {
	cha := make(chan string)
	go chanstringFuncNok(cha, gen)
	return cha
}

func chanstringFuncNok(out chan<- string, gen func() (string, bool)) {
	defer close(out)
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out <- res
	}
}

// stringChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func stringChanFuncErr(gen func() (string, error)) (out <-chan string) {
	cha := make(chan string)
	go chanstringFuncErr(cha, gen)
	return cha
}

func chanstringFuncErr(out chan<- string, gen func() (string, error)) {
	defer close(out)
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out <- res
	}
}

// End of stringChan producers
// ===========================================================================

// ===========================================================================
// Beg of stringPipe functions

// stringPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be stringPipeMap for functional people,
// but 'map' has a very different meaning in go lang.
func stringPipeFunc(inp <-chan string, act func(a string) string) (out <-chan string) {
	cha := make(chan string)
	if act == nil { // Make `nil` value useful
		act = func(a string) string { return a }
	}
	go pipestringFunc(cha, inp, act)
	return cha
}

func pipestringFunc(out chan<- string, inp <-chan string, act func(a string) string) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of stringPipe functions
// ===========================================================================

// ===========================================================================
// Beg of stringTube closures around stringPipe

// stringTubeFunc returns a closure around PipeStringFunc (_, act).
func stringTubeFunc(act func(a string) string) (tube func(inp <-chan string) (out <-chan string)) {

	return func(inp <-chan string) (out <-chan string) {
		return stringPipeFunc(inp, act)
	}
}

// End of stringTube closures around stringPipe
// ===========================================================================

// ===========================================================================
// Beg of stringDone terminators

// stringDone returns a channel to receive
// one signal
// upon close
// and after `inp` has been drained.
func stringDone(inp <-chan string) (done <-chan struct{}) {
	sig := make(chan struct{})
	go donestring(sig, inp)
	return sig
}

func donestring(done chan<- struct{}, inp <-chan string) {
	defer close(done)
	for i := range inp {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// stringDoneSlice returns a channel to receive
// a slice with every string received on `inp`
// upon close.
//
// Note: Unlike stringDone, stringDoneSlice sends the fully accumulated slice, not just an event, once upon close of inp.
func stringDoneSlice(inp <-chan string) (done <-chan []string) {
	sig := make(chan []string)
	go donestringSlice(sig, inp)
	return sig
}

func donestringSlice(done chan<- []string, inp <-chan string) {
	defer close(done)
	slice := []string{}
	for i := range inp {
		slice = append(slice, i)
	}
	done <- slice
}

// stringDoneFunc
// will apply `act` to every `inp` and
// returns a channel to receive
// one signal
// upon close.
func stringDoneFunc(inp <-chan string, act func(a string)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a string) { return }
	}
	go donestringFunc(sig, inp, act)
	return sig
}

func donestringFunc(done chan<- struct{}, inp <-chan string, act func(a string)) {
	defer close(done)
	for i := range inp {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of stringDone terminators
// ===========================================================================

// ===========================================================================
// Beg of stringFini closures

// stringFini returns a closure around `stringDone(_)`.
func stringFini() func(inp <-chan string) (done <-chan struct{}) {

	return func(inp <-chan string) (done <-chan struct{}) {
		return stringDone(inp)
	}
}

// stringFiniSlice returns a closure around `stringDoneSlice(_)`.
func stringFiniSlice() func(inp <-chan string) (done <-chan []string) {

	return func(inp <-chan string) (done <-chan []string) {
		return stringDoneSlice(inp)
	}
}

// stringFiniFunc returns a closure around `stringDoneFunc(_, act)`.
func stringFiniFunc(act func(a string)) func(inp <-chan string) (done <-chan struct{}) {

	return func(inp <-chan string) (done <-chan struct{}) {
		return stringDoneFunc(inp, act)
	}
}

// End of stringFini closures
// ===========================================================================

// ===========================================================================
// Beg of stringPair functions

// stringPair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func stringPair(inp <-chan string) (out1, out2 <-chan string) {
	cha1 := make(chan string)
	cha2 := make(chan string)
	go pairstring(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func pairstring(out1, out2 chan<- string, inp <-chan string) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func pairstring(out1, out2 chan<- string, inp <-chan string) {
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

// End of stringPair functions
// ===========================================================================

// ===========================================================================
// Beg of stringFork functions

// stringFork returns two channels
// either of which is to receive
// every result of inp
// before close.
func stringFork(inp <-chan string) (out1, out2 <-chan string) {
	cha1 := make(chan string)
	cha2 := make(chan string)
	go forkstring(cha1, cha2, inp)
	return cha1, cha2
}

/* not used - kept for reference only.
func forkstring(out1, out2 chan<- string, inp <-chan string) {
	defer close(out1)
	defer close(out2)
	for i := range inp {
		out1 <- i
		out2 <- i
	}
} */

func forkstring(out1, out2 chan<- string, inp <-chan string) {
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

// End of stringFork functions
// ===========================================================================

// ===========================================================================
// Beg of stringFanIn2 simple binary Fan-In

// stringFanIn2 returns a channel to receive
// all from both `inp1` and `inp2`
// before close.
func stringFanIn2(inp1, inp2 <-chan string) (out <-chan string) {
	cha := make(chan string)
	go fanIn2string(cha, inp1, inp2)
	return cha
}

/* not used - kept for reference only.
// fanin2string as seen in Go Concurrency Patterns
func fanin2string(out chan<- string, inp1, inp2 <-chan string) {
	for {
		select {
		case e := <-inp1:
			out <- e
		case e := <-inp2:
			out <- e
		}
	}
} */

func fanIn2string(out chan<- string, inp1, inp2 <-chan string) {
	defer close(out)

	var (
		closed bool   // we found a chan closed
		ok     bool   // did we read successfully?
		e      string // what we've read
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

// End of stringFanIn2 simple binary Fan-In
// ===========================================================================
