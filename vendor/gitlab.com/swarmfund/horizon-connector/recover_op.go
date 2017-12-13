package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
)

type RecoverOp struct {
	AccountID string
	OldSigner string
	NewSigner string
}

func (op RecoverOp) XDR() (*xdr.Operation, error) {
	var xAccountID, xOldSigner, xNewSigner xdr.AccountId

	err := xAccountID.SetAddress(op.AccountID)
	if err != nil {
		return nil, err
	}

	err = xOldSigner.SetAddress(op.OldSigner)
	if err != nil {
		return nil, err
	}

	err = xNewSigner.SetAddress(op.NewSigner)
	if err != nil {
		return nil, err
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeRecover,
			RecoverOp: &xdr.RecoverOp{
				Account:   xAccountID,
				OldSigner: xdr.PublicKey(xOldSigner),
				NewSigner: xdr.PublicKey(xNewSigner),
			},
		},
	}, nil
}
