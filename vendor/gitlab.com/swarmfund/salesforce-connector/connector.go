package salesforce

import (
	"net/url"

	"gitlab.com/distributed_lab/salesforce"
)

// Connector provides salesforce-interface to be used in PSIM services
type Connector struct {
	client salesforce.Client
}

// NewConnector construct a connector from arguments and gets accessToken
func NewConnector(apiURL *url.URL, secret string, id string, username string, password string) *Connector {
	client := salesforce.NewClient(apiURL, secret, id, username, password)

	return &Connector{
		client: client,
	}
}

// PostObjectResponse holds data received after SendEvent
type PostObjectResponse struct {
	SalesforceID string   `json:"id"`
	Success      bool     `json:"success"`
	Errors       []string `json:"errors"`
}
