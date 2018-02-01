package bearer

import (
	"context"

	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type SalesStateChecker struct {
	connector *horizon.Connector
	config    Config
}

func NewSalesStateChecker(connector *horizon.Connector, config Config) *SalesStateChecker {
	return &SalesStateChecker{
		connector: connector,
		config:    config,
	}
}

func (ssc *SalesStateChecker) GetSales() ([]horizon.Sale, error) {
	return ssc.connector.Sales().Sales()
}

func (ssc *SalesStateChecker) GetHorizonInfo() (info *horizon.Info, err error) {
	return ssc.connector.Info()
}

func (ssc *SalesStateChecker) BuildTx(info *horizon.Info, saleID uint64) (string, error) {
	builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	envelope, err := builder.
					 Transaction(ssc.config.Source).
					 Op(xdrbuild.CheckSaleOp{SaleID: saleID}).
					 Sign(ssc.config.Signer).
					 Marshal()

	return envelope, err
}

func (ssc *SalesStateChecker) SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult {
	return ssc.connector.Submitter().Submit(ctx, envelope)
}