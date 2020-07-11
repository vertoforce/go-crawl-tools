package crawltools

import (
	"context"
	"net/http"
	"sync"

	"github.com/vertoforce/go-crawl-tools/proxy"
)

// PageURLFunc is a function that returns the appropriate URL based on the page number we are requesting
type PageURLFunc func(page int64) string

// CrawlPages Crawls all pages using the supplied parameters.
// It starts at page 1, getting the total number of pages, then spawns multiple threads to call all the other pages
func CrawlPages(ctx context.Context, pageURLFunc PageURLFunc, parseFunction ParseFunc, totalPagesFunc TotalPagesFunc, p proxy.Proxy, maxThreads int) chan interface{} {
	ret := make(chan interface{})

	threadLimit := make(chan struct{}, maxThreads)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(ret)
		defer close(threadLimit)
		// Crawl first page
		req, err := http.NewRequest("GET", pageURLFunc(1), nil)
		if err != nil {
			return
		}
		html, err := CrawlPage(ctx, p, req, ret, parseFunction)
		if err != nil {
			// TODO: Return this somehow
			return
		}
		pageCount := totalPagesFunc(ctx, html)

		// Crawl all other pages
		for page := int64(2); page < pageCount; page++ {
			select {
			case <-ctx.Done():
				return
			case threadLimit <- struct{}{}: // Try to consume a new thread
			}

			wg.Add(1)
			go func(currentPage int64) {
				// Build request
				req, err := http.NewRequest("GET", pageURLFunc(currentPage), nil)
				if err != nil {
					return
				}

				_, _ = CrawlPage(ctx, p, req, ret, parseFunction)
				// Mark thread as done
				<-threadLimit
				wg.Done()
			}(page)
		}

		// wait for all threads to be done
		wg.Wait()
	}()

	return ret
}
