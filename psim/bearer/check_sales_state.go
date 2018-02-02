package bearer

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// CheckSalesState allows to perform check sale state operation for all sales in core DB
func (s *Service) checkSalesState(ctx context.Context) error {
	sales, err := s.helper.GetSales()
	if err != nil {
		return errors.Wrap(err, "failed to get sales")
	}

	if len(sales) == 0 {
		return nil
	}

	info, err := s.helper.GetHorizonInfo()
	if err != nil {
		return errors.Wrap(err, "failed to get horizon info")
	}

	for _, sale := range sales {
		envelope, err := s.helper.BuildTx(info, sale.ID)
		if err != nil {
			return errors.Wrap(err, "failed to build tx", logan.F{
				"sale_id": sale.ID,
			})
		}

		result := s.helper.SubmitTx(ctx, envelope)
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
