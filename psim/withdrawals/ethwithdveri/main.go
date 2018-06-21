package ethwithdveri

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHWithdrawVerify, func(ctx context.Context) (app.Service, error) {
		fields := logan.F{
			"service": conf.ServiceETHWithdrawVerify,
		}

		var config Config
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).GetRequired(conf.ServiceETHWithdrawVerify)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to figure out config", fields)
		}

		if err := config.Validate(); err != nil {
			return nil, errors.Wrap(err, "Config is invalid", fields)
		}

		log := app.Log(ctx)

		horizonConnector := app.Config(ctx).Horizon().WithSigner(config.Signer)

		info, err := horizonConnector.System().Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}
		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

		wallet := eth.NewWallet()
		ethAddress, err := wallet.ImportHEX(config.PrivateKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import key")
		}

		service, err := NewService(
			log,
			config,
			ethAddress,
			horizonConnector.Listener(),
			horizonConnector.Operations(),
			builder,
			horizonConnector.Submitter(),
			app.Config(ctx).Ethereum(),
			wallet,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create Service")
		}

		return service, nil
	})
}
