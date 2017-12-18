package ethwithdraw

import (
	"context"

	"math/big"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type Config struct {
	Keystore string
	From     string
	GasPrice *big.Int
	Source   keypair.KP
	Signer   keypair.KP
}

func init() {
	app.RegisterService(conf.ServiceETHWithdraw, func(ctx context.Context) (utils.Service, error) {
		config := Config{
			GasPrice: big.NewInt(1000000000),
		}
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.CommonHooks).
			From(app.Config(ctx).Get(conf.ServiceETHWithdraw)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		ks := keystore.NewKeyStore(config.Keystore, keystore.LightScryptN, keystore.LightScryptP)
		if !ks.HasAddress(common.HexToAddress(config.From)) {
			return nil, errors.New("now wallet for address")
		}

		for _, account := range ks.Accounts() {
			ks.Unlock(account, "foobar")
		}

		horizonV2 := app.Config(ctx).HorizonV2()
		horizon, err := app.Config(ctx).Horizon()
		if err != nil {
			return nil, errors.New("failed to get horizon")
		}

		eth := app.Config(ctx).Ethereum()

		return NewService(ctx, app.Log(ctx), config, horizonV2, ks, eth, horizon), nil
	})
}

type Service struct {
	ctx       context.Context
	log       *logan.Entry
	horizonV2 *horizonV2.Connector
	horizon   *horizon.Connector
	config    Config
	keystore  *keystore.KeyStore
	eth       *ethclient.Client
}

func NewService(
	ctx context.Context, log *logan.Entry, config Config, horizonV2 *horizonV2.Connector,
	keystore *keystore.KeyStore, eth *ethclient.Client, horizon *horizon.Connector,
) *Service {
	return &Service{
		ctx:       ctx,
		log:       log,
		config:    config,
		horizonV2: horizonV2,
		horizon:   horizon,
		keystore:  keystore,
		eth:       eth,
	}
}

func (s *Service) Run() chan error {
	requestCh := make(chan horizonV2.Request)
	go func() {
		errs := s.horizonV2.Listener().Requests(requestCh)
		for {
			select {
			case request := <-requestCh:
				fmt.Println("got request")
				fmt.Println(request.Details.RequestType)
				if request.Details.RequestType != 4 {
					// not a withdraw request
					continue
				}
				fmt.Println("it's withdraw!")
				if request.State != 1 {
					// not pending
					continue
				}
				// TODO check asset
				fmt.Println("and it's pending")
				fmt.Println("amount", int64(request.Details.Withdraw.Amount))
				nonce, err := s.eth.NonceAt(s.ctx, common.HexToAddress(s.config.From), nil)
				if err != nil {
					s.log.WithError(err).Error("failed to get nonce")
				}

				txgas := big.NewInt(21000)
				value := new(big.Int).
					Sub(
						new(big.Int).Mul(
							big.NewInt(int64(request.Details.Withdraw.Amount)),
							new(big.Int).Mul(
								big.NewInt(10000000),
								big.NewInt(10000000),
							),
						),
						new(big.Int).Mul(txgas, s.config.GasPrice),
					)
				fmt.Println(value)
				tx, err := s.keystore.SignTx(
					accounts.Account{
						Address: common.HexToAddress(s.config.From),
					},
					types.NewTransaction(
						nonce,
						common.HexToAddress(request.Details.Withdraw.ExternalDetails),
						value,
						txgas,
						s.config.GasPrice,
						nil,
					),
					nil,
				)
				if err != nil {
					s.log.WithError(err).Error("failed to sign tx")
					continue
				}
				txraw, err := rlp.EncodeToBytes(tx.Data())
				if err != nil {
					s.log.WithError(err).Error("failed to encode tx")
					continue
				}
				err = s.horizon.Transaction(&horizon.TransactionBuilder{
					Source: s.config.Source,
				}).Op(&horizon.ReviewRequestOp{
					ID:     request.ID,
					Hash:   request.Hash,
					Action: xdr.ReviewRequestOpActionApprove,
					Details: horizon.ReviewRequestOpDetails{
						Type: xdr.ReviewableRequestTypeWithdraw,
						Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
							ExternalDetails: fmt.Sprintf("%x", txraw),
						},
					},
				}).
					Sign(s.config.Signer).
					Submit()
				if err != nil {
					serr, ok := errors.Cause(err).(horizon.SubmitError)
					if ok {
						fmt.Println(string(serr.ResponseBody()))
					}
					s.log.WithError(err).Error("failed to submit review tx")
				}
				if err := s.eth.SendTransaction(s.ctx, tx); err != nil {
					s.log.WithError(err).Error("failed to send tx")
					continue
				}
				fmt.Println(tx.String())
			case err := <-errs:
				s.log.WithError(err).Warn("failed to get request")
			}
		}
	}()

	return make(chan error)
}
