package crawltools

import "context"

// ParseFunc takes a html page and parses it into some useful structure
type ParseFunc func(ctx context.Context, html string) (items []interface{})

// TotalPagesFunc Takes the html of a page and returns the total pages for paged crawls (total page count).
// This is because the total page count is usually listed on the page and we need to know how many pages there
// are total to determine when to stop.
type TotalPagesFunc func(ctx context.Context, html string) int64
