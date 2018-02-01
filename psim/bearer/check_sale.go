package bearer

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var errNoSales = errors.New("no sales")

// checkSaleState is create and submit `CheckSaleStateOp`.
func (s *Service) checkSaleState(ctx context.Context) error {
	sales, err := s.checker.GetSales()
	if err != nil {
		return errors.Wrap(err, "failed to get sales")
	}

	if len(sales) == 0 {
		return nil
	}

	info, err := s.checker.GetHorizonInfo()
	if err != nil {
		return errors.Wrap(err, "failed to get horizon info")
	}

	for _, sale := range sales {
		envelope, err := s.checker.BuildTx(info, sale.ID)
		if err != nil {
			return errors.Wrap(err, "failed to marshal tx")
		}

		result := s.checker.SubmitTx(ctx, envelope)
		if result.Err != nil {
			return errors.Wrap(result.Err, "failed to submit tx", logan.F{
				"submit_response_raw":      string(result.RawResponse),
				"submit_response_tx_code":  result.TXCode,
				"submit_response_op_codes": result.OpCodes,
			})
		}
	}

	return nil
}
