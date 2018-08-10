package balancereporter

import (
	"time"

	"github.com/pkg/errors"
	salesforce "gitlab.com/swarmfund/salesforce-connector"
	"gitlab.com/tokend/regources"
)

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

// SalesforceTarget represents a target with salesforce api
type SalesforceTarget struct {
	*salesforce.Connector
}

// NewSalesforceTarget constructs a target
func NewSalesforceTarget(sc *salesforce.Connector) *SalesforceTarget {
	return &SalesforceTarget{
		sc,
	}
}

// SendEvent uses salesforce client connector for sending event to analytics
func (st *SalesforceTarget) SendEvent(event *regources.BalancesReport, swmAmount int64, threshold int64, date *time.Time) error {
	_, err := st.Connector.SendReport(event, swmAmount, threshold, date)
	if err != nil {
		return errors.Wrap(err, "failed to post event")
	}
	return nil
}
