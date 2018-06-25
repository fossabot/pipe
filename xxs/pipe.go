// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.

package pipe

import (
	"github.com/cheekybits/genny/generic"
)

// anyThing is the generic type flowing thru the pipe network.
type anyThing generic.Type

// ===========================================================================
// Beg of anySupply channel object

// anySupply is a
// supply channel
type anySupply struct {
	dat chan anyThing
	//  chan struct{}
}

// anySupplyMakeChan returns
// a (pointer to a) fresh
// unbuffered
// supply channel
func anySupplyMakeChan() *anySupply {
	d := anySupply{
		dat: make(chan anyThing),
		// : make(chan struct{}),
	}
	return &d
}

// anySupplyMakeBuff returns
// a (pointer to a) fresh
// buffered (with capacity=`cap`)
// supply channel
func anySupplyMakeBuff(cap int) *anySupply {
	d := anySupply{
		dat: make(chan anyThing, cap),
		// : make(chan struct{}),
	}
	return &d
}

// Provide is the send method
// - aka "myAnyChan <- myAny"
func (c *anySupply) Provide(dat anyThing) {
	// .req
	c.dat <- dat
}

// Receive is the receive operator as method
// - aka "myAny := <-myAnyChan"
func (c *anySupply) Receive() (dat anyThing) {
	// eq <- struct{}{}
	return <-c.dat
}

// Request is the comma-ok multi-valued form of Receive and
// reports whether a received value was sent before the anyThing channel was closed
func (c *anySupply) Request() (dat anyThing, open bool) {
	// eq <- struct{}{}
	dat, open = <-c.dat
	return dat, open
}

// Close closes the underlying anyThing channel
func (c *anySupply) Close() {
	close(c.dat)
}

// Cap reports the capacity of the underlying anyThing channel
func (c *anySupply) Cap() int {
	return cap(c.dat)
}

// Len reports the length of the underlying anyThing channel
func (c *anySupply) Len() int {
	return len(c.dat)
}

// End of anySupply channel object
// ===========================================================================

// ===========================================================================
// Beg of anyThingChan producers

// anyThingChan returns a channel to receive
// all inputs
// before close.
func anyThingChan(inp ...anyThing) (out *anySupply) {
	cha := anySupplyMakeChan()
	go chananyThing(cha, inp...)
	return cha
}

func chananyThing(out *anySupply, inp ...anyThing) {
	defer out.Close()
	for i := range inp {
		out.Provide(inp[i])
	}
}

// anyThingChanSlice returns a channel to receive
// all inputs
// before close.
func anyThingChanSlice(inp ...[]anyThing) (out *anySupply) {
	cha := anySupplyMakeChan()
	go chananyThingSlice(cha, inp...)
	return cha
}

func chananyThingSlice(out *anySupply, inp ...[]anyThing) {
	defer out.Close()
	for i := range inp {
		for j := range inp[i] {
			out.Provide(inp[i][j])
		}
	}
}

// anyThingChanFuncNok returns a channel to receive
// all results of generator `gen`
// until `!ok`
// before close.
func anyThingChanFuncNok(gen func() (anyThing, bool)) (out *anySupply) {
	cha := anySupplyMakeChan()
	go chananyThingFuncNok(cha, gen)
	return cha
}

func chananyThingFuncNok(out *anySupply, gen func() (anyThing, bool)) {
	defer out.Close()
	for {
		res, ok := gen() // generate
		if !ok {
			return
		}
		out.Provide(res)
	}
}

// anyThingChanFuncErr returns a channel to receive
// all results of generator `gen`
// until `err != nil`
// before close.
func anyThingChanFuncErr(gen func() (anyThing, error)) (out *anySupply) {
	cha := anySupplyMakeChan()
	go chananyThingFuncErr(cha, gen)
	return cha
}

func chananyThingFuncErr(out *anySupply, gen func() (anyThing, error)) {
	defer out.Close()
	for {
		res, err := gen() // generate
		if err != nil {
			return
		}
		out.Provide(res)
	}
}

// End of anyThingChan producers
// ===========================================================================

// ===========================================================================
// Beg of anyThingPipe functions

// anyThingPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be anyThingPipeMap for functional people,
// but 'map' has a very different meaning in go lang.
func anyThingPipeFunc(inp *anySupply, act func(a anyThing) anyThing) (out *anySupply) {
	cha := anySupplyMakeChan()
	if act == nil {
		act = func(a anyThing) anyThing { return a }
	}
	go pipeanyThingFunc(cha, inp, act)
	return cha
}

func pipeanyThingFunc(out *anySupply, inp *anySupply, act func(a anyThing) anyThing) {
	defer out.Close()
	for i, ok := inp.Request(); ok; i, ok = inp.Request() {
		out.Provide(act(i))
	}
}

// End of anyThingPipe functions
// ===========================================================================

// ===========================================================================
// Beg of anyThingTube closures

// anyThingTubeFunc returns a closure around PipeanyThingFunc (_, act).
func anyThingTubeFunc(act func(a anyThing) anyThing) (tube func(inp *anySupply) (out *anySupply)) {

	return func(inp *anySupply) (out *anySupply) {
		return anyThingPipeFunc(inp, act)
	}
}

// End of anyThingTube closures
// ===========================================================================

// ===========================================================================
// Beg of anyThingDone terminators

// anyThingDone returns a channel to receive
// one signal before close after `inp` has been drained.
func anyThingDone(inp *anySupply) (done <-chan struct{}) {
	sig := make(chan struct{})
	go doitanyThing(sig, inp)
	return sig
}

func doitanyThing(done chan<- struct{}, inp *anySupply) {
	defer close(done)
	for i, ok := inp.Request(); ok; i, ok = inp.Request() {
		_ = i // Drain inp
	}
	done <- struct{}{}
}

// anyThingDoneSlice returns a channel to receive
// a slice with every anyThing received on `inp`
// before close.
//
//  Note: Unlike anyThingDone, anyThingDoneSlice sends the fully accumulated slice, not just an event, once upon close of inp.
func anyThingDoneSlice(inp *anySupply) (done <-chan []anyThing) {
	sig := make(chan []anyThing)
	go doitanyThingSlice(sig, inp)
	return sig
}

func doitanyThingSlice(done chan<- []anyThing, inp *anySupply) {
	defer close(done)
	slice := []anyThing{}
	for i, ok := inp.Request(); ok; i, ok = inp.Request() {
		slice = append(slice, i)
	}
	done <- slice
}

// anyThingDoneFunc returns a channel to receive
// one signal after `act` has been applied to every `inp`
// before close.
func anyThingDoneFunc(inp *anySupply, act func(a anyThing)) (done <-chan struct{}) {
	sig := make(chan struct{})
	if act == nil {
		act = func(a anyThing) { return }
	}
	go doitanyThingFunc(sig, inp, act)
	return sig
}

func doitanyThingFunc(done chan<- struct{}, inp *anySupply, act func(a anyThing)) {
	defer close(done)
	for i, ok := inp.Request(); ok; i, ok = inp.Request() {
		act(i) // apply action
	}
	done <- struct{}{}
}

// End of anyThingDone terminators
// ===========================================================================

// ===========================================================================
// Beg of anyThingFini closures

// anyThingFini returns a closure around `DoneanyThing(_)`.
func anyThingFini() func(inp *anySupply) (done <-chan struct{}) {

	return func(inp *anySupply) (done <-chan struct{}) {
		return anyThingDone(inp)
	}
}

// anyThingFiniSlice returns a closure around `DoneanyThingSlice(_)`.
func anyThingFiniSlice() func(inp *anySupply) (done <-chan []anyThing) {

	return func(inp *anySupply) (done <-chan []anyThing) {
		return anyThingDoneSlice(inp)
	}
}

// anyThingFiniFunc returns a closure around `DoneanyThingFunc(_, act)`.
func anyThingFiniFunc(act func(a anyThing)) func(inp *anySupply) (done <-chan struct{}) {

	return func(inp *anySupply) (done <-chan struct{}) {
		return anyThingDoneFunc(inp, act)
	}
}

// End of anyThingFini closures
// ===========================================================================

// ===========================================================================
// Beg of anyThingPair functions

// anyThingPair returns a pair of channels to receive every result of inp before close.
//  Note: Yes, it is a VERY simple fanout - but sometimes all You need.
func anyThingPair(inp *anySupply) (out1, out2 *anySupply) {
	cha1 := anySupplyMakeChan()
	cha2 := anySupplyMakeChan()
	go pairanyThing(cha1, cha2, inp)
	return cha1, cha2
}

func pairanyThing(out1, out2 *anySupply, inp *anySupply) {
	defer out1.Close()
	defer out2.Close()
	for i, ok := inp.Request(); ok; i, ok = inp.Request() {
		out1.Provide(i)
		out2.Provide(i)
	}
}

// End of anyThingPair functions
// ===========================================================================
