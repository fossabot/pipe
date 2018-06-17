// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import "time"

// ===========================================================================
// Beg of ScatterThing

// ScatterThing returns a slice (of size = size) of channels
// one of which shall receive any inp before close.
func ScatterThing(inp <-chan Thing, size int) (outS [](<-chan Thing)) {
	chaS := make([]chan Thing, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Thing)
	}

	go scatterThing(inp, chaS...)

	outS = make([]<-chan Thing, size)
	for i := 0; i < size; i++ {
		outS[i] = chaS[i] // convert `chan` to `<-chan`
	}

	return outS
}

// c scatterThing(inp <-chan Thing, outS ...chan<- Thing) {
// Note: go does not convert the passed slice `[]chan Thing` to `[]chan<- Thing` automatically.
// So, we do neither here, as we are lazy (we just call an internal helper function).
func scatterThing(inp <-chan Thing, outS ...chan Thing) {

	for i := range inp {
		for !trySendThing(i, outS...) {
			time.Sleep(time.Millisecond) // wait a little before retry
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

// End of FanThingOut
