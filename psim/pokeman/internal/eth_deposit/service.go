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
)

type Service struct {
	log           *logan.Entry
	balancePoller BalancePoller
	txProvider    TxProvider
	nativeTxProvider *NativeTxProvider
	ebdProvider ExternalBindingDataProvider
	esBinder ExternalSystemBinder
	currentBalanceProvider CurrentBalanceProvider
	esProvider ExternalSystemProvider
}

// ensureExternalBinding tries it's best to get you config.Source external system binding data for provided externalSystem
// TODO make sure callies handle ctx close and invalid outputs it will make us generate
func (s *Service) ensureExternalBinding(ctx context.Context, externalSystem int32) (string, error) {
	externalAddr, err := s.ebdProvider.CurrentExternalBindingData()
	if err != nil {
		return "", errors.Wrap(err, "failed to get external binding data")
	}

	// seems like account does not have external binding atm, let's fix that
	if externalAddr == nil {
		err := s.esBinder.Bind()
		if err != nil {
			return "", err
		}

		// probably better to parse tx result here to obtain external binding data,
		// but nobody loves to mess with txresult mess and it's also safer to check explicitly
		running.UntilSuccess(ctx, s.log, "external-data-getter", func(i context.Context) (bool, error) {
			externalAddr, err = s.ebdProvider.CurrentExternalBindingData()
			if err != nil {
				return false, errors.Wrap(err, "failed to get external binding data")
			}
			return externalAddr != nil, nil
		}, 5*time.Second, 5*time.Second)
	}
	return *externalAddr, nil
}

func (s *Service) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "poke-iter", func(i context.Context) error {
		// get asset external system type on every iteration, since it might change
		externalSystem, err := s.esProvider.GetExternalSystemType()
		if err != nil {
			return errors.Wrap(err, "failed to get external system type")
		}

		balanceBefore, err := s.currentBalanceProvider.CurrentBalance()
		if err != nil {
			return errors.Wrap(err, "failed to get account balance")
		}

		// get external address
		externalAddr, err := s.ensureExternalBinding(ctx, externalSystem)
		if err != nil {
			return errors.Wrap(err, "failed to get external address")
		}

		_, err = s.txProvider.Send(999, externalAddr)
		if err != nil {
			return errors.From(errors.Wrap(err, "failed to send such an amout to the external address"), logan.F{
				"amount":           999,
				"external_address": externalAddr,
			})
		}

		//
		// at this point we should buksovat, since the asset has been sent
		//

		// get updated balance, hopefully
		currentBalance, depositTook := s.balancePoller.PollBalance(balanceBefore.Balance)

		// TODO ensure balance is updated correctly
		// TODO check if external details are valid

		fmt.Printf("deposit took: %s\n", depositTook.String())

		//
		// withdraw flow, could ease on buksovanie for a bit
		//

		_, err = s.nativeTxProvider.Send()
		if err != nil {
			return errors.Wrap(err, "failed to submit withdraw tx")
		}
		_, withdrawTook := s.balancePoller.PollBalance(currentBalance)

		// TODO validate ETH balance
		// TODO validate tokend balance
		if withdrawTook.Minutes() >= 30 {
			http.Post("https://hooks.slack.com/services/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", "application/json", strings.NewReader("{\"text\":\"Take attention, withdraw took: "+withdrawTook.String()+"\"}"))
		}

		fmt.Printf("withdraw took: %s\n", withdrawTook.String())

		return nil
	}, 10*time.Second, 10*time.Second, 10*time.Second)
}
