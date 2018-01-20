package btcwithdraw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/conf"
)

const (
	RequestStatePending int32 = 1
	BTCAsset                  = "BTC"
)

var (
	ErrMissingAddress    = errors.New("Missing field in the ExternalDetails json of WithdrawalRequest.")
	ErrAddressNotAString = errors.New("Address field in ExternalDetails of WithdrawalRequest is not a string.")

	ErrNoVerifyServices    = errors.New("No BTC Withdraw Verify services were found.")
	ErrBadStatusFromVerify = errors.New("Unsuccessful status code from Verify.")
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
	WithdrawalRequests(result chan<- horizon.Request) <-chan error
}

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	CreateAndFundRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error)
	GetNetParams() *chaincfg.Params
}

type Service struct {
	log             *logan.Entry
	config          Config
	requestListener RequestListener
	horizon         *horizon.Connector
	btcClient       BTCClient
	discovery       *discovery.Client

	requests              chan horizon.Request
	requestListenerErrors <-chan error
}

// New is constructor for btcwithdraw Service.
func New(log *logan.Entry, config Config,
	requestListener RequestListener, horizonConnector *horizon.Connector, btc BTCClient, discoveryClient *discovery.Client) *Service {

	return &Service{
		log:             log.WithField("service", conf.ServiceBTCWithdraw),
		config:          config,
		requestListener: requestListener,
		horizon:         horizonConnector,
		btcClient:       btc,
		discovery:       discoveryClient,

		requests: make(chan horizon.Request),
	}
}

// Run is a blocking method, it returns closed channel only when it has finishing job.
func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	s.requestListenerErrors = s.requestListener.WithdrawalRequests(s.requests)

	app.RunOverIncrementalTimer(ctx, s.log, "request_processor", s.listenAndProcessRequests, 0, 5*time.Second)
}

func (s *Service) listenAndProcessRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case request := <-s.requests:
		err := s.processRequest(request)
		if err != nil {
			return errors.Wrap(err, "Failed to process Withdraw Request", GetRequestLoganFields("request", request))
		}

		return nil
	case err := <-s.requestListenerErrors:
		return errors.Wrap(err, "RequestListener sent error")
	}
}

func (s *Service) processRequest(request horizon.Request) error {
	if !IsPendingBTCWithdraw(request) {
		return nil
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).Debug("Found pending BTC Withdrawal Request.")

	withdrawAddress, err := GetWithdrawAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get BTC Address from the WithdrawalRequest")
	}
	// Divide by precision of the system.
	withdrawAmount := GetWithdrawAmount(request)
	fields := logan.F{
		"withdraw_address": withdrawAddress,
		"withdraw_amount":  withdrawAmount,
	}

	// Validate
	rejectReason := s.getRejectReason(withdrawAddress, withdrawAmount)
	if rejectReason != "" {
		s.log.WithFields(fields).WithFields(GetRequestLoganFields("request", request)).WithField("reject_reason", rejectReason).
			Warn("Got BTC Withdraw Request which is invalid due to some RejectReason.")

		err = s.sendRejectToVerify(request, rejectReason)
		if err != nil {
			fields["reject_reason"] = rejectReason
			return errors.Wrap(err, "Failed to submit Reject of Request to Verify", fields)
		}

		// Request is invalid, Reject was submitted successfully.
		return nil
	}

	err = s.processValidPendingWithdraw(withdrawAddress, withdrawAmount, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process valid pending WithdrawalRequest", fields)
	}

	return nil
}

func (s *Service) getRejectReason(withdrawAddress string, amount float64) RejectReason {
	err := ValidateBTCAddress(withdrawAddress, s.btcClient.GetNetParams())
	if err != nil {
		return RejectReasonInvalidAddress
	}

	if amount < s.config.MinWithdrawAmount {
		return RejectReasonTooLittleAmount
	}

	return ""
}

func (s *Service) processValidPendingWithdraw(withdrawAddress string, withdrawAmount float64, request horizon.Request) error {
	fields := GetRequestLoganFields("request", request).Merge(logan.F{
		"withdraw_address": withdrawAddress,
		"withdraw_amount":  withdrawAmount,
	})

	signedTXHex, err := s.prepareSignedBitcoinTX(withdrawAddress, withdrawAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare signed Bitcoin TX", fields)
	}

	fields["signed_tx_hex"] = signedTXHex

	err = s.sendApproveToVerify(request, signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to send pre-signed ReviewRequestOp to Verify", fields)
	}

	return nil
}

func (s *Service) prepareSignedBitcoinTX(withdrawAddress string, withdrawAmount float64) (signedTXHex string, err error) {
	unsignedTXHex, err := s.btcClient.CreateAndFundRawTX(withdrawAddress, withdrawAmount, s.config.HotWalletAddress)
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

func (s *Service) sendRejectToVerify(request horizon.Request, reason RejectReason) error {
	envelope, err := s.prepareAndSignTX(xdrbuild.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Reason: string(reason),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to prepare ReviewRequest Reject Transaction")
	}

	if err := s.sendTXToVerify(envelope); err != nil {
		return errors.Wrap(err, "Failed to send TX to Verify")
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).WithField("reject_reason", reason).
		Info("Sent PermanentReject to Verify successfully.")
	return nil
}

func (s *Service) sendApproveToVerify(request horizon.Request, signedTXHex string) error {
	externalDetails := ExternalDetails{
		TXHex: signedTXHex,
	}

	detailsBytes, err := json.Marshal(externalDetails)
	if err != nil {
		errors.Wrap(err, "Failed to marshal ExternalDetails for OpWithdrawal (containing hex and hash of BTC TX)")
	}

	envelope, err := s.prepareAndSignTX(xdrbuild.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.WithdrawalDetails{
			ExternalDetails: string(detailsBytes),
		},
	})
	if err != nil {
		return errors.Wrap(err, "Failed to prepare ReviewRequest Approve Transaction")
	}

	err = s.sendTXToVerify(envelope)
	if err != nil {
		return errors.Wrap(err, "Failed to send TX to Verify")
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).WithField("signed_tx_hex", signedTXHex).
		Info("Sent Approve to Verify successfully.")

	return nil
}

func (s *Service) prepareAndSignTX(op xdrbuild.ReviewRequestOp) (envelope string, err error) {
	info, err := s.horizon.Info()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Horizon info")
	}
	builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

	return builder.Transaction(s.config.SourceKP).
		Op(op).
		Sign(s.config.SignerKP).
		Marshal()
}

func (s *Service) sendTXToVerify(envelope string) error {
	body := ReviewRequest{
		Envelope: envelope,
	}

	rawRequestBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal ReviewRequest (with Envelope)")
	}

	// Find Verify
	services, err := s.discovery.DiscoverService(conf.ServiceBTCWithdrawVerify)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", conf.ServiceBTCWithdrawVerify))
	}
	if len(services) == 0 {
		return ErrNoVerifyServices
	}

	url := services[0].Address
	fields := logan.F{
		"verify_url":       url,
		"raw_request_body": string(rawRequestBody),
	}

	bodyReader := bytes.NewReader(rawRequestBody)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return errors.Wrap(err, "Failed to create Review Request to Verify", fields)
	}

	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to send the request", fields)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fields := logan.F{
			"verify_url":       url,
			"status_code":      response.StatusCode,
			"raw_request_body": string(rawRequestBody),
		}

		// TODO
		//defer func() { _ = resp.Body.Close() }()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "Failed to read the body of response from Verify", fields)
		}

		return errors.From(ErrBadStatusFromVerify, fields.Merge(logan.F{"response_body": string(responseBody)}))
	}

	return nil
}
