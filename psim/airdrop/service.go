package airdrop

import (
	"context"

	"sync"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type TXStreamer interface {
	StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
}

type Service struct {
	log     *logan.Entry
	config  Config
	builder *xdrbuild.Builder

	txSubmitter TXSubmitter
	txStreamer  TXStreamer

	createdAccounts        map[string]struct{}
	generalAccountsCh      chan string
	pendingGeneralAccounts SyncSet
}

func NewService(
	log *logan.Entry,
	config Config,
	builder *xdrbuild.Builder,
	txSubmitter TXSubmitter,
	txStreamer TXStreamer,
) *Service {

	return &Service{
		log:     log,
		config:  config,
		builder: builder,

		txSubmitter: txSubmitter,
		txStreamer:  txStreamer,

		createdAccounts:        make(map[string]struct{}),
		generalAccountsCh:      make(chan string, 100),
		pendingGeneralAccounts: SyncSet{mu: sync.Mutex{}},
	}
}

func (s *Service) Run(ctx context.Context) {
	go s.listenLedgerChangesInfinitely(ctx)

	go s.consumeGeneralAccounts(ctx)

	go s.processPendingGeneralAccounts(ctx)
}
