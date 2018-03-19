package mrefairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, error)
}

type LedgerStreamer interface {
	Run(ctx context.Context) <-chan airdrop.TimedLedgerChange
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
	ledgerStreamer    LedgerStreamer
	// TODO Consider substituting with some BalanceIDProvider entity.
	accountsConnector airdrop.AccountsConnector
	usersConnector    UsersConnector
	emailProcessor    EmailProcessor

	blackList map[string]struct{}
	// AccountID to Bonus map
	snapshot map[string]*bonusParams
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	ledgerStreamer LedgerStreamer,
	accountsConnector airdrop.AccountsConnector,
	usersConnector UsersConnector,
	emailProcessor EmailProcessor,
) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,
		ledgerStreamer:    ledgerStreamer,
		accountsConnector: accountsConnector,
		usersConnector:    usersConnector,
		emailProcessor:    emailProcessor,

		blackList: make(map[string]struct{}),
		snapshot:  make(map[string]*bonusParams),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.BlackList {
		s.log.WithField("account_address", accID).Debug("Added Account to BlackList.")
		s.blackList[accID] = struct{}{}
	}

	s.processChangesUpToSnapshotTime(ctx)

	go s.emailProcessor.Run(ctx)

	s.payOutSnapshot(ctx)

	<-ctx.Done()
}
