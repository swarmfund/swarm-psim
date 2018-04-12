package notifier

import (
	"context"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceOperationNotifier, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(globalConfig.GetRequired(conf.ServiceOperationNotifier)).
		With(figure.BaseHooks, utils.ETHHooks, EventHooks, EmailsHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out", logan.F{
			"service": conf.ServiceOperationNotifier,
		})
	}

	err = checkRequestTokenSuffixesValidity(config)
	if err != nil {
		return nil, errors.Wrap(err, "invalid 'email_request_token_suffix'", logan.F{
			"service": conf.ServiceOperationNotifier,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	checkSaleStateResponses := horizonConnector.Listener().StreamAllCheckSaleStateOps(ctx, 0)
	createKYCRequestOpResponses := horizonConnector.Listener().StreamAllCreateKYCRequestOps(ctx, 0)
	reviewRequestOpResponses := horizonConnector.Listener().StreamAllReviewRequestOps(ctx, 0)

	service, err := New(
		config,
		log,
		globalConfig.Notificator(),
		horizonConnector.Operations(),
		horizonConnector.Sales(),
		horizonConnector.Templates(),
		horizonConnector.Transactions(),
		horizonConnector.Users(),
		checkSaleStateResponses,
		createKYCRequestOpResponses,
		reviewRequestOpResponses,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create service", logan.F{
			"service": conf.ServiceOperationNotifier,
		})
	}

	return service, nil
}

func checkRequestTokenSuffixesValidity(config Config) error {
	if len(config.OrderCancelled.Emails.RequestTokenSuffix) == 0 {
		return errors.New("'email_request_token_suffix' in order_cancelled must not be empty")
	}
	if len(config.KYCCreated.Emails.RequestTokenSuffix) == 0 {
		return errors.New("'email_request_token_suffix' in kyc_created must not be empty")
	}
	if len(config.KYCApproved.Emails.RequestTokenSuffix) == 0 {
		return errors.New("'email_request_token_suffix' in kyc_approved must not be empty")
	}
	if len(config.KYCRejected.Emails.RequestTokenSuffix) == 0 {
		return errors.New("'email_request_token_suffix' in kyc_rejected must not be empty")
	}

	return nil
}
