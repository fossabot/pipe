// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipe

// ===========================================================================
// Beg of anyThingTube closures

// anyThingTubeFunc returns a closure around PipeanyThingFunc (_, act).
func anyThingTubeFunc(act func(a anyThing) anyThing) (tube func(inp anymode) (out anymode)) {

	return func(inp anymode) (out anymode) {
		return anyThingPipeFunc(inp, act)
	}
}

// anyThingTubeBuffer returns a closure around PipeanyThingBuffer (_, cap).
func anyThingTubeBuffer(cap int) (tube func(inp anymode) (out anymode)) {

	return func(inp anymode) (out anymode) {
		return anyThingPipeBuffer(inp, cap)
	}
}

// End of anyThingTube closures
// ===========================================================================
