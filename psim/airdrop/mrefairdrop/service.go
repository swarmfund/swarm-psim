package mrefairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/swarmfund/psim/psim/lchanges"
	"gitlab.com/tokend/horizon-connector"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, bool, error)
}

type LedgerStreamer interface {
	GetStream() <-chan lchanges.TimedLedgerChange
	Run(ctx context.Context, cursor string)
}

type BalanceIDProvider interface {
	GetBalanceID(accAddress, asset string) (*string, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type accountsConnector interface {
	ByAddress(address string) (*horizon.Account, error)
}

type usaChecker interface {
	CheckIsUSA(acc horizon.Account) (bool, error)
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
	balanceIDProvider BalanceIDProvider
	usersConnector    UsersConnector
	accountsConnector accountsConnector
	usaChecker        usaChecker
	emailProcessor    EmailProcessor

	blackList map[string]struct{}
	snapshot  map[string]*bonusParams // AccountID to Bonus map
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	ledgerStreamer LedgerStreamer,
	balanceIDProvider BalanceIDProvider,
	usersConnector UsersConnector,
	accountsConnector accountsConnector,
	usaChecker usaChecker,
	emailProcessor EmailProcessor,
) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,
		ledgerStreamer:    ledgerStreamer,
		balanceIDProvider: balanceIDProvider,
		usersConnector:    usersConnector,
		accountsConnector: accountsConnector,
		usaChecker:        usaChecker,
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
