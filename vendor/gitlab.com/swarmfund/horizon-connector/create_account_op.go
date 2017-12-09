package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
)

type CreateAccountOp struct {
	AccountID   string
	Referrer    *string
	AccountType xdr.AccountType
}

func (op CreateAccountOp) XDR() (*xdr.Operation, error) {
	var xAccountID, xReferrerID xdr.AccountId
	err := xAccountID.SetAddress(op.AccountID)
	if err != nil {
		return nil, err
	}

	xdrop := &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeCreateAccount,
			CreateAccountOp: &xdr.CreateAccountOp{
				Destination: xAccountID,
				AccountType: op.AccountType,
			},
		},
	}

	if op.Referrer != nil {
		err := xReferrerID.SetAddress(*op.Referrer)
		if err != nil {
			return nil, err
		}
		xdrop.Body.CreateAccountOp.Referrer = &xReferrerID
	}

	return xdrop, nil
}
