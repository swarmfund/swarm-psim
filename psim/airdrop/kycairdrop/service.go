package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/swarmfund/psim/psim/lchanges"
	horizon "gitlab.com/tokend/horizon-connector"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, bool, error)
}

type LedgerStreamer interface {
	Run(ctx context.Context, cursor string)
	GetStream() <-chan lchanges.TimedLedgerChange
}

type AccountsConnector interface {
	airdrop.AccountsConnector
	ByAddress(address string) (*horizon.Account, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type ReferencesProvider interface {
	References(accountID string) ([]horizon.Reference, error)
}

type USAChecker interface {
	CheckIsUSA(acc horizon.Account) (bool, error)
}

type EmailProcessor interface {
	Run(context.Context)
	AddEmailAddress(ctx context.Context, emailAddress string)
}

type Service struct {
	log    *logan.Entry
	config Config

	issuanceSubmitter  IssuanceSubmitter
	ledgerStreamer     LedgerStreamer
	accountsConnector  AccountsConnector
	usersConnector     UsersConnector
	usaChecker         USAChecker
	referencesProvider ReferencesProvider

	emailProcessor EmailProcessor

	blackList          map[string]struct{}
	generalAccountsCh  chan string
	existingReferences []string
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	ledgerStreamer LedgerStreamer,
	accountsConnector AccountsConnector,
	usersConnector UsersConnector,
	usaChecker USAChecker,
	referencesProvider ReferencesProvider,
	emailProcessor EmailProcessor,
) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,

		ledgerStreamer:     ledgerStreamer,
		accountsConnector:  accountsConnector,
		usersConnector:     usersConnector,
		usaChecker:         usaChecker,
		referencesProvider: referencesProvider,

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

	s.fetchAllReferences(ctx)
	s.log.Infof("Fetched (%d) References.", len(s.existingReferences))

	go s.listenLedgerChanges(ctx)
	go s.consumeGeneralAccounts(ctx)
	go s.emailProcessor.Run(ctx)

	<-ctx.Done()
}
