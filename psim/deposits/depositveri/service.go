package depositveri

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

const (
	RequestStatePending int32 = 1
)

var ErrInvalidRequest = errors.New("request is invalid in some way")

type IssuanceRequestsStreamer interface {
	StreamIssuanceRequests(ctx context.Context, opts horizon.IssuanceRequestStreamingOpts) <-chan horizon.ReviewableRequestEvent
}

type Service struct {
	log *logan.Entry

	source keypair.Address
	signer keypair.Full

	externalSystem     int32
	lastBlocksNotWatch uint64

	horizon          *horizon.Connector
	issuanceStreamer IssuanceRequestsStreamer
	addressProvider  deposit.AddressProvider
	builder          *xdrbuild.Builder
	offchainHelper   deposit.OffchainHelper

	streamingOpts horizon.IssuanceRequestStreamingOpts
}

type Opts struct {
	Log                *logan.Entry
	Source             keypair.Address
	Signer             keypair.Full
	ExternalSystem     int32
	LastBlocksNotWatch uint64
	Horizon            *horizon.Connector
	IssuanceStreamer   IssuanceRequestsStreamer
	AddressProvider    deposit.AddressProvider
	Builder            *xdrbuild.Builder
	OffchainHelper     deposit.OffchainHelper
}

func New(opts Opts) *Service {
	return &Service{
		log: opts.Log,

		source: opts.Source,
		signer: opts.Signer,

		externalSystem:     opts.ExternalSystem,
		lastBlocksNotWatch: opts.LastBlocksNotWatch,

		horizon:          opts.Horizon,
		issuanceStreamer: opts.IssuanceStreamer,
		addressProvider:  opts.AddressProvider,
		builder:          opts.Builder,
		offchainHelper:   opts.OffchainHelper,

		streamingOpts: horizon.IssuanceRequestStreamingOpts{
			StopOnEmptyPage: true,
			ReverseOrder:    false,
			AssetCode:       opts.OffchainHelper.GetAsset(),
			RequestState:    RequestStatePending,
		},
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	running.WithBackOff(ctx, s.log, "process_all_issuances_once", s.processAllIssuancesOnce, 30*time.Second, 5*time.Second, time.Hour)
}

func (s *Service) processAllIssuancesOnce(ctx context.Context) error {
	requestEvents := s.issuanceStreamer.StreamIssuanceRequests(ctx, s.streamingOpts)

	for {
		if running.IsCancelled(ctx) {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-requestEvents:
			if !ok {
				// No more requests
				return nil
			}

			request, err := event.Unwrap()
			if err != nil {
				return errors.Wrap(err, "received error from IssuanceRequest stream")
			}

			if err := s.processRequest(request); err != nil {
				s.log.WithError(err).WithField("request_id", request.ID).Error("failed to process request")
				continue
			}
		}
	}
}

func (s *Service) processRequest(request *regources.ReviewableRequest) error {
	// check request sanity
	if !s.isSaneRequest(request) {
		if err := s.rejectRequest(request, "invalid_in_some_way"); err != nil {
			return errors.Wrap(err, "failed to reject request")
		}
		s.log.WithField("request_id", request.ID).Info("request rejected due to request sanity")
		return nil
	}

	issuance := request.Details.IssuanceCreate
	details := issuance.DepositDetails

	txFindMeta, tx, err := s.offchainHelper.FindTX(
		context.TODO(), details.BlockNumber, details.TXHash)
	if err != nil {
		return errors.Wrap(err, "failed to find the TX")
	}

	if tx == nil {
		if txFindMeta.StopWaiting {
			if err := s.rejectRequest(request, "invalid_in_some_way"); err != nil {
				return errors.Wrap(err, "failed to reject request")
			}
			s.log.WithField("request_id", request.ID).Info("request rejected since tx is not found")
			return nil
		}
		return nil
	}

	// check if we have enough confirmations
	lastKnownBlock, err := s.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return errors.Wrap(err, "failed to get last known Block number")
	}
	if lastKnownBlock-txFindMeta.BlockWhereFound < s.lastBlocksNotWatch {
		// Not enough confirmations
		return nil
	}

	if len(tx.Outs) <= int(details.OutIndex) {
		return errors.Wrap(ErrInvalidRequest, "invalid out length", logan.F{
			"want_index": details.OutIndex,
			"len(outs)":  len(tx.Outs),
		})
	}

	out := tx.Outs[details.OutIndex]

	expectedReference := s.offchainHelper.BuildReference(details.BlockNumber, details.TXHash, out.Address, details.OutIndex, 64)
	if expectedReference != string(*request.Reference) {
		return errors.Wrap(ErrInvalidRequest, "reference mismatch", logan.F{
			"got":              *request.Reference,
			"expected":         expectedReference,
			"reference_inputs": []interface{}{details.BlockNumber, details.TXHash, out.Address, details.OutIndex, 64},
		})
	}

	err = s.checkAddress(context.TODO(), issuance.Receiver, out.Address, txFindMeta.BlockTime)
	if err != nil {
		return errors.Wrap(err, "Offchain address is invalid or failed to check")
	}

	err = s.checkValue(out, issuance.Amount)
	if err != nil {
		return errors.Wrap(err, "Deposit value is invalid or failed to check")
	}

	if err := s.approveRequest(request); err != nil {
		return errors.Wrap(err, "failed to approve request")
	}

	s.log.WithField("request_id", request.ID).Info("request approved")
	return nil
}

func (s *Service) isSaneRequest(request *regources.ReviewableRequest) bool {
	details := request.Details.IssuanceCreate.DepositDetails
	return details.TXHash != ""
}

func (s *Service) checkAddress(ctx context.Context, balanceID, offchainAddress string, ts time.Time) error {
	requestAccountID, err := s.horizon.Balances().AccountID(balanceID)
	if err != nil {
		return errors.Wrap(err, "Failed to get AccountID by Balance from Horizon")
	}
	if requestAccountID == nil {
		return errors.Errorf("No Account was found by provided BalanceID")
	}

	accountID := s.addressProvider.ExternalAccountAt(ctx, ts, s.externalSystem, offchainAddress)
	if accountID == nil {
		return errors.New("No Account was connected with this OffchainAddress at this moment of time.")
	}

	if *accountID != *requestAccountID {
		return errors.Wrap(ErrInvalidRequest, "invalid receiver", logan.F{
			"offchain": *accountID,
			"core":     *requestAccountID,
		})
	}

	return nil
}

func (s *Service) checkValue(out deposit.Out, issuanceAmount regources.Amount) error {
	if out.Value < s.offchainHelper.GetMinDepositAmount() {
		return errors.Errorf("Output value is less than MinDepositAmount (%d) - using offchain precision.", s.offchainHelper.GetMinDepositAmount())
	}

	valueWithoutFee := out.Value - s.offchainHelper.GetFixedDepositFee()
	systemValue := int64(s.offchainHelper.ConvertToSystem(valueWithoutFee))
	if systemValue != int64(issuanceAmount) {
		return errors.Errorf("Invalid Issuance amount, in Core - (%d), in Offchain - (%d) - using system precision.", issuanceAmount, systemValue)
	}

	return nil
}

func (s *Service) rejectRequest(request *regources.ReviewableRequest, rejectReason string) error {
	envelope, err := s.builder.
		Transaction(s.source).
		Op(xdrbuild.ReviewRequestOp{
			ID:      request.ID,
			Hash:    request.Hash,
			Action:  xdr.ReviewRequestOpActionPermanentReject,
			Reason:  rejectReason,
			Details: xdrbuild.IssuanceDetails{},
		}).
		Sign(s.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		return errors.Wrap(result.Err, "failed to submit tx", result.GetLoganFields())
	}

	return nil
}

func (s *Service) approveRequest(request *regources.ReviewableRequest) error {
	envelope, err := s.builder.
		Transaction(s.source).
		Op(xdrbuild.ReviewRequestOp{
			ID:      request.ID,
			Hash:    request.Hash,
			Action:  xdr.ReviewRequestOpActionApprove,
			Details: xdrbuild.IssuanceDetails{},
		}).
		Sign(s.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		return errors.Wrap(result.Err, "failed to submit tx", result.GetLoganFields())
	}

	return nil
}
