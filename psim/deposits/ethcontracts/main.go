package ethcontracts

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHContracts, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceETHContracts)).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceETHContracts,
		})
	}

	ethWallet := eth.NewWallet()
	ethAddress, err := ethWallet.ImportHEX(config.ETHPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to import PrivKey hex into ETH Wallet")
	}

	// TODO
	//horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	//horizonInfo, err := horizonConnector.Info()
	//if err != nil {
	//	return nil, errors.Wrap(err, "Failed to get Horizon info")
	//}

	//builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return NewService(
		log,
		config,
		ethAddress,
		globalConfig.Ethereum(),
		ethWallet,
	), nil
}
