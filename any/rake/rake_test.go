// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rake

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func ExampleRake_closure() {

	var r *Rake
	crawl := func(item Any) {
		if false {
			log.Println("have:", item)
		}
		for i := 0; i < rand.Intn(9)+2; i++ {
			r.Feed(rand.Intn(2000)) // up to 10 new numbers < 2.000
		}
		time.Sleep(time.Millisecond * 10)
	}

	r = New(crawl, nil, 80)

	r.Feed(1)

	<-r.Done()

	fmt.Println("Done")
	// Output:
	// Done
}

func ExampleRake_chained() {

	var r *Rake
	crawl := func(item Any) {
		if false {
			log.Println("have:", item)
		}
		for i := 0; i < rand.Intn(9)+2; i++ {
			r.Feed(rand.Intn(2000)) // up to 10 new numbers < 2.000
		}
		time.Sleep(time.Millisecond * 10)
	}

	r = New(nil, nil, 80).Rake(crawl).Feed(1)

	<-r.Done()

	fmt.Println("Done")
	// Output:
	// Done
}
