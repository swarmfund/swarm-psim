package btcfunnel

import (
	"context"
	"fmt"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/distributed_lab/logan/v3"
)

func init() {
	app.RegisterService(conf.ServiceBTCFunnel, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx).WithField("service", conf.ServiceBTCFunnel)

	config := Config{}

	err := figure.
		Out(&config).
		From(globalConfig.GetRequired(conf.ServiceBTCFunnel)).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBTCFunnel))
	}

	if config.BlocksToBeIncluded < 2 || config.BlocksToBeIncluded  > 25 {
		return nil, errors.Errorf("Invalid BocksToBeIncluded value (%d), must be from 2 to 25.", config.BlocksToBeIncluded)
	}
	if config.MaxFeePerKB <= 0 {
		return nil, errors.Errorf("Invalid MaxFeePerKB value (%.8f), must be grater than zero.", config.MaxFeePerKB)
	}

	netParams, err := bitcoin.GetNetParams(config.OffchainCurrency, config.OffchainBlockchain)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build NetParams by currency and blockchain", logan.F{
			"currency":   config.OffchainCurrency,
			"blockchain": config.OffchainBlockchain,
		})
	}

	return New(config, log, globalConfig.Bitcoin(), netParams, globalConfig.NotificationSender()), nil
}
