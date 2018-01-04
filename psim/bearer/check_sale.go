package bearer

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

func (s *Service) checkSaleState() (xdr.Operation, error) {
	var accountID xdr.AccountId
	err := accountID.SetAddress(s.config.Source.Address())
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "unable to set address")
	}

	return xdr.Operation{
		SourceAccount: &accountID,
		Body: xdr.OperationBody{
			Type:             xdr.OperationTypeCheckSaleState,
			CheckSaleStateOp: &xdr.CheckSaleStateOp{},
		},
	}, nil

}
