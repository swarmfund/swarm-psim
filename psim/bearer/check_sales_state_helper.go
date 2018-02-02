package bearer

import (
	"context"

	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

// CheckSalesStateHelper is a particular implementation of CheckSalesStateHelperInterface
type CheckSalesStateHelper struct {
	connector *horizon.Connector
	config    Config
}

// NewCheckSalesStateHelper is a constructor for CheckSalesStateHelper
func NewCheckSalesStateHelper(connector *horizon.Connector, config Config) *CheckSalesStateHelper {
	return &CheckSalesStateHelper{
		connector: connector,
		config:    config,
	}
}

// GetSales returns sales from core DB
func (ssc *CheckSalesStateHelper) GetSales() ([]horizon.Sale, error) {
	return ssc.connector.Sales().Sales()
}

// GetHorizonInfo retrieves horizon info using horizon-connector
func (ssc *CheckSalesStateHelper) GetHorizonInfo() (info *horizon.Info, err error) {
	return ssc.connector.Info()
}

// BuildTx builds transaction with check sale state operation
func (ssc *CheckSalesStateHelper) BuildTx(info *horizon.Info, saleID uint64) (string, error) {
	builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	envelope, err := builder.
		Transaction(ssc.config.Source).
		Op(xdrbuild.CheckSaleOp{SaleID: saleID}).
		Sign(ssc.config.Signer).
		Marshal()

	return envelope, err
}

// SubmitTx submits transaction to horizon, returns submit result
func (ssc *CheckSalesStateHelper) SubmitTx(ctx context.Context, envelope string) horizon.SubmitResult {
	return ssc.connector.Submitter().Submit(ctx, envelope)
}
