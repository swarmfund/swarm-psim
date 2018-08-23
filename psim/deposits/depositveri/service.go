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

			return errors.Wrap(s.processRequest(ctx, *request), "failed to process IssuanceRequest")
		}
	}
}

func (s *Service) processRequest(ctx context.Context, request regources.ReviewableRequest) error {
	issuance := request.Details.IssuanceCreate

	if string(issuance.Asset) != s.offchainHelper.GetAsset() {
		return nil
	}

	depositDetails := issuance.DepositDetails

	txFindMeta, tx, err := s.offchainHelper.FindTX(ctx, depositDetails.BlockNumber, depositDetails.TXHash)
	if err != nil {
		return errors.Wrap(err, "failed to find the TX")
	}

	if tx == nil {
		if txFindMeta.StopWaiting {
			return errors.Wrap(err, "failed to verify request (TX does not exit in Offchain and we don't expect it to appear)")
		} else {
			// No TX in Offchain yet, but there is a hope to find TX in future - ignoring this Request for now
			return nil
		}
	}

	// TX was found
	lastKnownBlock, err := s.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return errors.Wrap(err, "failed to get last known Block number")
	}
	if lastKnownBlock-txFindMeta.BlockWhereFound < s.lastBlocksNotWatch {
		// Not enough confirmations
		return nil
	}

	if len(tx.Outs) <= int(depositDetails.OutIndex) {
		return errors.Errorf("OutIndex (%d) is invalid, the Offchain TX has only (%d) Outputs.", depositDetails.OutIndex, len(tx.Outs))
	}
	out := tx.Outs[depositDetails.OutIndex]

	expectedReference := s.offchainHelper.BuildReference(depositDetails.BlockNumber, depositDetails.TXHash,
		out.Address, depositDetails.OutIndex, 64)
	if expectedReference != string(*request.Reference) {
		return errors.Errorf("Invalid reference - from Offchain - (%s), in Core - (%s).", expectedReference, *request.Reference)
	}

	err = s.checkAddress(ctx, issuance.Receiver, out.Address, txFindMeta.BlockTime)
	if err != nil {
		return errors.Wrap(err, "Offchain address is invalid or failed to check")
	}

	err = s.checkValue(out, issuance.Amount)
	if err != nil {
		return errors.Wrap(err, "Deposit value is invalid or failed to check")
	}

	return nil
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

	if accountID != requestAccountID {
		return errors.Errorf("Invalid Issuance receiver - in Offchain - (%s), in Core - (%s).", accountID, requestAccountID)
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

func (s *Service) rejectRequest(request regources.ReviewableRequest, rejectReason string) error {
	envelope, err := s.builder.
		Transaction(s.source).
		Op(xdrbuild.ReviewRequestOp{
			ID:     request.ID,
			Hash:   request.Hash,
			Action: xdr.ReviewRequestOpActionReject,
			// TODO Check that nil is OK
			Details: nil,
			Reason:  rejectReason,
		}).
		Sign(s.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	submitDetails, err := s.horizon.Submitter().SubmitE(envelope)
	if err != nil {
		return errors.Wrap(err, "failed to submit tx", logan.F{
			"envelope": envelope,
		})
	}

	if submitDetails.StatusCode < 200 || submitDetails.StatusCode >= 300 {
		return errors.From(errors.New("Error submitting TX."), logan.F{
			"envelope":                envelope,
			"submit_response_details": submitDetails,
		})
	}

	return nil
}
