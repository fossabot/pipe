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
// Beg of anyThingFan2 easy fan-in's

// anyThingFan2 returns a channel to receive
// everything from `inp`
// as well as
// all inputs
// before close.
func anyThingFan2(inp <-chan anyThing, inps ...anyThing) (out <-chan anyThing) {
	return anyThingFanIn2(inp, anyThingChan(inps...))
}

// anyThingFan2Slice returns a channel to receive
// everything from `inp`
// as well as
// all inputs
// before close.
func anyThingFan2Slice(inp <-chan anyThing, inps ...[]anyThing) (out <-chan anyThing) {
	return anyThingFanIn2(inp, anyThingChanSlice(inps...))
}

// anyThingFan2Chan returns a channel to receive
// everything from `inp`
// as well as
// everything from `inp2`
// before close.
//  Note: anyThingFan2Chan is nothing but anyThingFanIn2
func anyThingFan2Chan(inp <-chan anyThing, inp2 <-chan anyThing) (out <-chan anyThing) {
	return anyThingFanIn2(inp, inp2)
}

// anyThingFan2FuncNok returns a channel to receive
// everything from `inp`
// as well as
// all results of generator `gen`
// until `!ok`
// before close.
func anyThingFan2FuncNok(inp <-chan anyThing, gen func() (anyThing, bool)) (out <-chan anyThing) {
	return anyThingFanIn2(inp, anyThingChanFuncNok(gen))
}

// anyThingFan2FuncErr returns a channel to receive
// everything from `inp`
// as well as
// all results of generator `gen`
// until `err != nil`
// before close.
func anyThingFan2FuncErr(inp <-chan anyThing, gen func() (anyThing, error)) (out <-chan anyThing) {
	return anyThingFanIn2(inp, anyThingChanFuncErr(gen))
}

// End of anyThingFan2 easy fan-in's
// ===========================================================================
