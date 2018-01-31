package bearer

import (
	"context"
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
)

var errNoSales = errors.New("no sales")

// checkSaleState is create and submit `CheckSaleStateOp`.
func (s *Service) checkSaleState(ctx context.Context) error {
	sales, err := s.obtainSales()

	var builder *xdrbuild.Builder
	{
		info, err := s.horizon.Info()
		if err != nil {
			return errors.Wrap(err, "failed to get horizon info")
		}
		builder = xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	}

	for _, sale := range sales {
		envelope, err := builder.
			Transaction(s.config.Source).
			Op(xdrbuild.CheckSaleOp{SaleID: sale.ID,}).
			Sign(s.config.Signer).
			Marshal()

		if err != nil {
			return errors.Wrap(err, "failed to marshal tx")
		}

		result := s.horizon.Submitter().Submit(ctx, envelope)

		if len(result.OpCodes) == 1 && strings.Contains(result.OpCodes[0], xdr.CheckSaleStateResultCodeNotFound.ShortString()) {
			return errNoSales
		}

		if result.Err != nil {
			return errors.Wrap(result.Err, "failed to submit tx", logan.F{
				"tx_code": result.TXCode,
			})
		}
	}

	return err
}
