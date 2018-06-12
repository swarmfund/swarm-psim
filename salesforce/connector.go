package salesforce

import "net/url"

// EmptyConnector is used for signalizing about special conditions
var EmptyConnector = &Connector{}

// Connector provides salesforce-interface to be used in PSIM services
type Connector struct {
	client *Client
}

// NewConnector construct a connector from arguments and gets accessToken
func NewConnector(apiURL *url.URL, secret string, id string) (*Connector, error) {
	client := NewClient(apiURL, secret, id)
	authResponse, err := client.PostAuthRequest()
	if err != nil {
		return EmptyConnector, err
	}
	client.accessToken = authResponse.AccessToken
	return &Connector{
		client,
	}, nil
}

func (c *Connector) WithUserData(username string, password string) *Connector {
	return &Connector{
		c.client.WithUserData(username, password),
	}
}

// SendEvent sends an event from arguments to salesforce
func (c *Connector) SendEvent(sphere string, actionName string, timeString string, actorName string, actorEmail string, investmentAmount int64, investmentCountry string) (*EventResponse, error) {
	return c.client.PostEvent(sphere, actionName, timeString, actorName, actorEmail, investmentAmount, investmentCountry)
}
