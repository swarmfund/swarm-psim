package issuance

import (
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/tokend/keypair"
)

func CraftIssuanceTX(opt RequestOpt, builder *xdrbuild.Builder, source keypair.Address, signer keypair.Full) *xdrbuild.Transaction {
	return builder.
		Transaction(source).
		Op(xdrbuild.CreateIssuanceRequestOp{
			Reference: opt.Reference,
			Receiver:  opt.Receiver,
			Asset:     opt.Asset,
			Amount:    opt.Amount,
			Details:   opt.Details,
		}).
		Sign(signer)
}
