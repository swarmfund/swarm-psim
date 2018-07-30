package eth_deposit

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/regources"
)

type BalancePoller interface {
	PollBalance(current regources.Amount) (updated regources.Amount, took time.Duration)
}

type balancePoller struct {
	ctx                    context.Context
	logger                 *logan.Entry
	timeout                time.Duration
	currentBalanceProvider CurrentBalanceProvider
}

func NewBalancePoller(ctx context.Context, logger *logan.Entry, timout time.Duration, currentBalanceProvider CurrentBalanceProvider) BalancePoller {
	return &balancePoller{
		ctx,
		logger,
		timout,
		currentBalanceProvider,
	}
}

// pollBalance will endlessly poll for balance update in config.Asset for config.Source
// and return updated balance value as well as approximate time it took to update
// TODO make sure callies handle ctx close and invalid outputs it will make us generate
func (b *balancePoller) PollBalance(current regources.Amount) (updated regources.Amount, took time.Duration) {
	started := time.Now()
	defer func() {
		took = time.Now().Sub(started)
	}()
	running.UntilSuccess(b.ctx, b.logger, "balance-poller", func(i context.Context) (bool, error) {
		if time.Now().Sub(started) >= b.timeout {
			return false, errors.New("timed out")
		}

		balance, err := b.currentBalanceProvider.CurrentBalance()
		if err != nil {
			return false, errors.Wrap(err, "failed to get account balance")
		}
		if current != balance.Balance {
			return true, nil
		}

		updated = balance.Balance
		return false, nil
	}, 5*time.Second, 5*time.Second)
	return updated, took
}