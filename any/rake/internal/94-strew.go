// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rake

import "time"

// ===========================================================================
// Beg of itemStrew - scatter them

// itemStrew returns a slice (of size = size) of channels
// one of which shall receive each inp before close.
func (inp itemFrom) itemStrew(size int) (outS []itemFrom) {
	chaS := make(map[chan item]struct{}, size)
	for i := 0; i < size; i++ {
		chaS[make(chan item)] = struct{}{}
	}

	go inp.strewitem(chaS)

	outS = make([]itemFrom, size)
	i := 0
	for c := range chaS {
		outS[i] = (itemFrom)(c) // convert `chan item` to itemFrom
		i++
	}

	return outS
}

func (inp itemFrom) strewitem(outS map[chan item]struct{}) {

	for i := range inp {
		for !inp.trySenditem(i, outS) {
			time.Sleep(time.Millisecond * 10) // wait a little before retry
		} // !sent
	} // inp

	for o := range outS {
		close(o)
	}
}

func (static itemFrom) trySenditem(inp item, outS map[chan item]struct{}) bool {

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

// End of itemStrew - scatter them
// ===========================================================================
