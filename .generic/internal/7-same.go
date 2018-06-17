// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

// ===========================================================================
// Beg of SameThing comparator

// inspired by go/doc/play/tree.go

// SameThing reads values from two channels in lockstep
// and iff they have the same contents then
// `true` is sent on the returned bool channel
// before close.
func SameThing(same func(a, b Thing) bool, inp1, inp2 <-chan Thing) (out <-chan bool) {
	cha := make(chan bool)
	go sameThing(cha, same, inp1, inp2)
	return cha
}

func sameThing(out chan<- bool, same func(a, b Thing) bool, inp1, inp2 <-chan Thing) {
	defer close(out)
	for {
		v1, ok1 := <-inp1
		v2, ok2 := <-inp2

		if !ok1 || !ok2 {
			out <- ok1 == ok2
			return
		}
		if !same(v1, v2) {
			out <- false
			return
		}
	}
}

// End of SameThing comparator
