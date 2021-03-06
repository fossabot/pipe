// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// adjustments and embedding:
// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// original found in $GOROOT/test/chan/sieve2.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import "container/ring"

// Note: pipeThingAdjust imports "container/ring" for the expanding buffer.

// ===========================================================================
// Beg of ThingPipeAdjust

// ThingPipeAdjust returns a channel to receive
// all `inp`
// buffered by a ThingSendProxy process
// before close.
func (inp ThingFrom) ThingPipeAdjust(sizes ...int) (out ThingFrom) {
	cap, que := sendThingProxySizes(sizes...)
	cha := make(chan Thing, cap)
	go inp.pipeThingAdjust(cha, que)
	return cha
}

// ThingTubeAdjust returns a closure around ThingPipeAdjust (_, sizes ...int).
func (inp ThingFrom) ThingTubeAdjust(sizes ...int) (tube func(inp ThingFrom) (out ThingFrom)) {

	return func(inp ThingFrom) (out ThingFrom) {
		return inp.ThingPipeAdjust(sizes...)
	}
}

// End of ThingPipeAdjust
// ===========================================================================

// ===========================================================================
// Beg of sendThingProxy

func sendThingProxySizes(sizes ...int) (cap, que int) {

	// CAP is the minimum capacity of the buffered proxy channel in `ThingSendProxy`
	const CAP = 10

	// QUE is the minimum initially allocated size of the circular queue in `ThingSendProxy`
	const QUE = 16

	cap = CAP
	que = QUE

	if len(sizes) > 0 && sizes[0] > CAP {
		que = sizes[0]
	}

	if len(sizes) > 1 && sizes[1] > QUE {
		que = sizes[1]
	}

	if len(sizes) > 2 {
		panic("ThingSendProxy: too many sizes")
	}

	return
}

// ThingSendProxy returns a channel to serve as a sending proxy to 'out'.
// Uses a goroutine to receive values from 'out' and store them
// in an expanding buffer, so that sending to 'out' never blocks.
//  Note: the expanding buffer is implemented via "container/ring"
//
// Note: ThingSendProxy is kept for the Sieve example
// and other dynamic use to be discovered
// even so it does not fit the pipe tube pattern as ThingPipeAdjust does.
func ThingSendProxy(out ThingInto, sizes ...int) (send ThingInto) {
	cap, que := sendThingProxySizes(sizes...)
	cha := make(chan Thing, cap)
	go (ThingFrom)(cha).pipeThingAdjust(out, que)
	return cha
}

// pipeThingAdjust uses an adjusting buffer to receive from 'inp'
// even so 'out' is not ready to receive yet. The buffer may grow
// until 'inp' is closed and then will shrink by every send to 'out'.
//  Note: the adjusting buffer is implemented via "container/ring"
func (inp ThingFrom) pipeThingAdjust(out ThingInto, QUE int) {
	defer close(out)
	n := QUE // the allocated size of the circular queue
	first := ring.New(n)
	last := first
	var c ThingInto
	var e Thing
	ok := true
	for ok {
		c = out
		if first == last {
			c = nil // buffer empty: disable output
		} else {
			e = first.Value.(Thing)
		}
		select {
		case e, ok = <-inp:
			if ok {
				last.Value = e
				if last.Next() == first {
					last.Link(ring.New(n)) // buffer full: expand it
					n *= 2
				}
				last = last.Next()
			}
		case c <- e:
			first = first.Next()
		}
	}

	for first != last {
		out <- first.Value.(Thing)
		first = first.Unlink(1) // first.Next()
	}
}

// End of sendThingProxy
// ===========================================================================
