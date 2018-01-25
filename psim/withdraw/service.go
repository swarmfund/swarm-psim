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
	ErrNoVerifierServices    = errors.New("No Withdraw Verify services were found.")
	ErrBadStatusFromVerifier = errors.New("Unsuccessful status code from Verify.")
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	WithdrawalRequests(result chan<- horizon.Request) <-chan error
}

// OffchainHelper is the interface for specific Offchain(BTC or ETH) withdraw services to implement
// and parametrise the Service.
type OffchainHelper interface {
	// GetAsset must return string code of the Offchain asset the Helper works with.
	// The returned asset is used to filter WithdrawRequests.
	GetAsset() string
	// GetMinWithdrawAmount must return the threshold Withdraw amount value,
	// WithdrawRequests with amount less than provided by this method will be Rejected.
	GetMinWithdrawAmount() float64

	// ValidateAddress must return non-nil error if
	// provided string representation of Address is invalid for the Offchain network.
	//
	// This method is used during validation of WithdrawRequest and Requests with
	// invalid Addresses will be rejected.
	ValidateAddress(addr string) error
	// ValidateTx must return string explanation of validation error,
	// if the TX pays money not to the `withdrawAddress` or amount is bigger than `withdrawAmount`.
	//
	// If was impossible to validate the TX on some reason - method must return non-nil error,
	// otherwise returned error must be nil.
	ValidateTx(tx string, withdrawAddress string, withdrawAmount float64) (string, error)

	// CreateTX must prepare full transaction, without only signatures, everything else must be ready.
	// This TX is used to put into core when transforming a TowStepWithdraw into Withdraw.
	CreateTX(withdrawAddr string, withdrawAmount float64) (tx string, err error)
	// SignTX must sign the provided TX. Provided TX has all the data for transaction, except the signatures.
	SignTX(tx string) (string, error)
	// SendTX must spread the TX into Offchain network and return hash of already transmitted TX.
	SendTX(tx string) (txHash string, err error)
}

// Service is abstract withdraw service, which approves or rejects WithdrawRequest,
// communicating with withdraw verify service for multisig.
//
// To do all the Offchain(BTC or ETH) stuff, Service uses offchainHelper (implementor of the OffchainHelper interface).
type Service struct {
	verifierServiceName string
	signerKP            keypair.Full
	log                 *logan.Entry
	requestListener     RequestListener
	horizon             *horizon.Connector
	xdrbuilder          *xdrbuild.Builder
	discovery           *discovery.Client
	offchainHelper      OffchainHelper

	requests              chan horizon.Request
	requestListenerErrors <-chan error
}

// New is constructor for Service.
func New(
	serviceName string,
	verifierServiceName string,
	signerKP keypair.Full,
	log *logan.Entry,
	requestListener RequestListener,
	horizonConnector *horizon.Connector,
	builder *xdrbuild.Builder,
	discoveryClient *discovery.Client,
	helper OffchainHelper,
) *Service {

	return &Service{
		verifierServiceName: verifierServiceName,
		signerKP:            signerKP,
		log:                 log.WithField("service", serviceName),
		requestListener:     requestListener,
		horizon:             horizonConnector,
		xdrbuilder:          builder,
		discovery:           discoveryClient,
		offchainHelper:      helper,

		requests: make(chan horizon.Request),
	}
}

// Run is a blocking method, it returns only when ctx closes.
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

func (s *Service) processRequest(ctx context.Context, request horizon.Request) error {
	if ProvePendingRequest(request, nil, s.offchainHelper.GetAsset()) != "" {
		// Not a pending or asset doesn't match.
		return nil
	}
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) &&
		request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// Not a Withdraw at all.
		return nil
	}

	s.log.WithFields(GetRequestLoganFields("request", request)).Debugf("Found pending %s Withdrawal Request.", s.offchainHelper.GetAsset())

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// Only TwoStepWithdrawal can be rejected. If RequestType is already Withdraw - it means that it was PreliminaryApproved and needs Approve.
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
	}

	// processValidPendingRequest knows how to process both TwoStepWithdrawal and Withdraw RequestTypes.
	err := s.processValidPendingRequest(ctx, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process valid pending Request")
	}

	return nil
}

func (s *Service) sendRequestToVerifier(urlSuffix string, request interface{}) (*xdr.TransactionEnvelope, error) {
	rawRequestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal RequestBody")
	}

	url, err := s.getVerifierURL()
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
		return nil, errors.From(ErrBadStatusFromVerifier, fields)
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

func (s *Service) getVerifierURL() (string, error) {
	services, err := s.discovery.DiscoverService(s.verifierServiceName)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", s.verifierServiceName))
	}
	if len(services) == 0 {
		return "", ErrNoVerifierServices
	}

	return services[0].Address, nil
}
