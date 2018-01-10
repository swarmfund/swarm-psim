package bearer

import (
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
)

var errorNoSales = errors.New("no sales")

func createCheckSaleOp(address string) (*xdr.Operation, error) {
	var accountID xdr.AccountId
	err := accountID.SetAddress(address)
	if err != nil {
		return nil, err
	}

	return &xdr.Operation{
		SourceAccount: &accountID,
		Body: xdr.OperationBody{
			Type:             xdr.OperationTypeCheckSaleState,
			CheckSaleStateOp: &xdr.CheckSaleStateOp{},
		},
	}, nil
}

// checkSaleState is create and submit `CheckSaleStateOp`.
func (s *Service) checkSaleState() error {
	checkSale, err := createCheckSaleOp(s.config.Source.Address())
	if err != nil {
		return errors.Wrap(err, "unable to create checkSaleOp")
	}

	err = s.submitOperation(*checkSale)
	if err == nil {
		return nil
	}

	serr, ok := errors.Cause(err).(horizon.SubmitError)
	if !ok {
		return errors.Wrap(err, "unable to submit tx")
	}

	for _, code := range serr.OperationCodes() {
		if strings.Contains(code, xdr.CheckSaleStateResultCodeNoSalesFound.ShortString()) {
			return errorNoSales
		}
	}

	return errors.Wrap(err, "tx submission failed")
}
