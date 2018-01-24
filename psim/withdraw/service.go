package withdraw

import (
	"context"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/tokend/keypair"
)

var (
	ErrNoVerifyServices    = errors.New("No Withdraw Verify services were found.")
	ErrBadStatusFromVerify = errors.New("Unsuccessful status code from Verify.")
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	WithdrawalRequests(result chan<- horizon.Request) <-chan error
}

// TODO Comment
type OffchainHelper interface {
	// TODO Comment for all methods
	GetAsset() string
	GetHotWallerAddress() string
	GetMinWithdrawAmount() float64

	ValidateAddress(addr string) error
	ValidateTx(txHex string, withdrawAddress string, withdrawAmount float64) (string, error)

	CreateTX(addr string, amount float64) (txHex string, err error)
	SignTX(txHex string) (string, error)
	SendTX(txHex string) (txHash string, err error)
}

type Service struct {
	verifyServiceName string
	signerKP          keypair.Full `fig:"signer" mapstructure:"signer"`
	log               *logan.Entry
	requestListener   RequestListener
	horizon           *horizon.Connector
	xdrbuilder        *xdrbuild.Builder
	discovery         *discovery.Client
	offchainHelper    OffchainHelper

	requests              chan horizon.Request
	requestListenerErrors <-chan error
}

func New(
	serviceName string,
	verifyServiceName string,
	signerKP keypair.Full,
	log *logan.Entry,
	requestListener RequestListener,
	horizonConnector *horizon.Connector,
	builder *xdrbuild.Builder,
	discoveryClient *discovery.Client,
	helper OffchainHelper,
) *Service {

	return &Service{
		verifyServiceName: verifyServiceName,
		signerKP:          signerKP,
		log:               log.WithField("service", serviceName),
		requestListener:   requestListener,
		horizon:           horizonConnector,
		xdrbuilder:        builder,
		discovery:         discoveryClient,
		offchainHelper:    helper,

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
			return errors.Wrap(err, "Failed to process Withdraw Request", GetRequestLoganFields("request", request))
		}

		return nil
	case err := <-s.requestListenerErrors:
		return errors.Wrap(err, "RequestListener sent error")
	}
}

// FIXME
func (s *Service) processRequest(ctx context.Context, request horizon.Request) error {
	// FIXME Use new RequestType (pending)
	if ProvePendingRequest(request, int32(xdr.ReviewableRequestTypeWithdraw), s.offchainHelper.GetAsset()) != "" {
		return nil
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).Debugf("Found pending %s Withdrawal Request.", s.offchainHelper.GetAsset())

	rejectReason := s.getRejectReason(request)
	if rejectReason != "" {
		s.log.WithField("reject_reason", rejectReason).WithFields(GetRequestLoganFields("request", request)).
			Warn("Got Withdraw Request which is invalid due to the RejectReason.")

		err := s.processRequestReject(ctx, request, rejectReason)
		if err != nil {
			return errors.Wrap(err, "Failed to verify Reject of Request", logan.F{
				"reject_reason": rejectReason,
			})
		}

		// Request is invalid, Reject was submitted successfully.
		return nil
	}

	err := s.processValidPendingRequest(ctx, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process valid pending Request")
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

	envelopeResponse := EnvelopeResponse{}
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
	signedEnvelope, err := s.xdrbuilder.Sign(&envelope, s.signerKP)
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
	services, err := s.discovery.DiscoverService(s.verifyServiceName)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", s.verifyServiceName))
	}
	if len(services) == 0 {
		return "", ErrNoVerifyServices
	}

	return services[0].Address, nil
}
