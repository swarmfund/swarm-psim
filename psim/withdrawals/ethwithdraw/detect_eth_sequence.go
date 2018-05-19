package ethwithdraw

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

// DetectETHSequence tries to
// DetectETHSequence can take a while for execution
//
// Info log about successful detection of the sequence will happen inside,
// including details of where sequence(nonce) was found (Core/ETH blockchain).
func (s *Service) detectLastETHSequence(ctx context.Context) uint64 {
	ethSequence := s.obtainETHSequenceFromCore(ctx)
	if running.IsCancelled(ctx) {
		return 0
	}
	if ethSequence != 0 {
		return ethSequence
	}

	// Not found in core - need to look in ETH blockchain
	running.UntilSuccess(ctx, s.log, "pending_nonce_obtainer", func(ctx context.Context) (bool, error) {
		nonce, err := s.ethClient.PendingNonceAt(ctx, s.ethAddress)
		if err != nil {
			return false, errors.Wrap(err, "Failed to obtain PendingNonce for ETH Address", logan.F{
				"eth_address": s.ethAddress.String(),
			})
		}

		ethSequence = nonce
		return true, nil
	}, 10*time.Second, time.Hour)

	s.log.WithField("eth_sequence", ethSequence).Info("Successfully found ETH sequence(nonce) in ETH blockchain" +
		" (not found among raw TXs from PreConfirmationDetails of WithdrawRequests in Core).")
	return ethSequence
}

func (s *Service) obtainETHSequenceFromCore(ctx context.Context) uint64 {
	s.log.Info("Starting looking for ETH TXs sequence(nonce) among all WithdrawalRequests from Core.")
	requestsEvents := s.withdrawRequestsStreamer.StreamWithdrawalRequestsOfAsset(ctx, s.config.Asset, true, false)

	logger := s.log.WithField("runner", "last_eth_sequence_detector")
	var ethSequence uint64
	running.UntilSuccess(ctx, s.log, "last_eth_sequence_detector", func(ctx context.Context) (bool, error) {
		select {
		case <-ctx.Done():
			return true, nil
		case requestEvent, ok := <-requestsEvents:
			if !ok {
				// No more WithdrawalRequests, stopping
				return true, nil
			}

			request, err := requestEvent.Unwrap()
			if err != nil {
				return false, errors.Wrap(err, "Received erroneous WithdrawRequestEvent")
			}
			fields := logan.F{
				"request": request,
			}

			// This log is a bit redundant, but necessary sometimes.
			logger.WithFields(fields).Debug("Received WithdrawalRequest.")

			tx, err := getTX1(*request)
			if err != nil {
				return false, errors.Wrap(err, "Failed to get hex of raw ETH TX1", fields)
			}
			if tx == nil {
				// skip
				logger.WithFields(fields).Debug("ETH TX1 from WithdrawalRequest is nil.")
				return false, nil
			}

			addr, err := types.FrontierSigner{}.Sender(tx)
			if err != nil {
				return false, errors.Wrap(err, "Failed to check TX's sender", fields)
			}

			if addr == s.ethAddress {
				ethSequence = tx.Nonce()
				logger.WithFields(fields).WithField("eth_sequence", ethSequence).Info("Successfully found ETH sequence(nonce) among WithdrawalRequest from Core.")
				return true, nil
			}

			// skip
			return false, nil
		}
	}, 0, time.Hour)

	return ethSequence
}
