package ethwithdraw

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// SubmitETHTransactions is a blocking method
func (s *Service) submitETHTransactionsInfinitely(ctx context.Context) {
	// First time out of loop so that start immediately
	s.submitAllETHTransactionsOnce(ctx)

	t := time.NewTicker(time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.submitAllETHTransactionsOnce(ctx)
			continue
		}
	}
}

func (s *Service) submitAllETHTransactionsOnce(ctx context.Context) {
	s.log.Info("Starting iteration of submission all signed raw ETH TXs from WithdrawRequests into ETH blockchain.")
	requestsEvents := s.withdrawRequestsStreamer.StreamWithdrawalRequestsOfAsset(ctx, s.config.Asset, false, false)

	// This doneCtx is used to exit WithBackOff from inside.
	doneCtx, notifyNoMoreRequests := context.WithCancel(ctx)
	running.WithBackOff(doneCtx, s.log, "eth_txs_submitter", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case requestEvent, ok := <-requestsEvents:
			if !ok {
				// No more Requests
				notifyNoMoreRequests()
				return nil
			}

			request, err := requestEvent.Unwrap()
			if err != nil {
				return errors.Wrap(err, "Received erroneous WithdrawRequestEvent from Requests Streamer")
			}
			fields := logan.F{
				"request": request,
			}

			if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
				return nil
			}
			if request.State != RequestStatePending {
				// We are only interested in pending WithdrawRequests
				return nil
			}

			err = s.processPendingWithdrawRequest(ctx, *request)
			if err != nil {
				return errors.Wrap(err, "Failed to process pending WithdrawRequest", fields)
			}

			return nil
		}
	}, 0, 5*time.Second, time.Hour)
}

func (s *Service) processPendingWithdrawRequest(ctx context.Context, request horizon.Request) error {
	tx, err := getTX1(request)
	if err != nil {
		return errors.Wrap(err, "Failed to get hex of raw ETH TX1")
	}
	if tx == nil {
		// Request version is not 2, just skipping it, this service doesn't process old Requests.
		return nil
	}
	fields := logan.F{
		"eth_tx_hash":  tx.Hash().String(),
		"eth_tx_nonce": tx.Nonce(),
	}
	logger := s.log.WithFields(fields).WithField("request_id", request.ID)

	ethTX, isPending, err := s.ethClient.TransactionByHash(ctx, tx.Hash())
	if err != nil {
		errors.Wrap(err, "Failed to get Transaction by hash from ETH blockchain", fields)
	}
	if ethTX != nil && !isPending {
		// Everything is fine, TX is in ETH blockchain - nothing to do
		logger.Debug("Found already submitted ETH TX from WithdrawRequest from Core, skipping it.")
		return nil
	}

	logger.Debug("Found not submitted ETH TX in WithdrawRequest in Core, submitting it.")

	err = s.ethClient.SendTransaction(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "Failed to send Transaction into ETH blockchain", fields)
	}

	s.ensureMined(ctx, tx.Hash())
	logger.Info("Successfully submitted ETH TX into ETH blockchain.")

	return nil
}

func (s *Service) ensureMined(ctx context.Context, hash common.Hash) {
	running.UntilSuccess(ctx, s.log.WithField("eth_tx_hash", hash.String()), "ensure-mined", func(i context.Context) (bool, error) {
		tx, pending, err := s.ethClient.TransactionByHash(ctx, hash)
		if err != nil {
			return false, errors.Wrap(err, "Failed to get TX.")
		}

		if pending {
			return false, nil
		}

		if tx == nil {
			return false, errors.New("Transaction not found.")
		}

		return true, nil
	}, 10*time.Second, 10*time.Second)
}
