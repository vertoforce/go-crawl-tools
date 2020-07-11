package proxy

import (
	"context"
	"net/http"
)

// Proxy is some escrow for making a request (through a proxy, or not)
type Proxy interface {
	MakeRequest(ctx context.Context, req *http.Request) (string, error)
}
