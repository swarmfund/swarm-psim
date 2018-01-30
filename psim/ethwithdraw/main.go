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
	txGas = big.NewInt(21000)
	// DEPRECATED Now uses 12. Move to amount.One
	ethPrecision = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1000000))
)

func init() {
	app.RegisterService(conf.ServiceETHWithdraw, func(ctx context.Context) (app.Service, error) {
		config := Config{}
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).GetRequired(conf.ServiceETHWithdraw)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		if config.GasPrice == nil {
			return nil, errors.New("'gas_price' cannot be empty")
		}

		wallet := eth.NewWallet()
		address, err := wallet.ImportHEX(config.Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import key")
		}

		return NewService(
			app.Log(ctx),
			config,
			app.Config(ctx).Horizon().WithSigner(config.Signer),
			wallet,
			app.Config(ctx).Ethereum(),
			address), nil
	})
}
