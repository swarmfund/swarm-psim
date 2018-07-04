package salesforce

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/tokend/regources"
)

// EmptyConnector is used for signalizing about special conditions
var EmptyConnector = &Connector{}

// Connector provides salesforce-interface to be used in PSIM services
type Connector struct {
	client *Client
}

// NewConnector construct a connector from arguments and gets accessToken
func NewConnector(apiURL *url.URL, secret string, id string, username string, password string) (*Connector, error) {
	client := NewClient(apiURL, secret, id)
	authResponse, err := client.PostAuthRequest(username, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate while constructing salesforce connector")
	}
	return &Connector{
		client: &Client{
			httpClient:  client.httpClient,
			apiURL:      client.apiURL,
			secret:      client.secret,
			accessToken: authResponse.AccessToken,
			id:          client.id,
		},
	}, nil
}

// SendEvent sends an event from arguments to salesforce
func (c *Connector) SendEvent(sphere string, actionName string, timeString string, actorName string, actorEmail string, investmentAmount int64, investmentCountry string) (*EventResponse, error) {
	return c.client.PostEvent(sphere, actionName, timeString, actorName, actorEmail, investmentAmount, investmentCountry)
}

func (c *Connector) SendReport(report *regources.BalancesReport, swmAmount int64, threshold int64, date *time.Time) (*EventResponse, error) {
	return c.client.PostReport(report.TotalAccountsCount.PositiveBalance, report.TotalAccountsCount.ZeroBalance, swmAmount, threshold, report.TotalAccountsCount.AboveThreshold, report.TotalAccountsCount.BelowThreshold, date)
}
