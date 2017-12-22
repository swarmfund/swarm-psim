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
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"encoding/json"
	"github.com/piotrnar/gocoin/lib/btc"
	"encoding/hex"
)

const (
	requestStatePending int32 = 1
	btcAsset = "BTC"
)

type RequestListener interface{
	Requests(result chan<- horizonV2.Request) <-chan error
}

type BTCClient interface {
	CreateRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error)
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
			return errors.Wrap(err, "Failed to process Withdraw Request", logan.F{
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

	if request.Details.Withdraw.DestinationAsset != btcAsset {
		// Withdraw not to a BTC - not interesting for this Service.
		return nil
	}

	// Divide by 10^4 (precision of the system)
	amount := float64(int64(request.Details.Withdraw.DestinationAmount)) / 10000.0
	withdrawAddress := string(request.Details.Withdraw.ExternalDetails)

	err := s.processPendingWithdrawRequest(ctx, withdrawAddress, amount, request.ID, request.Hash)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) processPendingWithdrawRequest(ctx context.Context, withdrawAddress string, withdrawAmount float64,
		requestID uint64, requestHash string) error {

	fields := logan.F{
		"request_id":       requestID,
		"withdraw_address": withdrawAddress,
		"withdraw_amount":  withdrawAmount,
	}

	s.log.WithFields(fields).Info("Processing pending Withdraw Request.")

	signedTXHex, err := s.prepareSignedBitcoinTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare signed Bitcoin TX", fields)
	}

	fields = fields.Add("signed_tx_hex", signedTXHex)

	txBytes, err := hex.DecodeString(signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to decode signed TX hex into bytes", fields)
	}
	tx, _ := btc.NewTx(txBytes)
	signedTXHash := tx.Hash.String()
	fields = fields.Add("signed_tx_hash", signedTXHash)

	err = s.submitReviewRequestOp(requestID, requestHash, signedTXHash, signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to submit ReviewRequestOp to Horizon", fields)
	}

	sentTXHash, err := s.btcClient.SendRawTX(signedTXHex)
	if err != nil {
		// This problem should be fixed manually.
		// Transactions from approved requests not existing in the Bitcoin blockchain
		// should be submitted once more.
		// This process should probably be automated.
		s.log.WithFields(fields).WithError(err).Error("Failed to send withdraw TX into Bitcoin blockchain.")
		return nil
	}

	fields = fields.Add("sent_tx_hash", sentTXHash)


	s.log.WithFields(fields).Info("Sent withdraw TX to Bitcoin blockchain successfully.")
	return nil

}

func (s *Service) prepareSignedBitcoinTX(withdrawAddress string, withdrawAmount float64) (signedTXHex string, err error) {
	unsignedTXHex, err := s.btcClient.CreateRawTX(withdrawAddress, withdrawAmount, s.config.HotWalletAddress)
	if err != nil {
		if errors.Cause(err) == bitcoin.ErrInsufficientFunds {
			return "", errors.Wrap(err, "Could not create raw TX - not enough BTC on hot wallet")
		}

		return "", errors.Wrap(err, "Failed to create raw TX")
	}

	signedOnceTXHex, err := s.btcClient.SignAllTXInputs(unsignedTXHex, s.config.HotWalletScriptPubKey, &s.config.HotWalletRedeemScript, s.config.PrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sing raw TX using first PrivateKey", logan.F{"unsigned_tx_hex": unsignedTXHex})
	}

	// TODO Move signing by second PrivateKey to some verifier service.
	signedTXHex, err = s.btcClient.SignAllTXInputs(signedOnceTXHex, s.config.HotWalletScriptPubKey, &s.config.HotWalletRedeemScript, s.config.PrivateKey2)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sing raw TX using second PrivateKey", logan.F{"signed_once_tx_hex": signedOnceTXHex})
	}

	return signedTXHex, nil
}

func (s *Service) submitReviewRequestOp(requestID uint64, requestHash, signedTXHash, signedTXHex string) error {
	externalDetails := struct {
		TXHash string `json:"tx_hash"`
		TXHex  string `json:"tx_hex"`
	}{
		TXHash: signedTXHash,
		TXHex:  signedTXHex,
	}
	detailsBytes, err := json.Marshal(externalDetails)
	if err != nil {
		errors.Wrap(err, "Failed to marshal ExternalDetails for OpWithdrawal (containing hex and hash of BTC TX)")
	}

	err = s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.SourceKP,
	}).Op(&horizon.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: horizon.ReviewRequestOpDetails{
			Type: xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
				ExternalDetails: string(detailsBytes),
			},
		},
	}).
		Sign(s.config.SignerKP).
		Submit()

	if err != nil {
		var fields logan.F

		sErr, ok := errors.Cause(err).(horizon.SubmitError)
		if ok {
			fields = logan.F{"horizon_submit_error_response_body": string(sErr.ResponseBody())}
		}

		return errors.Wrap(err, "Failed to submit Transaction to Horizon", fields)
	}

	return nil
}
