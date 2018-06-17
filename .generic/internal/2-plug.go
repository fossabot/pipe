// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

// ===========================================================================
// Beg of PlugThing - graceful terminator

// PlugThing returns a channel to receive every `inp` before close and a channel to signal this closing.
// Upon receipt of a stop signal,
// output is immediately closed,
// and for graceful termination
// any remaining input is drained before done is signalled.
func PlugThing(inp <-chan Thing, stop <-chan struct{}) (out <-chan Thing, done <-chan struct{}) {
	cha := make(chan Thing)
	doit := make(chan struct{})
	go plugThing(cha, doit, inp, stop)
	return cha, doit
}

func plugThing(out chan<- Thing, done chan<- struct{}, inp <-chan Thing, stop <-chan struct{}) {
	defer close(done)

	var ok bool // did we read sucessfully?
	var e Thing // what we've read
	for {
		select {
		case e, ok = <-inp:
			if ok {
				out <- e
			} else {
				break
			}
		case <-stop:
			break
		}
	}

	close(out)

	for _ = range inp {
		// drain inp
	}

	done <- struct{}{}
}

// End of PlugThing - graceful terminator
