package contractfunnel

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHContractFunnel, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceETHContractFunnel)).
		With(figure.BaseHooks, utils.ETHHooks, hooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	//ethWallet := eth.NewWallet()
	//ethAddress, err := ethWallet.ImportHEX(config.ETHPrivateKey)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Failed to import PrivKey hex into ETH Wallet")
	//}

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
		globalConfig.Ethereum(),
		//ethWallet,
	)
}
