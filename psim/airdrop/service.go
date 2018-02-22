package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/tokend/keypair"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type Service struct {
	log         *logan.Entry
	builder     *xdrbuild.Builder
	txSubmitter *horizon.Submitter

	source keypair.Address
	signer keypair.Full
}

func NewService(
	log         *logan.Entry,
	builder     *xdrbuild.Builder,
	txSubmitter *horizon.Submitter,
	source keypair.Address,
	signer keypair.Full,
) *Service {

	return &Service{
		log: log,
		builder: builder,
		txSubmitter: txSubmitter,

		source: source,
		signer: signer,
	}
}

// TODO
func (s *Service) Run(ctx context.Context) {
	// TODO
}
