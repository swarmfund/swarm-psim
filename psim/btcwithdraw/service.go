package btcwithdraw

import (
	"context"
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

const (
	requestStatePending int32 = 1
)

type RequestListener interface{
	Requests(result chan<- horizonV2.Request) <-chan error
}

type BTCClient interface {

}

type Service struct {
	log *logan.Entry
	requestListener RequestListener
	horizon   *horizon.Connector
	btcClient BTCClient

	requests chan horizonV2.Request
	requestListenerErrors <-chan error
}

func New(log *logan.Entry, requestListener RequestListener, horizon *horizon.Connector, btc BTCClient) *Service {

	return &Service{
		log: log.WithField("service", conf.ServiceBTCWithdraw),
		requestListener: requestListener,
		horizon:         horizon,
		btcClient:       btc,

		requests: make(chan horizonV2.Request),
	}
}

// Run is a blocking method, it returns closed channel only when it's finishing.
func (s *Service) Run(ctx context.Context) chan error {
	s.requestListenerErrors = s.requestListener.Requests(s.requests)

	app.RunOverIncrementalTimer(ctx, s.log, "request_processor", s.listenAndProcessRequests, 1 * time.Second)

	errs := make(chan error)
	close(errs)
	return errs
}

func (s *Service) listenAndProcessRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case request := <- s.requests:
		err := s.processRequest(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to process Request", logan.F{
				"request_id": request.ID,
			})
		}

		return nil
	case err := <-s.requestListenerErrors:
		return errors.Wrap(err, "RequestListener sent error")
	}
}

func (s *Service) processRequest(ctx context.Context, request horizonV2.Request) error {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// not a withdraw request
		return nil
	}

	if request.State != requestStatePending {
		return nil
	}

	//err := s.horizon.Transaction(&horizon.TransactionBuilder{
	//	Source: s.config.Source,
	//}).Op(&horizon.ReviewRequestOp{
	//	ID:     request.ID,
	//	Hash:   request.Hash,
	//	Action: xdr.ReviewRequestOpActionApprove,
	//	Details: horizon.ReviewRequestOpDetails{
	//		Type: xdr.ReviewableRequestTypeWithdraw,
	//		Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
	//			// TODO Pass raw TX and hash
	//			ExternalDetails: fmt.Sprintf("%x", txraw),
	//		},
	//	},
	//}).
	//	Sign(s.config.Signer).
	//	Submit()
	//if err != nil {
	//	sErr, ok := errors.Cause(err).(horizon.SubmitError)
	//	if ok {
	//		fmt.Println(string(sErr.ResponseBody()))
	//	}
	//	s.log.WithError(err).Error("Failed to submit ReviewRequest Operation.")
	//}
	//
	//if err := s.btcClient.SendTransaction(ctx, tx); err != nil {
	//	s.log.WithError(err).Error("Failed to send withdraw TX into Bitcoin blockchain.")
	//	return nil
	//}

	return nil
}
