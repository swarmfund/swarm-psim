package internal

import "net/http"

// Client exists for testing purpose only
type Client interface {
	Get(string) (*http.Response, error)
}
