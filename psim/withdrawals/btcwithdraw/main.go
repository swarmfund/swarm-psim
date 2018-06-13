package btcwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdraw"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
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

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)
	btcHelper, err := NewBTCHelper(
		log,
		config.MinWithdrawAmount,
		config.HotWalletAddress,
		config.HotWalletScriptPubKey,
		config.HotWalletRedeemScript,
		config.PrivateKey,
		config.OffchainCurrency,
		config.OffchainBlockchain,
		globalConfig.Bitcoin(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create CommonBTCHelper")
	}

	// FIXME
	// FIXME
	// FIXME
	kp, err := keypair.ParseAddress("GDF6CDA63MU2IW6CQJPNOYEHQBHFF2XNHAPLR6ZUOJP3SBQRKROZFO7Z")
	if err != nil {
		panic(errors.Wrap(err, "failed to parse kp"))
	}

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
		btcHelper,
		// FIXME
		// FIXME
		// FIXME
		true,
		kp,
	), nil
}
