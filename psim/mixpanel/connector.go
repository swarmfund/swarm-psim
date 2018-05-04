package mixpanel

import (
	"time"

	"github.com/dukex/mixpanel"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

const (
	issuanceRequestEvent   = "Issuance Request"
	withdrawalRequestEvent = "Withdrawal Request"
)

type Connector struct {
	mixpanel mixpanel.Mixpanel
}

func NewConnector(token string) *Connector {
	return &Connector{
		mixpanel.New(token, ""),
	}
}

func (c *Connector) IssuanceRequest(id string, ts *time.Time, op *xdr.CreateIssuanceRequestOp) error {
	return c.mixpanel.Track(id, issuanceRequestEvent, &mixpanel.Event{
		IP:        "0",
		Timestamp: ts,
		Properties: map[string]interface{}{
			"reference":        op.Reference,
			"receiver":         op.Request.Receiver.AsString(),
			"asset":            op.Request.Asset,
			"amount":           amount.String(int64(op.Request.Amount)),
			"external_details": op.Request.ExternalDetails,
		},
	})
}

func (c *Connector) WithdrawalRequest(id string, ts *time.Time, op *xdr.CreateWithdrawalRequestOp) error {
	return c.mixpanel.Track(id, withdrawalRequestEvent, &mixpanel.Event{
		IP:        "0",
		Timestamp: ts,
		Properties: map[string]interface{}{
			"balance": op.Request.Balance.AsString(),
			"amount":  amount.String(int64(op.Request.Amount)),
			"details": op.Request.Details,
		},
	})
}
