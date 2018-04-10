package wallet_cleaner

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
)

type Service struct {
	connector *horizon.Connector
	config    Config
	log       *logan.Entry
}

func New(log *logan.Entry, connector *horizon.Connector, config Config) *Service {
	return &Service{
		connector: connector,
		log:       log,
		config:    config,
	}
}

func (s *Service) Run(ctx context.Context) {
	q := s.connector.Wallets()

	verified := false
	ops := types.GetOpts{
		Verified: &verified,
	}

	expireDuration := s.config.ExpireDuration

	do := func() (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(err)
			}
		}()
		tokens, _, err := q.Filter(&ops)
		if err != nil {
			return errors.Wrap(err, "failed to get email tokens")
		}

		if len(tokens) == 0 {
			return nil
		}

		for _, token := range tokens {
			if token.Attributes.LastSentAt == nil {
				// email has not been sent yet
				continue
			}
			if token.Attributes.LastSentAt.Add(expireDuration).Before(time.Now()) {
				if err := q.Delete(token.ID); err != nil {
					return errors.Wrap(err, "failed to delete wallets")
				}

				s.log.WithFields(logan.F{"wallet-cleaner": "deleted wallet"}).Info(token.ID)
			}
		}

		return nil
	}

	for ; ; time.Sleep(expireDuration / 10) {
		if err := do(); err != nil {
			s.log.WithError(err).Error("failed to clean wallets")
		}
	}
}
