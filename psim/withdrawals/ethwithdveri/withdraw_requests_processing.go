package ethwithdveri

import (
	"context"
	"time"

	"math/big"

	"encoding/json"

	"fmt"

	"encoding/hex"

	"strings"

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
	"gitlab.com/tokend/horizon-connector"
)

const (
	Confirmations = 12
)

func (s *Service) processWithdrawRequestsInfinitely(ctx context.Context) {
	s.log.Info("Starting processing(approve/reject) Withdraw Requests.")
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

			notProcessableReason := getRequestNotProcessableReason(*request)
			if notProcessableReason != "" {
				// Not a pending WithdrawRequests
				logger.WithField("request_not_processable_reason", notProcessableReason).Debug("Found not interesting Request.")
				return nil
			}

			logger.Info("Found interesting WithdrawRequest to approve/reject.")

			running.UntilSuccess(ctx, s.log, "pending_request_processor", func(ctx context.Context) (bool, error) {
				err = s.processPendingWithdrawRequest(ctx, *request)
				if err != nil {
					return false, errors.Wrap(err, "Failed to process pending Withdraw Request", fields)
				}

				return true, nil
			}, 5*time.Second, 10*time.Minute)

			return nil
		}
	}, 0, 5*time.Second, time.Hour)
}

// ProcessPendingWithdrawRequest prepares raw signed ETH TX and puts it into Request Approve.
func (s *Service) processPendingWithdrawRequest(ctx context.Context, request horizon.Request) error {
	if request.Details.RequestType == int32(xdr.ReviewableRequestTypeTwoStepWithdrawal) {
		running.UntilSuccess(ctx, s.log.WithField("request_id", request.ID), "withdraw_first_step_approval_waiter",
			func(ctx context.Context) (bool, error) {
				freshRequest, err := s.requestGetter.GetRequestByID(request.ID)
				if err != nil {
					return false, errors.Wrap(err, "Failed to get Request by ID from Horizon", logan.F{"request": request})
				}

				request = *freshRequest
				return request.Details.RequestType == int32(xdr.ReviewableRequestTypeWithdraw), nil
			}, 10*time.Second, time.Minute)
	}
	logger := s.log.WithField("request", request)
	logger.Info("Request has passed first step - ready to process it.")

	withdrawRequest := request.Details.Withdraw

	assetAmount := convertAmount(int64(withdrawRequest.Amount), s.config.AssetPrecision)

	// Reject
	rejectReason := s.getWithdrawRejectReason(request, assetAmount)
	if rejectReason != "" {
		err := s.rejectWithdrawRequest(request, rejectReason)
		if err != nil {
			return errors.Wrap(err, "Failed to reject Withdraw Request (due to initial validation fail)", logan.F{
				"reject_reason": rejectReason,
			})
		}

		logger.WithField("reject_reason", rejectReason).Warn("Rejected Withdraw Request successfully (due to initial validation fail).")
		return nil
	}

	tx1HashI, ok := withdrawRequest.PreConfirmationDetails[TX1HashPreConfirmDetailsKey]
	if !ok {
		return errors.New("Not found raw ETH TX_1 hash in the PreConfirmationDetails of WithdrawRequest.")
	}
	tx1Hash, ok := tx1HashI.(string)
	if !ok {
		return errors.New("Raw ETH TX_1 in the PreConfirmationDetails of WithdrawRequest is not of type string.")
	}
	if _, err := hex.DecodeString(tx1Hash); err != nil {
		return errors.Wrap(err, "Failed to decode bytes from hex of TX1 hash string", logan.F{
			"tx_1_hash_string": tx1Hash,
		})
	}

	transfer := s.waitForTXWithTransfer(ctx, tx1Hash)
	if running.IsCancelled(ctx) {
		return nil
	}

	// All checks, which guarantee panic-safe execution right here are made in `getWithdrawRejectReason()`
	addr := withdrawRequest.ExternalDetails[WithdrawAddressExtDetailsKey].(string)

	rejectReason = s.getTransferRejectReason(transfer, addr, assetAmount)
	if rejectReason != "" {
		err := s.rejectWithdrawRequest(request, rejectReason)
		if err != nil {
			return errors.Wrap(err, "Failed to reject Withdraw Request (due to invalid Transfer)", logan.F{
				"reject_reason": rejectReason,
			})
		}

		logger.WithField("reject_reason", rejectReason).Warn("Rejected Withdraw Request successfully (due to invalid Transfer).")
		return nil
	}

	// Approve
	tx, err := s.prepareSignedETHTx(ctx, transfer.Id)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare raw signed ETH TX hex")
	}
	txHex, err := eth.Marshal(*tx)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal ETH TX into hex")
	}

	err = s.approveWithdrawRequest(request, txHex, hex.EncodeToString(tx.Hash().Bytes()))
	if err != nil {
		return errors.Wrap(err, "Failed to approve Withdraw Request")
	}

	s.newETHSequence += 1

	return nil
}

// WaitForTX is a blocking method, it only returns when TX1 has 12 confirmations or ctx is cancelled.
func (s *Service) waitForTXWithTransfer(ctx context.Context, ethTX1Hash string) Transfer {
	var transfer Transfer

	running.UntilSuccess(ctx, s.log.WithField("eth_tx_hash", ethTX1Hash), "eth_tx_1_12confirmations_waiter", func(ctx context.Context) (bool, error) {
		// TransactionReceipt returns error if TX is still pending
		receipt, err := s.ethClient.TransactionReceipt(ctx, common.HexToHash(ethTX1Hash))
		if err != nil {
			return false, errors.Wrap(err, "Failed to obtain TransactionReceipt")
		}

		if len(receipt.Logs) == 0 {
			return false, errors.New("Obtained TX1 TransactionReceipt from ETH blockchain with empty Logs list.")
		}

		receiptLog := receipt.Logs[0]
		transferID := new(big.Int).SetBytes(receiptLog.Data)

		head, err := s.ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return false, errors.Wrap(err, "Failed to get last Block's header from ETHClient")
		}
		if head.Number.Uint64()-receiptLog.BlockNumber < Confirmations {
			// Need more confirmations for TX1 before we can proceed.
			return false, nil
		}

		transfer, err = s.multisigContractReader.GetPendingTransfer(nil, transferID)
		if err != nil {
			return false, errors.Wrap(err, "Failed to get pending Transfer")
		}

		return true, nil
	}, 5*time.Second, 30*time.Second)

	return transfer
}

func (s *Service) getTransferRejectReason(transfer Transfer, expectedAddress string, expectedAmount *big.Int) string {
	transferTo := strings.ToLower(transfer.To.String())
	// Just in case
	expectedAddress = strings.ToLower(expectedAddress)

	if transferTo != expectedAddress {
		return fmt.Sprintf("Invalid destination Address in Transfer, expected (%s), got (%s).", expectedAddress, transferTo)
	}
	if transfer.Amount.Cmp(expectedAmount) != 0 {
		return fmt.Sprintf("Invalid Amount in Transfer, expected (%s), got (%s).", expectedAmount.String(), transfer.Amount.String())
	}

	var expectedTokenAddr common.Address
	if s.config.TokenAddress != nil {
		expectedTokenAddr = *s.config.TokenAddress
	}
	if expectedTokenAddr != transfer.Token {
		return fmt.Sprintf("Invalid Token Address in Transfer, expected (%s), got (%s).", expectedTokenAddr.String(), transfer.Token.String())
	}

	return ""
}

// PrepareETHTx never returns nil error with empty TX hex.
func (s *Service) prepareSignedETHTx(ctx context.Context, transferID *big.Int) (*types.Transaction, error) {
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
		// GasLimit probably depends on the Contract and method
		GasLimit: s.config.GasLimit,
		Context:  ctx,
	}

	// Same confirmation method for both ETH and ERC20
	tx, err = s.multisigContractWriter.ConfirmTransfer(&opts, transferID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to confirm transfer in MultisigWallet Contract")
	}

	return tx, nil
}

func (s *Service) approveWithdrawRequest(request horizon.Request, rawETHTxHex, ethTXHash string) error {
	newReviewerDetails := make(map[string]string)
	newReviewerDetails[TX2ReviewerDetailsKey] = rawETHTxHex
	newReviewerDetails[TX2HashReviewerDetailsKey] = ethTXHash

	extDetailsBB, err := json.Marshal(newReviewerDetails)
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
		Details: xdrbuild.WithdrawalDetails{
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

	s.log.WithFields(fields).WithField("request", request).Info("Approved Withdraw Request successfully")
	return nil
}

func (s *Service) rejectWithdrawRequest(request horizon.Request, rejectReason string) error {
	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     request.ID,
		Hash:   request.Hash,
		Action: xdr.ReviewRequestOpActionPermanentReject,
		Details: xdrbuild.WithdrawalDetails{
			ExternalDetails: "{}",
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
