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

package sites

import "container/ring"

// Note: pipeSiteAdjust imports "container/ring" for the expanding buffer.

// ===========================================================================
// Beg of SitePipeAdjust

// SitePipeAdjust returns a channel to receive
// all `inp`
// buffered by a SiteSendProxy process
// before close.
func (my *Traffic) SitePipeAdjust(inp <-chan Site, sizes ...int) (out <-chan Site) {
	cap, que := my.sendSiteProxySizes(sizes...)
	cha := make(chan Site, cap)
	go my.pipeSiteAdjust(cha, inp, que)
	return cha
}

// SiteTubeAdjust returns a closure around SitePipeAdjust (_, sizes ...int).
func (my *Traffic) SiteTubeAdjust(sizes ...int) (tube func(inp <-chan Site) (out <-chan Site)) {

	return func(inp <-chan Site) (out <-chan Site) {
		return my.SitePipeAdjust(inp, sizes...)
	}
}

// End of SitePipeAdjust
// ===========================================================================

// ===========================================================================
// Beg of sendSiteProxy

func (my *Traffic) sendSiteProxySizes(sizes ...int) (cap, que int) {

	// CAP is the minimum capacity of the buffered proxy channel in `SiteSendProxy`
	const CAP = 10

	// QUE is the minimum initially allocated size of the circular queue in `SiteSendProxy`
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
		panic("SiteSendProxy: too many sizes")
	}

	return
}

// SiteSendProxy returns a channel to serve as a sending proxy to 'out'.
// Uses a goroutine to receive values from 'out' and store them
// in an expanding buffer, so that sending to 'out' never blocks.
//  Note: the expanding buffer is implemented via "container/ring"
//
// Note: SiteSendProxy is kept for the Sieve example
// and other dynamic use to be discovered
// even so it does not fit the pipe tube pattern as SitePipeAdjust does.
func (my *Traffic) SiteSendProxy(out chan<- Site, sizes ...int) chan<- Site {
	cap, que := my.sendSiteProxySizes(sizes...)
	cha := make(chan Site, cap)
	go my.pipeSiteAdjust(out, cha, que)
	return cha
}

// pipeSiteAdjust uses an adjusting buffer to receive from 'inp'
// even so 'out' is not ready to receive yet. The buffer may grow
// until 'inp' is closed and then will shrink by every send to 'out'.
//  Note: the adjusting buffer is implemented via "container/ring"
func (my *Traffic) pipeSiteAdjust(out chan<- Site, inp <-chan Site, QUE int) {
	defer close(out)
	n := QUE // the allocated size of the circular queue
	first := ring.New(n)
	last := first
	var c chan<- Site
	var e Site
	ok := true
	for ok {
		c = out
		if first == last {
			c = nil // buffer empty: disable output
		} else {
			e = first.Value.(Site)
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
		out <- first.Value.(Site)
		first = first.Unlink(1) // first.Next()
	}
}

// End of sendSiteProxy