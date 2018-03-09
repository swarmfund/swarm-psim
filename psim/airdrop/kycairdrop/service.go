package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/airdrop"
)

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type TXStreamer interface {
	StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
}

type AccountsConnector interface {
	Balances(address string) ([]horizon.Balance, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type Service struct {
	log     *logan.Entry
	config  Config
	builder *xdrbuild.Builder

	txSubmitter       TXSubmitter
	txStreamer        TXStreamer
	accountsConnector AccountsConnector
	usersConnector    UsersConnector
	notificator       airdrop.NotificatorConnector

	blackList         map[string]struct{}
	generalAccountsCh chan string

	emails airdrop.SyncSet
}

func NewService(
	log *logan.Entry,
	config Config,
	builder *xdrbuild.Builder,
	txSubmitter TXSubmitter,
	txStreamer TXStreamer,
	accountsConnector AccountsConnector,
	usersConnector UsersConnector,
	notificator airdrop.NotificatorConnector,
) *Service {

	return &Service{
		log:     log,
		config:  config,
		builder: builder,

		txSubmitter:       txSubmitter,
		txStreamer:        txStreamer,
		accountsConnector: accountsConnector,
		usersConnector:    usersConnector,

		notificator: notificator,

		blackList:         make(map[string]struct{}),
		generalAccountsCh: make(chan string, 100),

		emails: airdrop.NewSyncSet(),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.BlackList {
		s.blackList[accID] = struct{}{}
	}

	go s.listenLedgerChanges(ctx)

	go s.consumeGeneralAccounts(ctx)

	go s.processEmails(ctx)

	<-ctx.Done()
}
