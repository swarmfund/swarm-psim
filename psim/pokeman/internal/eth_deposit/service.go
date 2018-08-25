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

// TODO add throttling
// ensureExternalBinding tries it's best to get you config.Source external system binding data for provided externalSystem
func (s *Service) ensureExternalBinding(ctx context.Context, externalSystem int32) (string, error) {
	// FIXME: expose running.NewIncrementalTimer
	//timer := running.NewIncrementalTimer(1, 2, 2)
	externalAddr, err := s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
	if err != nil {
		return "", errors.Wrap(err, "failed to get external binding data")
	}
	if externalAddr == nil {
		// seems like account does not have external binding atm, let's fix that
		err := s.sendTxWithOp(ctx, &xdrbuild.BindExternalSystemAccountIDOp{externalSystem})
		if err != nil {
			return "", errors.Wrap(err, "failed to submit bind tx")
		}

		// probably better to parse tx result here to obtain external binding data,
		// but nobody loves to mess with txresult mess and it's also safer to check explicitly
		currentExternalBindingData := func(i context.Context) error {
			for externalAddr == nil {
				if running.IsCancelled(i) {
					return errors.New("interrupted")
				}
				externalAddr, err = s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
				if err != nil {
					return errors.Wrap(err, "failed to get external binding data")
				}
				// FIXME: expose running.IncrementalTimer.Next()
				//<-timer.Next()
			}
			return nil
		}
		err = currentExternalBindingData(ctx)
		if err != nil {
			return "", errors.Wrap(err, "failed to ensure external binding")
		}
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

func (s *Service) pollBalance(i context.Context, currentBalance regources.Amount) (regources.Amount, error) {
	balance, err := s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get account balance")
	}
	return balance.Balance, nil
}

func (s *Service) ensureBalanceChanged(ctx context.Context, balanceBeforeOperation regources.Amount, changedTo regources.Amount) error {
	// FIXME:
	// timer := running.NewIncrementalTimer(1, 2, 2)

	timedCtx, cancelTimedCtx := context.WithTimeout(ctx, s.config.PollingTimeout)
	defer cancelTimedCtx()

	operationStarted := time.Now()
	operationTook := func() string {
		return time.Now().Sub(operationStarted).String()
	}

	var balanceAfterOperation regources.Amount

	var balanceChangedAfterOperation bool
	for !balanceChangedAfterOperation {
		if running.IsCancelled(timedCtx) {
			if running.IsCancelled(ctx) {
				return errors.New("interrupted")
			}
			s.sendMessage(fmt.Sprintf("balance change polling timed out after: %s\n", operationTook()))
			return errors.New("timed out")
		}

		var err error
		balanceAfterOperation, err = s.pollBalance(timedCtx, balanceBeforeOperation)
		// FIXME:
		//<-timer.Next()
		if err != nil {
			s.sendMessage(fmt.Sprintf("balance change polling failed with error after: %s\n", operationTook()))
			return errors.Wrap(err, "balance polling failed")
		}
		balanceChangedAfterOperation = balanceAfterOperation != balanceBeforeOperation
	}

	if balanceAfterOperation != changedTo {
		s.sendMessage(fmt.Sprintf("op failed, took: %s\n", operationTook()))
		return errors.New("invalid balance change")
	}
	return nil
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

		foreignBalanceBefore, err := s.foreignTxProvider.GetCurrentBalance(ctx)

		balanceBefore := balance.Balance

		externalAddr, err := s.ensureExternalBinding(ctx, externalSystem)
		if err != nil {
			return errors.Wrap(err, "failed to get external address")
		}
		if !common.IsHexAddress(externalAddr) {
			return errors.New("invalid hex address")
		}

		s.foreignTxProvider.Send(ctx, 5, externalAddr)

		/* at this point we should buksovat, since ETH has been sent */

		err = s.ensureBalanceChanged(ctx, balanceBefore, balanceBefore+5)
		if err != nil {
			return errors.Wrap(err, "failed to ensure balance has been changed after deposit")
		}
		balanceBefore = balanceBefore + 5

		foreignBalanceAfter, err := s.foreignTxProvider.GetCurrentBalance(ctx)
		if foreignBalanceAfter.Uint64()-foreignBalanceBefore.Uint64() != 5 {
			return errors.New("foreign balance not changed")
		}

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

		//withdraw
		err = s.ensureBalanceChanged(ctx, balanceBefore, balanceBefore-2)
		if err != nil {
			return errors.Wrap(err, "failed to ensure balance has been changed after withdraw")
		}

		foreignBalanceAfter2, err := s.foreignTxProvider.GetCurrentBalance(ctx)
		if foreignBalanceAfter2.Uint64()-foreignBalanceAfter.Uint64() != 2 {
			return errors.New("foreign balance not changed")
		}

		return nil
	}, 10*time.Second, 10*time.Second, 10*time.Second)
}
