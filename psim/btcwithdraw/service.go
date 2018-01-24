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
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

var (
	ErrNoVerifyServices    = errors.New("No BTC Withdraw Verify services were found.")
	ErrBadStatusFromVerify = errors.New("Unsuccessful status code from Verify.")
)

// ExternalDetails is used to marshal and unmarshal external
// details of Withdrawal Details for ReviewRequest Operation
// during approve.
type ExternalDetails struct {
	//TXHash string `json:"tx_hash"`
	TXHex string `json:"tx_hex"`
}

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	WithdrawalRequests(result chan<- horizon.Request) <-chan error
}

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	// CreateAndFundRawTX sets Change position in Outputs to 1.
	CreateAndFundRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	GetNetParams() *chaincfg.Params
}

type Service struct {
	log             *logan.Entry
	config          Config
	requestListener RequestListener
	horizon         *horizon.Connector
	xdrbuilder      *xdrbuild.Builder
	btcClient       BTCClient
	discovery       *discovery.Client

	requests              chan horizon.Request
	requestListenerErrors <-chan error
}

// New is constructor for btcwithdraw Service.
func New(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	horizonConnector *horizon.Connector,
	builder *xdrbuild.Builder,
	btc BTCClient,
	discoveryClient *discovery.Client,
) *Service {

	return &Service{
		log:             log.WithField("service", conf.ServiceBTCWithdraw),
		config:          config,
		requestListener: requestListener,
		horizon:         horizonConnector,
		xdrbuilder:      builder,
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
		err := s.processRequest(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to process Withdraw Request", withdraw.GetRequestLoganFields("request", request))
		}

		return nil
	case err := <-s.requestListenerErrors:
		return errors.Wrap(err, "RequestListener sent error")
	}
}

// FIXME
func (s *Service) processRequest(ctx context.Context, request horizon.Request) error {
	// FIXME Use new RequestType (pending)
	if withdraw.ProvePendingRequest(request, int32(xdr.ReviewableRequestTypeWithdraw), withdraw.BTCAsset) != "" {
		return nil
	}

	s.log.WithFields(withdraw.GetRequestLoganFields("request", request)).Debug("Found pending BTC Withdrawal Request.")

	withdrawAddress, err := withdraw.GetWithdrawAddress(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get BTC Address from the WithdrawalRequest")
	}
	// Divide by precision of the system.
	withdrawAmount := withdraw.GetWithdrawAmount(request)
	fields := logan.F{
		"withdraw_address": withdrawAddress,
		"withdraw_amount":  withdrawAmount,
	}

	rejectReason := s.getRejectReason(withdrawAddress, withdrawAmount)
	if rejectReason != "" {
		fields["reject_reason"] = rejectReason

		s.log.WithFields(fields).WithFields(withdraw.GetRequestLoganFields("request", request)).
			Warn("Got BTC Withdraw Request which is invalid due to the RejectReason.")

		err = s.processRequestReject(ctx, request, rejectReason)
		if err != nil {
			return errors.Wrap(err, "Failed to verify Reject of Request", fields)
		}

		// Request is invalid, Reject was submitted successfully.
		return nil
	}

	err = s.processValidPendingRequest(ctx, withdrawAddress, withdrawAmount, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process valid pending Request", fields)
	}

	return nil
}

func (s *Service) sendRequestToVerify(urlSuffix string, request interface{}) (*xdr.TransactionEnvelope, error) {
	rawRequestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal RequestBody")
	}

	url, err := s.getVerifyURL()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get URL of Verify")
	}
	url = url + urlSuffix

	fields := logan.F{
		"verify_url":       url,
		"raw_request_body": string(rawRequestBody),
	}

	bodyReader := bytes.NewReader(rawRequestBody)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Review Request to Verify", fields)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send the request", fields)
	}
	fields["status_code"] = resp.StatusCode

	defer func() { _ = resp.Body.Close() }()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the body of response from Verify", fields)
	}
	fields["response_body"] = string(responseBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.From(ErrBadStatusFromVerify, fields)
	}

	envelopeResponse := withdraw.EnvelopeResponse{}
	err = json.Unmarshal(responseBytes, &envelopeResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response body", fields)
	}

	envelope := xdr.TransactionEnvelope{}
	err = envelope.Scan(envelopeResponse.Envelope)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Scan TransactionEnvelope from response from Verify", fields)
	}

	return &envelope, nil
}

func (s *Service) signAndSubmitEnvelope(ctx context.Context, envelope xdr.TransactionEnvelope) error {
	signedEnvelope, err := s.xdrbuilder.Sign(&envelope, s.config.SignerKP)
	if err != nil {
		return errors.Wrap(err, "Failed to sign Envelope")
	}

	envelopeBase64, err := xdr.MarshalBase64(signedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal fully signed Envelope")
	}
	submitResult := s.horizon.Submitter().Submit(ctx, envelopeBase64)
	if submitResult.Err != nil {
		return errors.Wrap(err, "Error submitting signed Envelope to Horizon", logan.F{"submit_result": submitResult})
	}

	return nil
}

func (s *Service) getVerifyURL() (string, error) {
	services, err := s.discovery.DiscoverService(conf.ServiceBTCWithdrawVerify)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", conf.ServiceBTCWithdrawVerify))
	}
	if len(services) == 0 {
		return "", ErrNoVerifyServices
	}

	return services[0].Address, nil
}
