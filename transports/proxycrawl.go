package transports

import (
	"fmt"
	"net/http"
	"net/url"
)

// ProxyCrawlTransport passes request through proxycrawl
type ProxyCrawlTransport struct {
	token      string
	maxRetries int
}

// GenerateProxyCrawlTransport builds a new proxy crawl transporter based on a token.
// maxRetries is the number of times we should try re-making the request if proxycrawl returns HTTP 520
func GenerateProxyCrawlTransport(proxyCrawlToken string, maxRetries int) http.RoundTripper {
	return &ProxyCrawlTransport{
		token:      proxyCrawlToken,
		maxRetries: maxRetries,
	}
}

// RoundTrip performs a request using proxy crawl
func (p *ProxyCrawlTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	values := url.Values{}
	values.Add("token", p.token)
	values.Add("url", req.URL.String())
	URL, err := url.Parse(fmt.Sprintf("https://api.proxycrawl.com/?%s", values.Encode()))
	if err != nil {
		return nil, err
	}

	req2 := *req
	req2.URL = URL

	// Keep trying to make requests until proxycrawl succeeds
	for i := 0; i < p.maxRetries; i++ {
		resp, err := http.DefaultClient.Do(&req2)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 520 {
			return resp, err
		}

		// Check context
		select {
		case <-req.Context().Done():
			break
		default:
		}
	}

	return nil, fmt.Errorf("proxycrawl could not make the request without getting a 520.  Try increasing maxRetries")
}
