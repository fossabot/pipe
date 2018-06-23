// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

// ===========================================================================
// Beg of anyThingPipe functions

// anyThingPipeFunc returns a channel to receive
// every result of action `act` applied to `inp`
// before close.
// Note: it 'could' be PipeanyThingMap for functional people,
// but 'map' has a very different meaning in go lang.
func (my anyOwner) anyThingPipeFunc(inp <-chan anyThing, act func(a anyThing) anyThing) (out <-chan anyThing) {
	cha := make(chan anyThing)
	if act == nil { // Make `nil` value useful
		act = func(a anyThing) anyThing { return a }
	}
	go my.pipeanyThingFunc(cha, inp, act)
	return cha
}

func (my anyOwner) pipeanyThingFunc(out chan<- anyThing, inp <-chan anyThing, act func(a anyThing) anyThing) {
	defer close(out)
	for i := range inp {
		out <- act(i) // apply action
	}
}

// End of anyThingPipe functions
// ===========================================================================
