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
// It starts at page 1, getting the total number of pages, then spawns multiple threads to crawl all the other pages.
//
// Each item it finds it sends to itemsChannel.  If you know the type the ParseFunc returns, you can generatea a new channel type casting each item to that type.
// It will close the channel once it's finished
func CrawlPages(ctx context.Context, pageURLFunc PageURLFunc, parseFunction ParseFunc, totalPagesFunc TotalPagesFunc, p proxy.Proxy, maxThreads int, itemsChannel chan interface{}) {
	threadLimit := make(chan struct{}, maxThreads)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(itemsChannel)
		defer close(threadLimit)
		// Crawl first page
		req, err := http.NewRequestWithContext(ctx, "GET", pageURLFunc(1), nil)
		if err != nil {
			return
		}
		html, err := CrawlPage(ctx, p, req, itemsChannel, parseFunction)
		if err != nil {
			// TODO: Return this somehow
			return
		}
		pageCount := totalPagesFunc(ctx, html)

		// Crawl all other pages
		for page := int64(2); page < pageCount; page++ {
			select {
			case <-ctx.Done():
				// Wait for all our child threads to finish before returning
				// and closing the channel they rely on
				wg.Wait()
				return
			case threadLimit <- struct{}{}: // Try to consume a new thread
			}

			wg.Add(1)
			go func(currentPage int64) {
				// Build request
				req, err := http.NewRequestWithContext(ctx, "GET", pageURLFunc(currentPage), nil)
				if err != nil {
					return
				}

				_, _ = CrawlPage(ctx, p, req, itemsChannel, parseFunction)
				// Mark thread as done
				<-threadLimit
				wg.Done()
			}(page)
		}

		// wait for all threads to be done
		wg.Wait()
	}()
}
