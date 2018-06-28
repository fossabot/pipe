// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import "time"

// ===========================================================================
// Beg of intStrew - scatter them

// intStrew returns a slice (of size = size) of channels
// one of which shall receive each inp before close.
func intStrew(inp <-chan int, size int) (outS [](<-chan int)) {
	chaS := make(map[chan int]struct{}, size)
	for i := 0; i < size; i++ {
		chaS[make(chan int)] = struct{}{}
	}

	go strewint(inp, chaS)

	outS = make([]<-chan int, size)
	i := 0
	for c := range chaS {
		outS[i] = (<-chan int)(c) // convert `chan` to `<-chan`
		i++
	}

	return outS
}

// c strewint(inp <-chan int, outS ...chan<- int) {
// Note: go does not convert the passed slice `[]chan int` to `[]chan<- int` automatically.
// So, we do neither here, as we are lazy (we just call an internal helper function).
func strewint(inp <-chan int, outS map[chan int]struct{}) {

	for i := range inp {
		for !trySendint(i, outS) {
			time.Sleep(time.Millisecond * 10) // wait a little before retry
		} // !sent
	} // inp

	for o := range outS {
		close(o)
	}
}

func trySendint(inp int, outS map[chan int]struct{}) bool {

	for o := range outS {

		select { // try to send
		case o <- inp:
			return true
		default:
			// keep trying
		}

	} // outS
	return false
}

// End of intStrew - scatter them
// ===========================================================================
