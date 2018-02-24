package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type TXStreamer interface {
	StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
	Users(ctx context.Context) (<-chan horizon.User, <-chan error)
}

type AccountsConnector interface {
	Balances(address string) ([]horizon.Balance, error)
}

type Service struct {
	log     *logan.Entry
	config  Config
	builder *xdrbuild.Builder

	txSubmitter       TXSubmitter
	txStreamer        TXStreamer
	usersConnector    UsersConnector
	accountsConnector AccountsConnector

	notificator *notificator.Connector

	createdAccounts        map[string]struct{}
	generalAccountsCh      chan string
	pendingGeneralAccounts SyncSet

	emails SyncSet
}

func NewService(
	log *logan.Entry,
	config Config,
	builder *xdrbuild.Builder,
	txSubmitter TXSubmitter,
	txStreamer TXStreamer,
	usersConnector UsersConnector,
	accountsConnector AccountsConnector,
	notificator *notificator.Connector,
) *Service {

	return &Service{
		log:     log,
		config:  config,
		builder: builder,

		txSubmitter:       txSubmitter,
		txStreamer:        txStreamer,
		usersConnector:    usersConnector,
		accountsConnector: accountsConnector,

		notificator: notificator,

		createdAccounts:        make(map[string]struct{}),
		generalAccountsCh:      make(chan string, 100),
		pendingGeneralAccounts: NewSyncSet(),

		emails: NewSyncSet(),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	go s.listenLedgerChangesInfinitely(ctx)

	go s.consumeGeneralAccounts(ctx)

	go s.processPendingGeneralAccounts(ctx)

	go s.processEmails(ctx)

	<-ctx.Done()
}
