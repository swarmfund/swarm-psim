package salesforce

import (
	"net/http"
	"net/url"
)

type Client interface {
	PostObject(json []byte, endpoint *url.URL) (statusCode int, body []byte, err error)
}

// NewClient constructs a salesforce Client from arguments and
func NewClient(apiURL *url.URL, secret, id, username, password string) Client {
	return &client{
		httpClient: http.DefaultClient,
		apiURL:     apiURL,
		secret:     secret,
		id:         id,
		username:   username,
		password:   password,
	}
}
