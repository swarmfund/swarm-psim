package withdraw

import (
	"context"
	"time"

	"fmt"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/verification"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

var (
	ErrNoVerifierServices = errors.New("No Withdraw Verify services were found.")
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	StreamWithdrawalRequestsOfAsset(ctx context.Context, destAssetCode string, reverseOrder, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

// RequestsConnector is the interface implemented in HorizonConnector,
// it allows to obtain Request by ID and is used for fetching Withdraw Request
// after preliminary approve (TwoStep Request).
type RequestsConnector interface {
	GetRequestByID(requestID uint64) (*regources.ReviewableRequest, error)
}

type TXSubmitter interface {
	// TODO Change deprecated Submit
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

// CommonOffchainHelper is the interface for specific Offchain(BTC or ETH)
// withdraw and withdraw-verify services to implement
// and parametrise the Service.
type CommonOffchainHelper interface {
	// GetAsset must return string code of the Offchain asset the Helper works with.
	// The returned asset is used to filter WithdrawRequests.
	GetAsset() string
	// GetMinWithdrawAmount must return the threshold Withdraw amount value,
	// WithdrawRequests with amount less than provided by this method will be Rejected.
	//
	// Must use the same precision system the other methods in this interface do
	GetMinWithdrawAmount() int64

	// ValidateAddress must return non-nil error if
	// provided string representation of Address is invalid for the Offchain network.
	//
	// This method is used during validation of WithdrawRequest and Requests with
	// invalid Addresses will be rejected.
	ValidateAddress(addr string) error
	// ValidateTX must return string explanation of validation error,
	// if the TX pays money not to the `withdrawAddress` or amount is bigger than `withdrawAmount`.
	//
	// If was impossible to validate the TX on some reason - method must return non-nil error,
	// otherwise returned error must be nil.
	ValidateTX(tx string, withdrawAddress string, withdrawAmount int64) (string, error)

	// ConvertAmount must convert DestinationAmount of the Withdraw
	// provided in integer with system precision to the amount in integer with Offchain precision
	// such as satoshis int BTC for instance.
	//
	// This value will be used to pass into other methods of this interface,
	// so all methods must use the same precision system.
	//
	// Normally this method should do
	//
	// 		return destinationAmount * ((10^N) / amount.One)
	//
	// where N - is the precision of the Offchain system (8 for Bitcoin - satoshis).
	//
	// Though wouldn't be so easy, if your precision is 18 (like for Ethereum) :/
	// TODO Rename to ConvertToOffchain
	ConvertAmount(destinationAmount int64) int64

	// SignTX must sign the provided TX. Provided offchain TX has all the data for transaction, except the signatures.
	SignTX(tx string) (string, error)
	// GetHash must return the hash of the provided TX.
	GetHash(tx string) (string, error)
}

// OffchainHelper is the interface for specific Offchain(BTC or ETH) withdraw services to implement
// and parametrise the Service.
type OffchainHelper interface {
	CommonOffchainHelper
	// Those Helpers which have background threads must start them on call of this method.
	// Run must be blocking. Only if ctx is cancelled
	// For those Helpers, which don't need any background threads - just return immediately.
	Run(context.Context)

	// CreateTX must prepare full transaction, without only signatures, everything else must be ready.
	// This offchain TX is used to put into core when transforming a TowStepWithdraw into Withdraw.
	CreateTX(ctx context.Context, withdrawAddr string, withdrawAmount int64) (tx string, err error)
	// SendTX must spread the offchain TX into Offchain network and return hash of already transmitted TX.
	SendTX(ctx context.Context, tx string) (txHash string, err error)
}

type VerificationConfig struct {
	Verify              bool
	SourceKP            keypair.Address // Needed only if Verify is false.
	VerifierServiceName string
	Discovery           *discovery.Client
}

// Service is abstract withdraw service, which approves or rejects WithdrawRequest,
// communicating with withdraw verify service for multisig.
//
// To do all the Offchain(BTC or ETH) stuff, Service uses offchainHelper (implementor of the OffchainHelper interface).
type Service struct {
	signerKP keypair.Full

	log               *logan.Entry
	requestListener   RequestListener
	requestsConnector RequestsConnector
	txSubmitter       TXSubmitter
	xdrbuilder        *xdrbuild.Builder
	verification      VerificationConfig
	offchainHelper    OffchainHelper

	requestEvents <-chan horizon.ReviewableRequestEvent
}

// New is constructor for Service.
func New(
	serviceName string,
	signerKP keypair.Full,
	log *logan.Entry,
	requestListener RequestListener,
	requestsConnector RequestsConnector,
	txSubmitter TXSubmitter,
	builder *xdrbuild.Builder,
	verification VerificationConfig,
	helper OffchainHelper,
) *Service {

	return &Service{
		signerKP:          signerKP,
		log:               log.WithField("service", serviceName),
		requestListener:   requestListener,
		requestsConnector: requestsConnector,
		txSubmitter:       txSubmitter,
		xdrbuilder:        builder,
		verification:      verification,
		offchainHelper:    helper,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	// TODO Wait
	go s.offchainHelper.Run(ctx)

	s.requestEvents = s.requestListener.StreamWithdrawalRequestsOfAsset(ctx, s.offchainHelper.GetAsset(), false, true)

	running.WithBackOff(ctx, s.log, "request_processor", s.listenAndProcessRequest, 0, 5*time.Second, 10*time.Minute)
}

func (s *Service) listenAndProcessRequest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case requestEvent := <-s.requestEvents:
		request, err := requestEvent.Unwrap()
		if err != nil {
			return errors.Wrap(err, "RequestsStreamer sent error")
		}
		fields := logan.F{"request": request}

		err = s.processRequest(ctx, *request)
		if err != nil {
			return errors.Wrap(err, "Failed to process Withdraw Request", fields)
		}

		return nil
	}
}

func (s *Service) processRequest(ctx context.Context, request regources.ReviewableRequest) error {
	proveErr := ProvePendingRequest(request, s.offchainHelper.GetAsset(), int32(xdr.ReviewableRequestTypeTwoStepWithdrawal), int32(xdr.ReviewableRequestTypeWithdraw))
	if proveErr != "" {
		// Not a pending or asset doesn't match.
		s.log.WithFields(logan.F{
			"request_id": request.ID,
			"prove_err":  proveErr,
		}).Debug("Found not interesting Request.")
		return nil
	}

	switch request.Details.RequestType {
	case int32(xdr.ReviewableRequestTypeTwoStepWithdrawal):
		s.log.WithField("request", request).Debugf("Found pending %s TwoStepWithdrawal Request.", s.offchainHelper.GetAsset())
	case int32(xdr.ReviewableRequestTypeWithdraw):
		s.log.WithField("request", request).Debugf("Found pending %s Withdraw Request.", s.offchainHelper.GetAsset())
	}

	if !s.verification.Verify && request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// TwoStepWithdraw Request is unprocessable by withdraw service without Verify.
		return nil
	}

	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		// Only TwoStepWithdrawal can be rejected. If RequestType is already Withdraw - it means that it was PreliminaryApproved and needs Approve.
		rejectReason := s.getRejectReason(request)
		if rejectReason != "" {
			s.log.WithField("reject_reason", rejectReason).WithField("request", request).
				Warn("Got WithdrawalRequest which is invalid due to the RejectReason.")

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
	return errors.Wrap(err, "Failed to process valid pending Request")
}

func (s *Service) sendRequestToVerifier(urlSuffix string, request interface{}) (*xdr.TransactionEnvelope, error) {
	url, err := s.getVerifierURL()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get URL of Verify")
	}
	url = url + urlSuffix

	envelope, err := verification.VerifyRequest(url, request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send request to Verifier", logan.F{"verifier_url": url})
	}

	return envelope, nil
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

	// TODO Fix: Submit method is deprecated
	submitResult := s.txSubmitter.Submit(ctx, envelopeBase64)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{"submit_result": submitResult})
	}

	return nil
}

func (s *Service) getVerifierURL() (string, error) {
	time.Sleep(15 * time.Second)
	services, err := s.verification.Discovery.DiscoverService(s.verification.VerifierServiceName)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to discover %s service.", s.verification.VerifierServiceName))
	}
	if len(services) == 0 {
		return "", ErrNoVerifierServices
	}

	return services[0].Address, nil
}
