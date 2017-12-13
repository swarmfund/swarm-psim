package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
)

type ManageAccountOp struct {
	AccountID     string
	AccountType   xdr.AccountType
	AddReasons    xdr.BlockReasons
	RemoveReasons xdr.BlockReasons
}

func (op ManageAccountOp) XDR() (*xdr.Operation, error) {
	var xAccountID xdr.AccountId
	err := xAccountID.SetAddress(op.AccountID)
	if err != nil {
		return nil, err
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageAccount,
			ManageAccountOp: &xdr.ManageAccountOp{
				Account:              xAccountID,
				BlockReasonsToAdd:    xdr.Uint32(op.AddReasons),
				BlockReasonsToRemove: xdr.Uint32(op.RemoveReasons),
				AccountType:          xdr.AccountType(op.AccountType),
			},
		},
	}, nil
}
