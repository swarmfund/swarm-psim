package ethwithdraw

import (
	"context"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

const (
	requestStatePending int32 = 1
)

var (
	txGas        = big.NewInt(21000)
	ethPrecision = new(big.Int).Mul(big.NewInt(10000000), big.NewInt(10000000))
)

func init() {
	app.RegisterService(conf.ServiceETHWithdraw, func(ctx context.Context) (utils.Service, error) {
		config := Config{
			GasPrice: big.NewInt(1000000000),
		}
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).GetRequired(conf.ServiceETHWithdraw)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		ks := keystore.NewKeyStore(config.Keystore, keystore.LightScryptN, keystore.LightScryptP)
		if !ks.HasAddress(config.From) {
			return nil, errors.New("now wallet for address")
		}

		for _, account := range ks.Accounts() {
			if err := ks.Unlock(account, config.Passphrase); err != nil {
				return nil, errors.Wrap(err, "failed to unlock account")
			}
		}

		horizonV2 := app.Config(ctx).HorizonV2()
		horizon, err := app.Config(ctx).Horizon()
		if err != nil {
			return nil, errors.New("failed to get horizon")
		}

		eth := app.Config(ctx).Ethereum()

		return NewService(app.Log(ctx), config, horizonV2, ks, eth, horizon), nil
	})
}
