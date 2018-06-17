// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

// ===========================================================================
// Beg of FanThingOut

// FanThingOut returns a slice (of size = size) of channels
// each of which shall receive any inp before close.
func FanThingOut(inp <-chan Thing, size int) (outS [](<-chan Thing)) {
	chaS := make([]chan Thing, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan Thing)
	}

	go fanThingOut(inp, chaS...)

	outS = make([]<-chan Thing, size)
	for i := 0; i < size; i++ {
		outS[i] = chaS[i] // convert `chan` to `<-chan`
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

// End of FanThingOut
