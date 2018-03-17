package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64) (*issuance.RequestOpt, error)
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

type EmailProcessor interface {
	Run(context.Context)
	AddEmailAddress(ctx context.Context, emailAddress string)
}

type Service struct {
	log    *logan.Entry
	config Config

	issuanceSubmitter IssuanceSubmitter

	txStreamer        TXStreamer
	accountsConnector AccountsConnector
	usersConnector    UsersConnector

	emailProcessor EmailProcessor

	blackList         map[string]struct{}
	generalAccountsCh chan string
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	txStreamer TXStreamer,
	accountsConnector AccountsConnector,
	usersConnector UsersConnector,
	emailProcessor EmailProcessor,
) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,

		txStreamer:        txStreamer,
		accountsConnector: accountsConnector,
		usersConnector:    usersConnector,

		emailProcessor: emailProcessor,

		blackList:         make(map[string]struct{}),
		generalAccountsCh: make(chan string, 100),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.BlackList {
		s.log.WithField("account_address", accID).Debug("Added Account to BlackList.")
		s.blackList[accID] = struct{}{}
	}

	go s.listenLedgerChanges(ctx)

	go s.consumeGeneralAccounts(ctx)

	go s.emailProcessor.Run(ctx)

	<-ctx.Done()
}
