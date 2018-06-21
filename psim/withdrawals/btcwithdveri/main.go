// BTC Withdraw Verify
package btcwithdveri

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/withdrawals/btcwithdraw"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdveri"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdrawVerify, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(app.Config(ctx).GetRequired(conf.ServiceBTCWithdrawVerify))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
			"service": conf.ServiceBTCWithdrawVerify,
		})
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.SignerKP)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	btcHelper, err := btcwithdraw.NewBTCHelper(
		log,
		config.TokendAsset,
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

	return withdveri.New(
		conf.ServiceBTCWithdrawVerify,
		log,
		config.SourceKP,
		config.SignerKP,
		horizonConnector.Operations(),
		xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		listener,
		globalConfig.Discovery(),
		btcHelper,
	), nil
}
