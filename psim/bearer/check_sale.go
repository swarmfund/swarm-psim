package bearer

import (
	"context"
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"fmt"
	"encoding/json"
	//"golang.org/x/net/html/atom"
)

var errorNoSales = errors.New("no sales")

func obtainSales(horizonClient *horizon.Client) ([]horizon.Sale, error) {
	respBytes, err := horizonClient.Get(fmt.Sprintf("/core_sales"))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get core sales from Horizon")
	}

	var sales []horizon.Sale
	err = json.Unmarshal(respBytes, &sales)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Sales from Horizon response", logan.F{
			"horizon_response": string(respBytes),
		})
	}

	return sales, nil
}

// checkSaleState is create and submit `CheckSaleStateOp`.
func (s *Service) checkSaleState() error {
	sales, err := obtainSales(s.horizon.Client())

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
			Op(xdrbuild.CheckSaleOp{
			SaleID: sale.ID,
		}).
			Sign(s.config.Signer).
			Marshal()

		if err != nil {
			return errors.Wrap(err, "failed to marshal tx")
		}

		result := s.horizon.Submitter().Submit(context.TODO(), envelope)
		if result.Err != nil {
			return errors.Wrap(result.Err, "failed to submit tx")
		}

		if len(result.OpCodes) == 1 && strings.Contains(result.OpCodes[0], xdr.CheckSaleStateResultCodeNotFound.ShortString()) {
			return errorNoSales
		}

		err = errors.Wrap(result.Err, "tx submission failed")
	}

	return err
}