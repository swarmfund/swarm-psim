package bearer

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/tokend/keypair"
)

// CheckSalesStateHelper is a particular implementation of CheckSalesStateHelperInterface
type CheckSalesStateHelper struct {
	SalesQ
	connector *horizon.Connector
	builder   *xdrbuild.Builder
	source    keypair.Address
	signer    keypair.Full
}

func NewCheckSalesStateHelper(
	connector *horizon.Connector, builder *xdrbuild.Builder, source keypair.Address, signer keypair.Full,
) *CheckSalesStateHelper {
	return &CheckSalesStateHelper{
		SalesQ:    connector.Sales(),
		connector: connector,
		builder:   builder,
		source:    source,
		signer:    signer,
	}
}

func (h *CheckSalesStateHelper) CloseSale(id uint64) (bool, error) {
	envelope, err := h.builder.
		Transaction(h.source).
		Op(xdrbuild.CheckSaleOp{
			SaleID: id,
		}).
		Sign(h.signer).
		Marshal()
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal tx")
	}
	result := h.connector.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		if len(result.OpCodes) == 1 {
			switch result.OpCodes[0] {
			case "op_not_ready":
				return false, nil
			}
		}
		return false, errors.Wrap(result.Err, "failed to submit tx", result.GetLoganFields())
	}
	return true, nil
}
