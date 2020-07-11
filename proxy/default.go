package proxy

import (
	"context"
	"net/http"

	browser "github.com/eddycjy/fake-useragent"
)

// DefaultProxy is no proxy, it simply makes the request
//
// It will add a default user agent if non is present
type DefaultProxy struct{}

// MakeRequest Makes a request without a proxy
func (d *DefaultProxy) MakeRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", browser.Random())
	}
	return http.DefaultClient.Do(req)
}
