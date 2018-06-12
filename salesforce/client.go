package salesforce

import (
	"net/http"
	"net/url"
)

var authEndpointURL = &url.URL{
	Path: "services/oauth2/token",
}

// Client is a custom salesforce client implementation
type Client struct {
	httpClient  *http.Client
	apiURL      *url.URL
	username    string
	password    string
	secret      string
	accessToken string
	id          string
}

// NewClient constructs a salesforce Client from arguments and
func NewClient(apiURL *url.URL, secret string, id string) *Client {
	salesforceClient := &Client{
		httpClient: http.DefaultClient,
		apiURL:     apiURL,
		secret:     secret,
		id:         id,
	}
	return salesforceClient
}

func (c *Client) WithUserData(username string, password string) *Client {
	return &Client{
		httpClient:  c.httpClient,
		apiURL:      c.apiURL,
		username:    username,
		password:    password,
		secret:      c.secret,
		accessToken: c.accessToken,
		id:          c.id,
	}
}
