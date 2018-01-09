package btcwithdraw

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/conf"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	RequestStatePending int32 = 1
	BTCAsset                  = "BTC"
)

var (
	ErrBadStatusFromVerify = errors.New("Unsuccessful status code from Verify.")
	ErrMissingAddress      = errors.New("Missing field in the ExternalDetails json of WithdrawalRequest.")
	ErrAddressNotAString   = errors.New("Address field in ExternalDetails of WithdrawalRequest is not a string.")
)

// ExternalDetails is used to marshal and unmarshal external
// details of Withdrawal Details for ReviewRequest Operation
// during approve.
type ExternalDetails struct {
	TXHash string `json:"tx_hash"`
	TXHex  string `json:"tx_hex"`
}

// ReviewRequest is the data structure to send pre-signed Request
// to Verify (Service btcwithdveri)
type ReviewRequest struct {
	Envelope string `json:"envelope"`
}

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	WithdrawalRequests(result chan<- horizonV2.Request) <-chan error
}

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	CreateRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error)
}

type Service struct {
	log             *logan.Entry
	config          Config
	requestListener RequestListener
	horizon         *horizon.Connector
	btcClient       BTCClient
	discovery       *discovery.Client

	requests              chan horizonV2.Request
	requestListenerErrors <-chan error
}

// New is constructor for btcwithdraw Service.
func New(log *logan.Entry, config Config,
	requestListener RequestListener, horizon *horizon.Connector, btc BTCClient, discoveryClient *discovery.Client) *Service {

	return &Service{
		log:             log.WithField("service", conf.ServiceBTCWithdraw),
		config:          config,
		requestListener: requestListener,
		horizon:         horizon,
		btcClient:       btc,
		discovery:       discoveryClient,

		requests: make(chan horizonV2.Request),
	}
}

// Run is a blocking method, it returns closed channel only when it has finishing job.
func (s *Service) Run(ctx context.Context) chan error {
	s.log.Info("Starting.")

	s.requestListenerErrors = s.requestListener.WithdrawalRequests(s.requests)

	app.RunOverIncrementalTimer(ctx, s.log, "request_processor", s.listenAndProcessRequests, 0, 5*time.Second)

	errs := make(chan error)
	close(errs)
	return errs
}

func (s *Service) listenAndProcessRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case request := <-s.requests:
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

	if request.State != RequestStatePending {
		return nil
	}

	if request.Details.Withdraw.DestinationAsset != BTCAsset {
		// Withdraw not to a BTC - not interesting for this Service.
		return nil
	}

	s.log.WithFields(getRequestLoganFields("request", request)).Debug("Found pending BTC Withdrawal Request.")

	withdrawAddress, err := ObtainAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain BTC Address from the WithdrawalRequest.")
	}
	// Divide by precision of the system.
	withdrawAmount := float64(int64(request.Details.Withdraw.DestinationAmount)) / amount.One

	// Validate
	isValid, err := s.validateOrReject(withdrawAddress, withdrawAmount, request.ID, request.Hash)
	if err != nil {
		return errors.Wrap(err, "Failed to validateOrReject Request")
	}
	if !isValid {
		// Request is invalid, but PermanentReject was successfully submitted.
		return nil
	}

	err = s.processValidPendingWithdraw(withdrawAddress, withdrawAmount, request.ID, request.Hash)
	if err != nil {
		return err
	}

	return nil
}

// TODO Consider moving to so common, as this logic is common for BTC and ETH.
func ObtainAddress(request horizonV2.Request) (string, error) {
	addrValue, ok := request.Details.Withdraw.ExternalDetails["address"]
	if !ok {
		return "", ErrMissingAddress
	}

	addr, ok := addrValue.(string)
	if !ok {
		return "", errors.From(ErrAddressNotAString, logan.F{"raw_address_value": addrValue})
	}

	return addr, nil
}

func (s *Service) validateOrReject(withdrawAddress string, amount float64, requestID uint64, requestHash string) (isValid bool, err error) {
	rejectReason := s.getRejectReason(withdrawAddress, amount, requestID)

	if rejectReason == "" {
		return true, nil
	}

	err = s.sendRejectToVerify(requestID, requestHash, rejectReason)
	if err != nil {
		return false, errors.Wrap(err, "Failed to submit Reject of Request to Verify",
			logan.F{
				"withdraw_address": withdrawAddress,
				"reject_reason":    rejectReason,
			})
	}

	// Request is invalid, Reject was submitted successfully.
	return false, nil
}

func (s *Service) getRejectReason(withdrawAddress string, amount float64, requestID uint64) RejectReason {
	_, err := btc.NewAddrFromString(withdrawAddress)
	if err != nil {
		s.log.WithField("withdraw_address", withdrawAddress).WithField("amount", amount).WithField("request_id", requestID).WithError(err).
			Warn("Got BTC Withdraw Request with wrong BTC Address.")
		return RejectReasonInvalidAddress
	}

	if amount < s.config.MinWithdrawAmount {
		s.log.WithField("withdraw_address", withdrawAddress).WithField("amount", amount).WithField("request_id", requestID).
			Warn("Got BTC Withdraw Request with too little amount.")
		return RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) processValidPendingWithdraw(withdrawAddress string, withdrawAmount float64,
	requestID uint64, requestHash string) error {

	fields := logan.F{
		"request_id":       requestID,
		"withdraw_address": withdrawAddress,
		"withdraw_amount":  withdrawAmount,
	}

	s.log.WithFields(fields).Info("Processing valid pending Withdraw Request.")

	signedTXHex, err := s.prepareSignedBitcoinTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare signed Bitcoin TX", fields)
	}

	fields["signed_tx_hex"] = signedTXHex

	err = s.sendApproveToVerify(requestID, requestHash, signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to send pre-signed ReviewRequestOp to Verify", fields)
	}

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

	signedTXHex, err = s.btcClient.SignAllTXInputs(unsignedTXHex, s.config.HotWalletScriptPubKey, &s.config.HotWalletRedeemScript, s.config.PrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sing raw TX", logan.F{"unsigned_tx_hex": unsignedTXHex})
	}

	return signedTXHex, nil
}

func (s *Service) sendRejectToVerify(requestID uint64, requestHash string, reason RejectReason) error {
	tx := s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.SourceKP,
	}).Op(&horizon.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Reason: string(reason),
		Details: horizon.ReviewRequestOpDetails{
			Type:       xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{},
		},
	}).
		Sign(s.config.SignerKP)

	err := s.sendTXToVerify(tx)
	if err != nil {
		return errors.Wrap(err, "Failed to send TX to Verify")
	}

	return nil
}

func (s *Service) sendApproveToVerify(requestID uint64, signedTXHash, signedTXHex string) error {
	externalDetails := ExternalDetails{
		TXHash: signedTXHash,
		TXHex:  signedTXHex,
	}

	detailsBytes, err := json.Marshal(externalDetails)
	if err != nil {
		errors.Wrap(err, "Failed to marshal ExternalDetails for OpWithdrawal (containing hex and hash of BTC TX)")
	}

	tx := s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.SourceKP,
	}).Op(&horizon.ReviewRequestOp{
		ID:     requestID,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: horizon.ReviewRequestOpDetails{
			Type: xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
				ExternalDetails: string(detailsBytes),
			},
		},
	}).
		Sign(s.config.SignerKP)

	err = s.sendTXToVerify(tx)
	if err != nil {
		return errors.Wrap(err, "Failed to send TX to Verify")
	}

	return nil
}

func (s *Service) sendTXToVerify(horizonTX *horizon.TransactionBuilder) error {
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	url := "http://localhost:8101/"

	// FIXME
	//services, err := s.discovery.DiscoverService(conf.ServiceBTCWithdrawVerify)
	//if err != nil {
	//	return errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", conf.ServiceBTCWithdrawVerify))
	//}
	//if len(services) == 0 {
	//	// TODO
	//}

	xdrTX, err := horizonTX.Marshal64()
	if err != nil {
		return errors.Wrap(err, "Failed to Marshal64 the horizonTX")
	}
	if xdrTX == nil {
		return errors.Wrap(err, "Marshal64 returned nil value without an error")
	}
	body := ReviewRequest{
		Envelope: *xdrTX,
	}

	rawRequestBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal ReviewRequest (with Envelope)")
	}

	bodyReader := bytes.NewReader(rawRequestBody)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return errors.Wrap(err, "Failed to create Review Request to Verify", logan.F{
			"url":              url,
			"raw_request_body": string(rawRequestBody),
		})
	}

	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to send the request", logan.F{
			"url":              url,
			"raw_request_body": string(rawRequestBody),
		})
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		//defer func() { _ = resp.Body.Close() }()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			// TODO
		}

		return errors.From(ErrBadStatusFromVerify, logan.F{
			"status_code":      response.StatusCode,
			"raw_request_body": string(rawRequestBody),
			"response_body":    string(responseBody),
		})
	}

	return nil
}
