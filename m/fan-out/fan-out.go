// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import (
	"github.com/cheekybits/genny/generic"
)

// anyThing is the generic type flowing thru the pipe network.
type anyThing generic.Type

// ===========================================================================
// Beg of FanOut

// FanOut returns a slice (of size = size) of channels
// each of which shall receive any inp before close.
func (inp anyThingFrom) FanOut(size int) (outS [](anyThingFrom)) {
	chaS := make([]chan anyThing, size)
	for i := 0; i < size; i++ {
		chaS[i] = make(chan anyThing)
	}

	go inp.fanOut(chaS...)

	outS = make([]anyThingFrom, size)
	for i := 0; i < size; i++ {
		outS[i] = (anyThingFrom)(chaS[i]) // convert `chan` to `<-chan`
	}

	return outS
}

// c (inp anyThingFrom) fanOut(outs ...anyThingInto) {
func (inp anyThingFrom) fanOut(outs ...chan anyThing) {

	for i := range inp {
		for o := range outs {
			outs[o] <- i
		}
	}

	for o := range outs {
		close(outs[o])
	}

}

// End of FanOut
// ===========================================================================
