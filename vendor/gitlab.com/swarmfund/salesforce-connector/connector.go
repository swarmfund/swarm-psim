package salesforce

import (
	"net/url"
	"time"

	"gitlab.com/distributed_lab/salesforce"
	"gitlab.com/tokend/regources"
)

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

// EmptyConnector is used for signalizing about special conditions
var EmptyConnector = &Connector{}

// Connector provides salesforce-interface to be used in PSIM services
type Connector struct {
	client salesforce.Client
}

// NewConnector construct a connector from arguments and gets accessToken
func NewConnector(apiURL *url.URL, secret string, id string, username string, password string) (*Connector, error) {
	client := salesforce.NewClient(apiURL, secret, id, username, password)

	return &Connector{
		client: client,
	}, nil
}

// SendEvent sends an event from arguments to salesforce
func (c *Connector) SendEvent(sphere string, actionName string, occuredAt time.Time, actorName string, actorEmail string, investmentAmount int64, depositAmount int64, depositCurrency, referral, country string) (*EventResponse, error) {
	return c.PostEvent(sphere, actionName, occuredAt.Format(salesforceTimeLayout), actorName, actorEmail, investmentAmount, depositAmount, depositCurrency, referral, country)
}

// SendReport is used to send balance report to its own endpoint
func (c *Connector) SendReport(report *regources.BalancesReport, swmAmount int64, threshold int64, date *time.Time) (*EventResponse, error) {
	return c.PostReport(report.TotalAccountsCount.PositiveBalance, report.TotalAccountsCount.ZeroBalance, swmAmount, threshold, report.TotalAccountsCount.AboveThreshold, report.TotalAccountsCount.BelowThreshold, date)
}
