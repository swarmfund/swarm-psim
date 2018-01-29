package ethwithdraw

import (
	"context"

	"math/big"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
)

const (
	requestStatePending int32 = 1
)

var (
	txGas        = big.NewInt(21000)
	// DEPRECATED Now uses 12. Move to amount.One
	ethPrecision = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1000000))
)

func init() {
	app.RegisterService(conf.ServiceETHWithdraw, func(ctx context.Context) (app.Service, error) {
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

		wallet := eth.NewWallet()
		address, err := wallet.ImportHEX(config.Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import key")
		}

		horizon := app.Config(ctx).Horizon()
		eth := app.Config(ctx).Ethereum()

		return NewService(app.Log(ctx), config, horizon, wallet, eth, address), nil
	})
}
