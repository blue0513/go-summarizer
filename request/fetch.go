package request

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

var ErrURLRequired = fmt.Errorf("URL is required")

func Fetch(url string) (*html.Node, error) {
	if url == "" {
		return nil, ErrURLRequired
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code=%d: %w", resp.StatusCode, err)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("html parse: %w", err)
	}

	return node, nil
}
