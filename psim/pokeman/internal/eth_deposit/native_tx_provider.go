package eth_deposit

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

type NativeTxProvider struct {
	horizon *horizon.Connector
	builder *xdrbuild.Builder
	source keypair.Address
	kp eth.Keypair
	signer keypair.Full
	asset string
	balanceID string
	ctx context.Context
}

func NewNativeTxProvider(horizon *horizon.Connector, builder *xdrbuild.Builder, source keypair.Address, kp eth.Keypair, signer keypair.Full, asset string, balanceID string, ctx context.Context) *NativeTxProvider {
	return &NativeTxProvider{
		horizon,
		builder,
		source,
		kp,
		signer,
		asset,
		balanceID,
		ctx,
	}
}

func (n *NativeTxProvider) Send() (bool, error) {
	envelope, err := n.builder.Transaction(n.source).Op(xdrbuild.CreateWithdrawRequestOp{
		Balance: n.balanceID,
		Asset:   n.asset,
		Amount:  2,
		Details: &xdrbuild.ETHWithdrawRequestDetails{
			Address: n.kp.Address().Hex(),
		},
	}).Sign(n.signer).Marshal()
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal withdraw request")
	}

	result := n.horizon.Submitter().Submit(n.ctx, envelope)
	if result.Err != nil {
		return false, errors.Wrap(result.Err, "failed to submit withdraw tx")
	}
	return true, nil
}