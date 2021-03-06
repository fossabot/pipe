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
// Beg of anyThingDaisyChain

// anyThingProc is the signature of the inner process of any linear pipe-network
//  Example: the identity proc:
// samesame := func(into anyThingInto, from anyThingFrom) { into <- <-from }
//  Note: type anyThingProc is provided for documentation purpose only.
// The implementation uses the explicit function signature
// in order to avoid some genny-related issue.
//  Note: In https://talks.golang.org/2012/waza.slide#40
//  Rob Pike uses a anyThingProc named `worker`.
type anyThingProc func(into anyThingInto, from anyThingFrom)

// Example: the identity proc - see `samesame` below
var _ anyThingProc = func(out anyThingInto, inp anyThingFrom) {
	// `out <- <-inp` or `into <- <-from`
	defer close(out)
	for i := range inp {
		out <- i
	}
}

// daisyanyThing returns a channel to receive all inp after having passed thru process `proc`.
func daisyanyThing(
	inp anyThingFrom, // a daisy to be chained
	proc func(into anyThingInto, from anyThingFrom), // a process function
) (
	out chan anyThing, // to receive all results
) { //  Body:

	cha := make(chan anyThing)
	go proc(cha, inp)
	return cha
}

// anyThingDaisyChain returns a channel to receive all inp
// after having passed
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making a null value useful.
func anyThingDaisyChain(
	inp chan anyThing, // a daisy to be chained
	procs ...func(out anyThingInto, inp anyThingFrom), // a process function
) (
	out chan anyThing, // to receive all results
) { //  Body:

	cha := inp

	if len(procs) < 1 {
		samesame := func(out anyThingInto, inp anyThingFrom) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyanyThing(cha, samesame)
	} else {
		for _, proc := range procs {
			cha = daisyanyThing(cha, proc)
		}
	}
	return cha
}

// anyThingDaisyChaiN returns a channel to receive all inp
// after having passed
// `somany` times
// thru the process(es) (`from` right `into` left)
// before close.
//
// Note: If `somany` is less than 1 or no `tubes` are provided,
// `out` shall receive elements from `inp` unaltered (as a convenience),
// thus making null values useful.
//
//  Note: anyThingDaisyChaiN(inp, 1, procs) <==> anyThingDaisyChain(inp, procs)
func anyThingDaisyChaiN(
	inp chan anyThing, // a daisy to be chained
	somany int, // how many times? so many times
	procs ...func(out anyThingInto, inp anyThingFrom), // a process function
) (
	out chan anyThing, // to receive all results
) { //  Body:

	cha := inp

	if somany < 1 {
		samesame := func(out anyThingInto, inp anyThingFrom) {
			// `out <- <-inp` or `into <- <-from`
			defer close(out)
			for i := range inp {
				out <- i
			}
		}
		cha = daisyanyThing(cha, samesame)
	} else {
		for i := 0; i < somany; i++ {
			cha = anyThingDaisyChain(cha, procs...)
		}
	}
	return cha
}

// End of anyThingDaisyChain
// ===========================================================================
