package crawltools

import "context"

// ParseFunc takes a html page and parses it into some useful structure
type ParseFunc func(ctx context.Context, html string) (items []interface{})

// TotalPagesFunc Takes the html of a page and returns the total pages for paged crawls (total page count)
type TotalPagesFunc func(ctx context.Context, html string) int64
