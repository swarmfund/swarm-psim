package btcwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdraw"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdraw, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(app.Config(ctx).GetRequired(conf.ServiceBTCWithdraw))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
			"service": conf.ServiceBTCWithdraw,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.SignerKP)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return withdraw.New(
		conf.ServiceBTCWithdraw,
		conf.ServiceBTCWithdrawVerify,
		config.SignerKP,
		log,
		horizonConnector.Listener(),
		horizonConnector.Operations(),
		horizonConnector.Submitter(),
		builder,
		globalConfig.Discovery(),
		NewBTCHelper(
			log,
			config.MinWithdrawAmount,
			config.HotWalletAddress,
			config.HotWalletScriptPubKey,
			config.HotWalletRedeemScript,
			config.PrivateKey,
			globalConfig.Bitcoin(),
		),
	), nil
}
