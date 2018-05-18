package ethfunnel

import (
	"context"

	"math/big"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHFunnel, func(ctx context.Context) (app.Service, error) {
		config := Config{
			Confirmations: big.NewInt(12),
		}
		err := figure.
			Out(&config).
			From(app.Config(ctx).GetRequired(conf.ServiceETHFunnel)).
			With(figure.BaseHooks, utils.ETHHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		wallet, err := eth.NewHDWallet(config.Seed, config.AccountsToDerive)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init wallet")
		}

		return NewService(app.Log(ctx), config, wallet, app.Config(ctx).Ethereum()), nil
	})
}
