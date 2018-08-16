package btcdepositveri

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/btcdeposit"
	"gitlab.com/swarmfund/psim/psim/deposits/depositveri"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceBTCDepositVerify, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(globalConfig.GetRequired(conf.ServiceBTCDepositVerify))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
			"service": conf.ServiceBTCDepositVerify,
		})
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	btcHelper, err := btcdeposit.NewBTCHelper(
		log,
		config.DepositAsset,
		config.MinDepositAmount,
		config.FixedDepositFee,
		config.OffchainCurrency,
		config.OffchainBlockchain,
		config.BlocksToSearchForTX,
		globalConfig.Bitcoin(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create BTCHelper")
	}

	return depositveri.New(
		config.ExternalSystem,
		conf.ServiceBTCDepositVerify,
		log,
		config.Signer,
		config.LastBlocksNotWatch,
		horizonConnector,
		xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		listener,
		globalConfig.Discovery(),
		btcHelper,
	), nil
}
