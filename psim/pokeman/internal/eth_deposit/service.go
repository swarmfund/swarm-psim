package eth_deposit

import (
	"context"
	"fmt"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	types2 "gitlab.com/tokend/horizon-connector/types"
)

type Service struct {
	log     *logan.Entry
	eth     *ethclient.Client
	horizon *horizon.Connector
	builder *xdrbuild.Builder
	config  Config
}

func NewService(log *logan.Entry, eth *ethclient.Client, horizon *horizon.Connector, config Config, builder *xdrbuild.Builder) *Service {
	service := Service{
		log:     log,
		eth:     eth,
		horizon: horizon,
		config:  config,
		builder: builder,
	}
	return &service

}

func (s *Service) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "deposit-iter", func(i context.Context) error {
		// get asset external system type
		// it's better to update it on every iteration in case it might change
		externalSystem, err := internal.GetExternalSystemType(s.horizon.Assets(), s.config.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to get external system type")
		}

		// get current account balance
		balanceBefore, err := s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to get balance")
		}

		// get external address
		externalAddr, err := s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
		if err != nil {
			return errors.Wrap(err, "failed to get external binding data")
		}
		if externalAddr == nil {
			// seems like account does not have external binding atm, let's fix that
			envelope, err := s.builder.Transaction(s.config.Source).Op(
				&xdrbuild.BindExternalSystemAccountIDOp{externalSystem},
			).Sign(s.config.Signer).Marshal()
			if err != nil {
				return errors.Wrap(err, "failed to marshal bind tx")
			}

			result := s.horizon.Submitter().Submit(context.Background(), envelope)
			if result.Err != nil {
				return errors.Wrap(result.Err, "failed to submit bind tx", result.GetLoganFields())
			}

			// probably better to parse tx result here to obtain external binding data,
			// but nobody loves to mess with txresult mess and it's also safer to check explicitly
			running.UntilSuccess(ctx, s.log, "external-data-getter", func(i context.Context) (bool, error) {
				externalAddr, err = s.horizon.Accounts().CurrentExternalBindingData(s.config.Source.Address(), externalSystem)
				if err != nil {
					return false, errors.Wrap(err, "failed to get external binding data")
				}
				return externalAddr != nil, nil
			}, 5*time.Second, 1*time.Minute)
		}

		// transfer some ETH
		nonce, err := s.eth.NonceAt(ctx, s.config.Keypair.Address(), nil)
		if err != nil {
			return errors.Wrap(err, "failed to get address nonce")
		}

		tx := types.NewTransaction(
			nonce,
			common.HexToAddress(*externalAddr),
			eth.FromGwei(big.NewInt(2000)),
			22000,
			eth.FromGwei(big.NewInt(5)),
			nil,
		)

		tx, err = s.config.Keypair.SignTX(tx)
		if err != nil {
			return errors.Wrap(err, "failed to sign tx")
		}

		if err = s.eth.SendTransaction(ctx, tx); err != nil {
			return errors.Wrap(err, "failed to send transfer tx")
		}

		//
		// at this point we should buksovat, since ETH has been sent
		//

		eth.EnsureHashMined(ctx, s.log, s.eth, tx.Hash())

		minedAt := time.Now()

		// poll account balance for update
		var currentBalance types2.Amount
		running.UntilSuccess(ctx, s.log, "balance-poller", func(i context.Context) (bool, error) {
			currentBalance, err = s.horizon.Accounts().CurrentBalanceIn(s.config.Source.Address(), s.config.Asset)
			if err != nil {
				return false, errors.Wrap(err, "failed to get account balance")
			}
			if currentBalance != balanceBefore {
				return true, nil
			}
			return false, nil
		}, 5*time.Second, 5*time.Second)

		issuedAt := time.Now()

		// TODO ensure balance is updated correctly
		// TODO check if external details are valid

		fmt.Printf("deposit took: %s\n", issuedAt.Sub(minedAt).String())

		return nil
	}, 10*time.Second, 10*time.Second, 10*time.Second)
}
