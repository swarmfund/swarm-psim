package telegram

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, bool, error)
}

type BalanceIDProvider interface {
	GetBalanceID(accAddress, asset string) (*string, error)
}

type Service struct {
	log    *logan.Entry
	config Config

	issuanceSubmitter IssuanceSubmitter
	balanceIDProvider BalanceIDProvider

	blackList map[string]struct{}
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	balanceIDProvider BalanceIDProvider) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,
		balanceIDProvider: balanceIDProvider,

		blackList: make(map[string]struct{}),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.BlackList {
		s.log.WithField("account_address", accID).Debug("Added Account to BlackList.")
		s.blackList[accID] = struct{}{}
	}

	<-ctx.Done()
}
