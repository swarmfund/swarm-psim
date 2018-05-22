package ethwithdveri

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
func (s *Service) detectNewETHSequence(ctx context.Context) (uint64, error) {
	lastUsedSequenceFromCore := s.obtainETHSequenceFromCore(ctx)
	if running.IsCancelled(ctx) {
		return 0, nil
	}

	var newETHSequence uint64
	// Not found in core - need to look in ETH blockchain
	running.UntilSuccess(ctx, s.log, "pending_nonce_obtainer", func(ctx context.Context) (bool, error) {
		nonce, err := s.ethClient.PendingNonceAt(ctx, s.ethAddress)
		if err != nil {
			return false, errors.Wrap(err, "Failed to obtain PendingNonce for ETH Address", logan.F{
				"eth_address": s.ethAddress.String(),
			})
		}

		newETHSequence = nonce
		return true, nil
	}, 10*time.Second, time.Hour)
	if running.IsCancelled(ctx) {
		return 0, nil
	}

	if lastUsedSequenceFromCore == 0 {
		// No ETH TX sequence in Core
		s.log.WithField("new_eth_sequence", newETHSequence).Info("Successfully found ETH sequence(nonce) in ETH blockchain" +
			" (not found among raw TXs from PreConfirmationDetails of WithdrawRequests in Core).")
		return newETHSequence, nil
	}

	if lastUsedSequenceFromCore+1 < newETHSequence {
		// There are some Transactions in ETH blockchain, which are not known for Core.
		return 0, errors.From(errors.New("Last used sequence found in ETH blockchain is grater than last used sequence from Core."), logan.F{
			"last_used_sequence_from_core": lastUsedSequenceFromCore,
			"last_used_sequence_in_eth":    newETHSequence - 1,
		})
	}

	// Success info log about finding ETH sequence in Core is inside the `obtainETHSequenceFromCore()`.
	return lastUsedSequenceFromCore + 1, nil
}

// TODO think of making some common helpers for ethwithdraw and ethwithveri services, as only difference is getTX1/getTX2 function
// ObtainETHSequenceFromCore returns last used ETH TX sequence.
func (s *Service) obtainETHSequenceFromCore(ctx context.Context) uint64 {
	s.log.Info("Starting looking for ETH TXs sequence(nonce) among all WithdrawalRequests from Core.")
	requestsEvents := s.withdrawRequestsStreamer.StreamWithdrawalRequestsOfAsset(ctx, s.config.Asset, true, false)

	logger := s.log.WithField("runner", "last_eth_sequence_detector")
	var ethSequence uint64
	// TODO Change to WithBackOff and exit it by closing new ctx from inside of WithBackOff
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

			_, tx, err := getTX2(*request)
			if err != nil {
				return false, errors.Wrap(err, "Failed to get hex of raw ETH TX2", fields)
			}
			if tx == nil {
				// skip
				logger.WithFields(fields).Debug("ETH TX2 from WithdrawalRequest is nil.")
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
