package eth_deposit

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type ExternalSystemTypeProvider func() (int32, error)

type CurrentBalanceProv func() (result horizon.Balance, err error)

type Submitter func(ctx context.Context, envelope string) horizon.SubmitResult

type CurrentExternalBindingDataProvider func(externalSystem int32) (*string, error)

type TxBuilder func(op xdrbuild.Operation) (string, error)

type Service struct {
	log                           *logan.Entry
	getExternalSystemType         ExternalSystemTypeProvider
	getCurrentBalance             CurrentBalanceProv
	getCurrentExternalBindingData CurrentExternalBindingDataProvider
	submit                        Submitter
	buildTx                       TxBuilder
	nativeTxProvider *NativeTxProvider
	foreignTxProvider TxProvider
}

// ensureExternalBinding tries it's best to get you config.Source external system binding data for provided externalSystem
// TODO make sure callies handle ctx close and invalid outputs it will make us generate
func (s *Service) ensureExternalBinding(ctx context.Context, externalSystem int32) (string, error) {
	externalAddr, err := s.getCurrentExternalBindingData(externalSystem)
	if err != nil {
		return "", errors.Wrap(err, "failed to get external binding data")
	}

	// seems like account does not have external binding atm, let's fix that
	if externalAddr == nil {
		envelope, err := s.buildTx(&xdrbuild.BindExternalSystemAccountIDOp{externalSystem})
		if err != nil {
			return "", errors.Wrap(err, "failed to marshal bind tx")
		}

		result := s.submit(context.Background(), envelope)
		if result.Err != nil {
			return "", errors.Wrap(result.Err, "failed to submit bind tx", result.GetLoganFields())
		}

		// probably it is better to parse tx result here to obtain external binding data,
		// but nobody loves to mess with tx result mess and it's also safer to check explicitly
		running.UntilSuccess(ctx, s.log, "external-data-getter", func(i context.Context) (bool, error) {
			externalAddr, err = s.getCurrentExternalBindingData(externalSystem)
			if err != nil {
				return false, errors.Wrap(err, "failed to get external binding data")
			}
			return externalAddr != nil, nil
		}, 5*time.Second, 5*time.Second)
	}
	return *externalAddr, nil
}

// pollBalance will endlessly poll for balance update in config.Asset for config.Source
// and return updated balance value as well as approximate time it took to update
// TODO make sure callees handle ctx close and invalid outputs it will make us generate
func (s *Service) pollBalance(ctx context.Context, current regources.Amount, timeout time.Duration) (updated regources.Amount, took time.Duration) {
	started := time.Now()
	defer func() {
		took = time.Now().Sub(started)
	}()
	running.UntilSuccess(ctx, s.log, "balance-poller", func(i context.Context) (bool, error) {
		if time.Now().Sub(started) >= timeout {
			return false, errors.New("timed out")
		}

		balance, err := s.getCurrentBalance()
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

func (s *Service) run(timeout time.Duration) func(ctx context.Context) {
	return func(ctx context.Context) {
		running.WithBackOff(ctx, s.log, "poke-iter", func(i context.Context) error {
			// getting asset external system type on every iteration, since it might change
			externalSystem, err := s.getExternalSystemType()
			if err != nil {
				return errors.Wrap(err, "failed to get external system type")
			}

			balanceBefore, err := s.getCurrentBalance()
			if err != nil {
				return errors.Wrap(err, "failed to get account balance")
			}

			externalAddr, err := s.ensureExternalBinding(ctx, externalSystem)
			if err != nil {
				return errors.Wrap(err, "failed to get external address")
			}

			_, err = s.foreignTxProvider.Send(ctx, 2000, externalAddr)
			if err != nil {
				return errors.From(errors.Wrap(err, "failed to send such an amout to the external address"), logan.F{
					"amount":           999,
					"external_address": externalAddr,
				})
			}

			/* at this point we should buksovat, since the asset has been sent */

			/*currentBalance*/_, depositTook := s.pollBalance(ctx, balanceBefore.Balance, timeout)

			// TODO ensure balance is updated correctly
			// TODO check if external details are valid

			fmt.Printf("deposit took: %s\n", depositTook.String())
			if depositTook >= timeout {
				http.Post("https://hooks.slack.com/services/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", "application/json", strings.NewReader("{\"text\":\"Take attention, deposit took: "+depositTook.String()+"\"}"))
			}

			/* withdraw flow, could ease on buksovanie for a bit */

			_, err = s.nativeTxProvider.Send(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to submit withdraw tx")
			}
			_, withdrawTook := s.pollBalance(ctx, balanceBefore.Balance, timeout)

			// TODO validate ETH balance
			// TODO validate tokend balance
			fmt.Printf("withdraw took: %s\n", withdrawTook.String())
			if withdrawTook >= timeout {
				http.Post("https://hooks.slack.com/services/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", "application/json", strings.NewReader("{\"text\":\"Take attention, withdraw took: "+withdrawTook.String()+"\"}"))
			}

			return nil
		}, 10*time.Second, 10*time.Second, 10*time.Second)
	}
}

func (s *Service) Run(ctx context.Context) {
	s.run(1*time.Minute)
}

type timedService struct {
	run func(ctx context.Context)
}

func (k *timedService) Run(ctx context.Context) {
	k.run(ctx)
}

func (s *Service) WithTimeout(timeout time.Duration) (*timedService) {
	return &timedService{s.run(timeout)}
}
