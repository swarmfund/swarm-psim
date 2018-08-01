package eth_deposit

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/multiplay/go-slack"
	"github.com/multiplay/go-slack/chat"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type Service struct {
	log               *logan.Entry
	foreignTxProvider TxProvider
	slack             slack.Client
	horizon           *horizon.Connector
	builder           *xdrbuild.Builder
	config            Config
}

func NewService(log *logan.Entry, eth TxProvider, slack slack.Client, horizon *horizon.Connector, config Config, builder *xdrbuild.Builder) *Service {
	service := Service{
		log:               log,
		foreignTxProvider: eth,
		slack:             slack,
		horizon:           horizon,
		config:            config,
		builder:           builder,
	}
	return &service
}

// pollBalance will endlessly pollBalanceChange for balance update in config.Asset for config.Source
// and return updated balance value as well as approximate time it took to update
// TODO make sure callers handle ctx close and invalid outputs it will make us generate
func (s *Service) pollBalance(ctx context.Context, current regources.Amount) (updated regources.Amount, took time.Duration) {
	started := time.Now()
	defer func() {
		took = time.Now().Sub(started)
	}()
	running.UntilSuccess(ctx, s.log, "balance-poller", func(i context.Context) (bool, error) {
		balance, err := s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
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

// ensureExternalBinding tries it's best to get you config.Source external system binding data for provided externalSystem
// TODO make sure callers handle ctx close and invalid outputs it will make us generate
func (s *Service) ensureExternalBinding(ctx context.Context, externalSystem int32) (string, error) {
	externalAddr, err := s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
	if err != nil {
		return "", errors.Wrap(err, "failed to get external binding data")
	}
	if externalAddr == nil {
		// seems like account does not have external binding atm, let's fix that
		s.sendTxWithOp(ctx, &xdrbuild.BindExternalSystemAccountIDOp{externalSystem})
		if err != nil {
			return "", errors.Wrap(err, "failed to submit bind tx")
		}

		// probably better to parse tx result here to obtain external binding data,
		// but nobody loves to mess with txresult mess and it's also safer to check explicitly
		running.UntilSuccess(ctx, s.log, "external-data-getter", func(i context.Context) (bool, error) {
			externalAddr, err = s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
			if err != nil {
				return false, errors.Wrap(err, "failed to get external binding data")
			}
			return externalAddr != nil, nil
		}, 5*time.Second, 5*time.Second)
	}
	return *externalAddr, nil
}

func (s *Service) sendTxWithOp(ctx context.Context, op xdrbuild.Operation) error {
	envelope, err := s.builder.Transaction(s.config.Signer).Op(op).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal withdraw request")
	}

	result := s.horizon.Submitter().Submit(ctx, envelope)
	if result.Err != nil {
		return errors.Wrap(err, "failed to submit tx", result.GetLoganFields())
	}
	return nil
}

func (s *Service) sendMessage(message string) {
	msg := fmt.Sprint(message)
	slackMsg := &chat.Message{Text: msg}
	slackMsg.Send(s.slack)
	fmt.Println(msg)
}

func (s *Service) pollBalanceChange(i context.Context, currentBalance regources.Amount) (regources.Amount, error) {
	balance, err := s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get account balance")
	}
	return balance.Balance, nil
}

func (s *Service) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "poke-iter", func(i context.Context) error {
		// it's better to update asset external system type on every iteration in case it might change
		externalSystem, err := internal.GetExternalSystemType(s.horizon.Assets(), s.config.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to get external system type")
		}

		balance, err := s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to get account balance")
		}

		balanceBefore := balance.Balance

		externalAddr, err := s.ensureExternalBinding(ctx, externalSystem)
		if err != nil {
			return errors.Wrap(err, "failed to get external address")
		}
		if !common.IsHexAddress(externalAddr) {
			return errors.New("invalud hex address")
		}

		s.foreignTxProvider.Send(ctx, 5, externalAddr)

		/* at this point we should buksovat, since ETH has been sent */

		// deposit
		depositPollingStarted := time.Now()
		depositTook := func() time.Duration {
			return time.Now().Sub(depositPollingStarted)
		}
		var balanceAfterDeposit regources.Amount
		var balanceChangedOnDeposit bool
		for !balanceChangedOnDeposit {
			if running.IsCancelled(ctx) {
				s.sendMessage(fmt.Sprintf("withdraw polling interrupted after: %s\n", depositTook().String()))
				return nil
			}
			if depositTook() >= 10*time.Minute {
				s.sendMessage(fmt.Sprintf("withdraw polling timed out\n"))
				return nil
			}
			balanceAfterDeposit, err = s.pollBalanceChange(ctx, balanceBefore)
			if err != nil {
				s.sendMessage(fmt.Sprintf("withdraw polling failed with error after: %s\n", depositTook().String()))
				return errors.Wrap(err, "failed to poll balance changes")
			}
			balanceChangedOnDeposit = balanceAfterDeposit != balanceBefore
		}
		if balanceAfterDeposit-balanceBefore != 5 {
			s.sendMessage(fmt.Sprintf("withdraw failed: %s\n", depositTook().String()))
			return nil
		}

		s.sendMessage(fmt.Sprintf("deposit took: %s\n", depositTook().String()))

		/* withdraw flow, could ease on buksovanie for a bit */

		err = s.sendTxWithOp(ctx, xdrbuild.CreateWithdrawRequestOp{
			Balance: balance.BalanceID,
			Asset:   s.config.Asset,
			Amount:  2,
			Details: s.foreignTxProvider.GetWithdrawRequestDetails(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to submit withdraw tx")
		}

		// withdraw
		updatedBalance := balanceAfterDeposit
		withdrawPollingStarted := time.Now()
		withdrawTook := func() time.Duration {
			return time.Now().Sub(withdrawPollingStarted)
		}
		var balanceAfterWithdraw regources.Amount
		var balanceChangedOnWithdraw bool
		for !balanceChangedOnWithdraw {
			if running.IsCancelled(ctx) {
				s.sendMessage(fmt.Sprintf("withdraw polling interrupted after: %s\n", withdrawTook().String()))
				return nil
			}
			if withdrawTook() >= 10*time.Minute {
				s.sendMessage(fmt.Sprintf("withdraw polling timed out\n"))
				return nil
			}
			balanceAfterWithdraw, err = s.pollBalanceChange(ctx, updatedBalance)
			if err != nil {
				s.sendMessage(fmt.Sprintf("withdraw polling failed with error after: %s\n", withdrawTook().String()))
				return errors.Wrap(err, "failed to poll balance changes")
			}
			balanceChangedOnWithdraw = balanceAfterWithdraw != updatedBalance
		}
		if updatedBalance-balanceAfterWithdraw != 2 {
			s.sendMessage(fmt.Sprintf("withdraw failed: %s\n", withdrawTook().String()))
			return nil
		}

		// TODO validate ETH balance

		s.sendMessage(fmt.Sprintf("withdraw took: %s\n", withdrawTook().String()))

		return nil
	}, 10*time.Second, 10*time.Second, 10*time.Second)
}
