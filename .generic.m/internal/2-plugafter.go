// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import "time"

// ===========================================================================
// Beg of ThingPlugAfter - graceful terminator

// ThingPlugAfter returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a time signal
// (e.g. from `time.After(...)`),
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func (inp ThingFrom) ThingPlugAfter(after <-chan time.Time) (out ThingFrom, done <-chan struct{}) {
	cha := make(chan Thing)
	doit := make(chan struct{})
	go inp.plugThingAfter(cha, doit, after)
	return cha, doit
}

func (inp ThingFrom) plugThingAfter(out ThingInto, done chan<- struct{}, after <-chan time.Time) {
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
// ===========================================================================