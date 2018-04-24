package wallet_cleaner

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/horizon-connector/types"
)

type Service struct {
	connector      *horizon.Connector
	expireDuration time.Duration
	log            *logan.Entry
}

func New(log *logan.Entry, connector *horizon.Connector, expireDuration time.Duration) *Service {
	return &Service{
		connector:      connector,
		log:            log,
		expireDuration: expireDuration,
	}
}

func (s *Service) Run(ctx context.Context) {
	q := s.connector.Wallets()

	verified := false
	ops := types.GetOpts{
		Verified: &verified,
	}

	do := func() (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(err)
			}
		}()
		wallets, page, err := q.Filter(&ops)
		if err != nil {
			return errors.Wrap(err, "failed to wallets")
		}

		if len(wallets) == 0 {
			ops.Page = nil
			return nil
		}

		for _, wallet := range wallets {
			if wallet.Attributes.LastSentAt == nil {
				// email has not been sent yet
				continue
			}
			if wallet.Attributes.LastSentAt.Add(s.expireDuration).Before(time.Now()) {
				if err := q.Delete(wallet.ID); err != nil {
					return errors.Wrap(err, "failed to delete wallets")
				}

				s.log.WithFields(logan.F{"wallet-id": wallet.ID}).Info("deleted wallet")
			}
		}

		ops.Page = &page

		return nil
	}

	for ; ; time.Sleep(s.expireDuration / 10) {
		if err := do(); err != nil {
			s.log.WithError(err).Error("failed to clean wallets")
		}
	}
}
