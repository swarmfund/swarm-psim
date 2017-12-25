package ethfunnel

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHFunnel, func(ctx context.Context) (utils.Service, error) {
		config := Config{}
		err := figure.
			Out(&config).
			From(app.Config(ctx).GetRequired(conf.ServiceETHFunnel)).
			With(figure.BaseHooks, utils.ETHHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		keystore := keystore.NewKeyStore(
			config.Keystore, keystore.LightScryptN, keystore.LightScryptP,
		)
		for _, account := range keystore.Accounts() {
			err := keystore.Unlock(account, config.Passphrase)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unlock")
			}
		}

		eth := app.Config(ctx).Ethereum()

		return NewService(ctx, app.Log(ctx), config, keystore, eth), nil
	})
}
