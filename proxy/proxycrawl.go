package proxy

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vertoforce/go-proxycrawl"
)

// ProxyCrawlProxy filters requests through proxycrawl
type ProxyCrawlProxy struct {
	// Number of times to try crawling the page before giving up
	TryCount int
	// The timeout for waiting on proxycrawl before ditching that request and making a new one
	PageTimeout time.Duration
	// The amount of time to tell proxycrawl to wait to load the page (for javascript requests)
	PageWait time.Duration

	proxyCrawlClient *proxycrawl.Client
}

// MakeRequest through proxycrawl.
// It will try over and over TryCount times until proxycrawl returns 200.
// It will use javascript or non-javascript, whatever it is configured to use.
func (p *ProxyCrawlProxy) MakeRequest(ctx context.Context, req *http.Request) (html string, err error) {
	for i := 0; i < p.TryCount; i++ { // Try this many times before giving up on this page
		ctxProxyCrawl, cancel := context.WithTimeout(ctx, p.PageTimeout)
		resp, err := p.proxyCrawlClient.MakeRequest(ctxProxyCrawl, &proxycrawl.RequestParameters{
			URL:      req.URL.String(),
			PageWait: p.PageWait.Milliseconds(),
		}, proxycrawl.JavascriptRequest)
		if err != nil {
			cancel()

			// If the context was canceled, break out of trying this page
			if errors.Is(err, context.Canceled) {
				return "", ctx.Err()
			}

			// Sleep and try again
			time.Sleep(time.Second * 3)
			continue
		}

		if resp.StatusCode != 200 {
			// This didn't work, try again
			cancel()
			time.Sleep(time.Second * 3)
			continue
		}

		// Success, return response
		// Read all body
		body, err := ioutil.ReadAll(resp.Body)
		cancel()
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("failed to get successful response after max tries (%d)", p.TryCount)
}
