// Copyright 2017 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpsyet

import (
	"fmt"
	"net/url"
	"time"

	"github.com/GoLangsam/pipe/examples/httpsyet/v4/sites"
)

// ---------------------------------------------------------------------------

// alias types in order to leave source code unaffected
type site = sites.Site

// ---------------------------------------------------------------------------
// make `result` an explicit type, and
// teach `Crawler` how to `report` a result.

// result represents (a secondary) observation
type result string

// report prints a result to Crawler.Out;
// used by `DoneResultFunc`.
func (c *Crawler) report(r result) {
	if _, err := fmt.Fprintln(c.Out, r); err != nil {
		c.Log.Printf("failed to write output '%s': %v\n", r, err)
	}
}

// ---------------------------------------------------------------------------

// crawling represents a crawling Crawler ... busy crawling ...
type crawling struct {
	*Crawler                   // config
	*sites.Traffic             // the cirular network
	results        chan result // to be reported
}

// ---------------------------------------------------------------------------
// teach `*crawling` how to crawl :-)

// crawl performs a crawling Crawler's main function: crawl.
func (c *crawling) crawl(s site) {
	urls := c.crawlSite(s)         // the core crawl process
	c.Feed(urls, s.URL, s.Depth-1) // new urls enter crawling - circular feedback
	time.Sleep(c.Delay)            // have a gentle nap
}

// ===========================================================================

// a Crawler crawling: crawl traffic, emit results
// and signal when done
func (c *Crawler) crawling(urls []*url.URL) (done <-chan struct{}) {
	crawling := crawling{
		c,                 // "Crawler is used as configuration ..."
		sites.New(),       // the cirular network
		make(chan result), // results - the (secondary) output
	}

	parallel := parallel(c.Parallel) // no idea what keeps Crawler from setting `Parallel` upon validation

	crawlingDone := crawling.Processor(crawling.crawl, parallel) // build the site traffic processing network

	go func() { // launch the results closer
		<-crawlingDone // block 'till crawling.Processor is done
		close(crawling.results)
	}()

	crawling.Feed(urls, nil, c.Depth) // feed initial urls

	return (resultFrom)(crawling.results).resultDoneFunc(c.report) // signal when results report are done
}

// ===========================================================================
