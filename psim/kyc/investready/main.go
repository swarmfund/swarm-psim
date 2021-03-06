// Sequence diagram for this service can be found here:
// https://gitlab.com/tokend/knowledge_base/blob/master/flows/kyc/invest-ready-sequence-diagram.jpeg
package investready

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceInvestReady, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceInvestReady)).
		With(figure.BaseHooks, utils.ETHHooks, hooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceInvestReady,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return NewService(
		log,
		config,
		horizonConnector.Listener(),
		horizonConnector.Operations(),
		horizonConnector.Accounts(),
		kyc.NewRequestPerformer(builder, config.Source, config.Signer, horizonConnector.Submitter()),
		horizonConnector.Blobs(),
		horizonConnector.Users(),
		NewConnector(log.WithField("service", conf.ServiceInvestReady), config.Connector),
		doorman.New(!config.RedirectsConfig.CheckSignature, horizonConnector.Accounts()),
	), nil
}
