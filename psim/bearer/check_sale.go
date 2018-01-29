package bearer

import (
	"context"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
)

var errorNoSales = errors.New("no sales")

// checkSaleState is create and submit `CheckSaleStateOp`.
func (s *Service) checkSaleState() error {
	var builder *xdrbuild.Builder
	{
		info, err := s.horizon.Info()
		if err != nil {
			return errors.Wrap(err, "failed to get horizon info")
		}
		builder = xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	}

	envelope, err := builder.
		Transaction(s.config.Source).
		Op(xdrbuild.CheckSaleOp{}).
		Sign(s.config.Signer).
		Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		return errors.Wrap(err, "failed to submit tx")
	}

	if len(result.OpCodes) == 1 && strings.Contains(result.OpCodes[0], xdr.CheckSaleStateResultCodeNoSalesFound.ShortString()) {
		return errorNoSales
	}

	return errors.Wrap(result.Err, "tx submission failed")
}
