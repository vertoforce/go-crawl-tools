package crawltools

import (
	"context"
	"net/http"

	"github.com/vertoforce/go-crawl-tools/proxy"
)

// CrawlPage Crawls a single page, parsing out any generic items using the parseFunction, and sends the resulting items to the channel
func CrawlPage(ctx context.Context, p proxy.Proxy, req *http.Request, itemsChannel chan interface{}, parseFunction ParseFunc) (html string, err error) {
	html, err = p.MakeRequest(ctx, req)
	if err != nil {
		return "", err
	}

	// Parse the result
	items := parseFunction(ctx, html)

	// Dump the found items
	for _, item := range items {
		select {
		case itemsChannel <- item:
		case <-ctx.Done():
			return html, ctx.Err()
		}
	}

	return html, nil
}
