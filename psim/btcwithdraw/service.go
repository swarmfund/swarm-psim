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
	"fmt"
)

const (
	requestStatePending int32 = 1
)

type RequestListener interface{
	Requests(result chan<- horizonV2.Request) <-chan error
}

type BTCClient interface {
	CreateRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignRawTX(txHex, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
}

type Service struct {
	log *logan.Entry
	config Config
	requestListener RequestListener
	horizon   *horizon.Connector
	btcClient BTCClient

	requests chan horizonV2.Request
	requestListenerErrors <-chan error
}

func New(log *logan.Entry, config Config, requestListener RequestListener, horizon *horizon.Connector, btc BTCClient) *Service {

	return &Service{
		log: log.WithField("service", conf.ServiceBTCWithdraw),
		config:          config,
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

// TODO Refactor me.
func (s *Service) processRequest(ctx context.Context, request horizonV2.Request) error {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// not a withdraw request
		return nil
	}

	if request.State != requestStatePending {
		return nil
	}

	// Divide by 10^8 (satoshi)
	amount := float64(int64(request.Details.Withdraw.Amount)) / 10000000.0
	withdrawAddress := string(request.Details.Withdraw.ExternalDetails)
	fields := logan.F{
		"withdraw_address": withdrawAddress,
		"withdraw_amount":           amount,
	}

	txHex, err := s.btcClient.CreateRawTX(withdrawAddress, amount, s.config.HotWalletAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to create raw TX", fields)
	}

	txHex, err = s.btcClient.SignRawTX(txHex, s.config.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to sing raw TX using first PrivateKey", fields)
	}
	// TODO Move signing by second PrivateKey to some verifier service.
	txHex, err = s.btcClient.SignRawTX(txHex, s.config.PrivateKey2)
	if err != nil {
		return errors.Wrap(err, "Failed to sing raw TX using second PrivateKey", fields)
	}

	// Now txHex is hex of ready and signed TX.
	fields = fields.Add("signed_tx_hex", txHex)

	err = s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.SourceKP,
	}).Op(&horizon.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: horizon.ReviewRequestOpDetails{
			Type: xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
				// TODO Pass raw TX and hash in some JSON.
				ExternalDetails: fmt.Sprintf("%x", txHex),
			},
		},
	}).
		Sign(s.config.SignerKP).
		Submit()
	if err != nil {
		sErr, ok := errors.Cause(err).(horizon.SubmitError)
		if ok {
			fields = fields.Add("horizon_submit_error_response_body", string(sErr.ResponseBody()))
		}
		s.log.WithFields(fields).WithError(err).Error("Failed to submit ReviewRequest Operation.")

		return errors.Wrap(err, "Failed to submit Transaction to Horizon")
	}

	txHash, err := s.btcClient.SendRawTX(txHex)
	if err != nil {
		// This problem should be fixed manually.
		// Transactions from approved requests not existing in the Bitcoin blockchain
		// should be submitted once more.
		// This process should probably be automated.
		s.log.WithFields(fields).WithError(err).Error("Failed to send withdraw TX into Bitcoin blockchain.")
		return nil
	}

	s.log.WithField("tx_hash", txHash).Info("Sent withdraw TX to Bitcoin blockchain successfully.")
	return nil
}
