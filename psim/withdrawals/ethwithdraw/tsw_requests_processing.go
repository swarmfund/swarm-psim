package ethwithdraw

import (
	"context"
	"time"

	"math/big"

	"encoding/json"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/regources"
)

const (
	ETHAssetCode = "ETH"
)

func (s *Service) processTSWRequestsInfinitely(ctx context.Context) {
	s.log.Info("Starting processing(approve/reject) TwoStepWithdraw Requests.")
	requestsEvents := s.withdrawRequestsStreamer.StreamWithdrawalRequestsOfAsset(ctx, s.config.Asset, false, true)

	running.WithBackOff(ctx, s.log, "requests_approver", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case requestEvent := <-requestsEvents:
			// To be sure now work is done once ctx is cancelled.
			if running.IsCancelled(ctx) {
				return nil
			}

			request, err := requestEvent.Unwrap()
			if err != nil {
				return errors.Wrap(err, "Received erroneous WithdrawRequestEvent")
			}
			fields := logan.F{
				"request": request,
			}
			logger := s.log.WithField("request", request)

			// TODO is it safe to skip requests we could not fulfill ATM
			requestIsInteresting := isProcessablePendingRequest(*request)
			if !requestIsInteresting {
				// Not a pending TwoStepWithdrawRequests
				logger.Debug("Found not interesting Request.")
				return nil
			}

			logger.Info("Found interesting TSWRequest to approve/reject.")

			running.UntilSuccess(ctx, s.log, "pending_request_processor", func(ctx context.Context) (bool, error) {
				err = s.processPendingTSWRequest(ctx, *request)
				if err != nil {
					return false, errors.Wrap(err, "Failed to process pending TwoStepWithdraw Request", fields)
				}

				return true, nil
			}, 5*time.Second, 10*time.Minute)

			return nil
		}
	}, 0, 5*time.Second, time.Hour)
}

// ProcessPendingTSWRequest prepares raw signed ETH TX and puts it into Request Approve.
//
// TSWRequest stands from TwoStepWithdraw Request
func (s *Service) processPendingTSWRequest(ctx context.Context, request regources.ReviewableRequest) error {
	tswRequest := request.Details.TwoStepWithdraw

	assetAmount := convertAmount(int64(tswRequest.Amount), s.config.AssetPrecision)

	// Reject
	rejectReason := s.getTSWRejectReason(request, assetAmount)
	if rejectReason != "" {
		err := s.rejectTSWRequest(request, rejectReason)
		if err != nil {
			return errors.Wrap(err, "Failed to reject TwoStepWithdraw Request", logan.F{
				"reject_reason": rejectReason,
			})
		}

		s.log.WithField("request", request).Warn("Rejected TwoStepWithdraw Request successfully", logan.F{
			"request":       request,
			"reject_reason": rejectReason,
		})
		return nil
	}

	// All checks, which guarantee panic-safe execution right here are made in `getTSWRejectReason()`
	addr := tswRequest.ExternalDetails[WithdrawAddressExtDetailsKey].(string)

	// Approve
	tx, err := s.prepareSignedETHTx(ctx, addr, assetAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare raw signed ETH TX hex")
	}
	txHex, err := eth.Marshal(*tx)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal ETH TX into hex")
	}

	err = s.approveTSWRequest(request, txHex, hex.EncodeToString(tx.Hash().Bytes()))
	if err != nil {
		return errors.Wrap(err, "Failed to approve TwoStepWithdraw Request")
	}

	s.newETHSequence += 1

	return nil
}

// PrepareETHTx never returns nil error with empty TX hex.
func (s *Service) prepareSignedETHTx(ctx context.Context, addr string, amount *big.Int) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error

	opts := bind.TransactOpts{
		From:  s.ethAddress,
		Nonce: big.NewInt(int64(s.newETHSequence)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			signedTX, err := s.ethWallet.SignTX(address, tx)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to sign ETH TX with ETHWallet")
			}

			return signedTX, nil
		},

		Value:    big.NewInt(0),
		GasPrice: s.config.GasPrice,
		// GasLimit probably depends on the Contract, for current MultisigWallet Contract GasPrice was 186552 in Ropsten.
		GasLimit: s.config.GasLimit,
		Context:  ctx,
	}

	if s.config.Asset == ETHAssetCode {
		tx, err = s.multisigContract.CreateEtherTransfer(&opts, common.HexToAddress(addr), amount)
	} else {
		tx, err = s.multisigContract.CreateTokenTransfer(&opts, common.HexToAddress(addr), *s.config.TokenAddress, amount)
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create transfer in MultisigWallet Contract")
	}

	return tx, nil
}

// TSWRequest stands from TwoStepWithdraw Request
func (s *Service) approveTSWRequest(request regources.ReviewableRequest, rawETHTxHex, ethTXHash string) error {
	newPreConfirmDetails := make(map[string]interface{})
	newPreConfirmDetails[VersionPreConfirmDetailsKey] = 3
	newPreConfirmDetails[TX1PreConfirmDetailsKey] = rawETHTxHex
	newPreConfirmDetails[TX1HashPreConfirmDetailsKey] = ethTXHash

	extDetailsBB, err := json.Marshal(newPreConfirmDetails)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal ExternalDetails into JSON bytes")
	}
	fields := logan.F{
		"approval_external_details": string(extDetailsBB),
	}

	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.TwoStepWithdrawalDetails{
			ExternalDetails: string(extDetailsBB),
		},
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope", fields)
	}

	_, err = s.txSubmitter.SubmitE(signedEnvelope)
	// TODO Check error, maybe we need to retry here
	if err != nil {
		return errors.Wrap(err, "Error submitting signed Envelope to Horizon")
	}

	s.log.WithFields(fields).WithField("request", request).Info("Approved TwoStepWithdraw Request successfully")
	return nil
}

func (s *Service) rejectTSWRequest(request regources.ReviewableRequest, rejectReason string) error {
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Details: xdrbuild.TwoStepWithdrawalDetails{
			ExternalDetails: "",
		},
		Reason: rejectReason,
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	_, err = s.txSubmitter.SubmitE(signedEnvelope)
	// TODO Check error, maybe we need to retry here
	if err != nil {
		return errors.Wrap(err, "Error submitting signed Envelope to Horizon")
	}

	return nil
}

func convertAmount(tokendAmount int64, assetPrecision uint) *big.Int {
	bigAmount := big.NewInt(tokendAmount)
	oneAsset := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(assetPrecision)), big.NewInt(0))
	return bigAmount.Mul(bigAmount, oneAsset).Div(bigAmount, big.NewInt(amount.One))
}
