// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

import (
	"sync"

	"github.com/cheekybits/genny/generic"
)

// anyThing is the generic type flowing thru the pipe network.
type anyThing generic.Type

// ===========================================================================
// Beg of anyThingFanIn

// anyThingFanIn returns a channel to receive all inputs arriving
// on variadic inps
// before close.
//
//  Note: For each input one go routine is spawned to forward arrivals.
//
//  Ref: https://blog.golang.org/pipelines
//  Ref: https://github.com/QuentinPerez/go-stuff/channel/Fan-out-Fan-in/main.go
func anyThingFanIn(inps ...<-chan anyThing) (out <-chan anyThing) {
	cha := make(chan anyThing)

	wg := new(sync.WaitGroup)
	wg.Add(len(inps))

	go func(wg *sync.WaitGroup, out chan anyThing) { // Spawn "close(out)" once all inps are done
		wg.Wait()
		close(out)
	}(wg, cha)

	for i := range inps {
		go func(out chan<- anyThing, inp <-chan anyThing) { // Spawn "output(c)"s
			defer wg.Done()
			for i := range inp {
				out <- i
			}
		}(cha, inps[i])
	}

	return cha
}

// End of anyThingFanIn
// ===========================================================================
