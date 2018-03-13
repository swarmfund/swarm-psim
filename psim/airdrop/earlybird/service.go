package earlybird

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/airdrop"
)

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

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type Service struct {
	log     *logan.Entry
	config  Config
	builder *xdrbuild.Builder

	txSubmitter       TXSubmitter
	txStreamer        TXStreamer
	usersConnector    UsersConnector
	accountsConnector AccountsConnector

	notificator airdrop.NotificatorConnector

	createdAccounts        map[string]struct{}
	generalAccountsCh      chan string
	pendingGeneralAccounts airdrop.SyncSet

	emails airdrop.SyncSet
}

func NewService(
	log *logan.Entry,
	config Config,
	builder *xdrbuild.Builder,
	txSubmitter TXSubmitter,
	txStreamer TXStreamer,
	usersConnector UsersConnector,
	accountsConnector AccountsConnector,
	notificator airdrop.NotificatorConnector,
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
		pendingGeneralAccounts: airdrop.NewSyncSet(),

		emails: airdrop.NewSyncSet(),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.WhiteList {
		s.log.WithField("account_address", accID).Debug("Added created Account from WhiteList.")
		s.createdAccounts[accID] = struct{}{}
	}

	go s.listenLedgerChangesInfinitely(ctx)

	go s.consumeGeneralAccounts(ctx)

	go s.processPendingGeneralAccounts(ctx)

	go s.processEmails(ctx)

	<-ctx.Done()
}
