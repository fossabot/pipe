// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import (
	"github.com/cheekybits/genny/generic"
)

// anyThing is the generic type flowing thru the pipe network.
type anyThing generic.Type

// anyOwner is the generic who shall own the methods.
//  Note: Need to use `generic.Number` here as `generic.Type` is an interface and cannot have any method.
type anyOwner generic.Number

// ===========================================================================
// Beg of has nil versions

// Functions suitable only for types which can be == nil.
// Thus, do not use for basic built-in's such as int, string, ...

// anyThingChanFuncNil returns a channel to receive
// all results of generator `gen`
// until nil
// before close.
func (my anyOwner) anyThingChanFuncNil(gen func() anyThing) (out <-chan anyThing) {
	cha := make(chan anyThing)
	go my.chananyThingFuncNil(cha, gen)
	return cha
}

func (my anyOwner) chananyThingFuncNil(out chan<- anyThing, gen func() anyThing) {
	defer close(out)
	for {
		res := gen() // generate
		if res == nil {
			return
		}
		out <- res
	}
}

// End of has nil versions
// ===========================================================================
