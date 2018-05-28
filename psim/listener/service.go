package listener

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/keypair"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
)

type ServiceConfig struct {
	Signer keypair.Full
}

type Service struct {
	config      ServiceConfig
	extractor   internal.Extractor
	handler     TokendHandler
	broadcaster internal.Broadcaster
	logger      *logan.Entry
}

func NewService(config ServiceConfig, extractor internal.Extractor, handler TokendHandler, broadcaster internal.Broadcaster, log *logan.Entry) *Service {
	return &Service{
		config:      config,
		extractor:   extractor,
		handler:     handler,
		broadcaster: broadcaster,
		logger:      log,
	}
}

func (s *Service) Run(ctx context.Context) {
	running.UntilSuccess(ctx, s.logger, conf.ListenerService, s.DispatchEvents, defaultServiceRetryTimeIncrement, defaultMaxServiceRetryTime)
}

func (s *Service) registerProcessors() {
	th := s.handler
	th.SetProcessor(xdr.OperationTypeCreateKycRequest, th.processKYCCreateUpdateRequestOp)
	th.SetProcessor(xdr.OperationTypeReviewRequest, th.processReviewRequestOp)
	th.SetProcessor(xdr.OperationTypeCreateIssuanceRequest, th.processCreateIssuanceRequestOp)
	th.SetProcessor(xdr.OperationTypeManageOffer, th.processManageOfferOp)
	th.SetProcessor(xdr.OperationTypePayment, th.processPayment)
	th.SetProcessor(xdr.OperationTypePaymentV2, th.processPaymentV2)
	th.SetProcessor(xdr.OperationTypeCreateWithdrawalRequest, th.processWithdrawRequest)
	th.SetProcessor(xdr.OperationTypeCreateAccount, th.processCreateAccountOp)

}

func (s *Service) DispatchEvents(ctx context.Context) (bool, error) {
	extractedTxData, err := s.extractor.Extract(ctx)
	if err != nil {
		return false, err
	}

	s.registerProcessors()
	emittedEvents, err := s.handler.Process(extractedTxData)
	if err != nil {
		return false, err
	}
	err = s.broadcaster.BroadcastEvents(ctx, emittedEvents)
	if err != nil {
		return false, err
	}

	return false, nil
}
