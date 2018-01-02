package btcwithdveri

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/distributed_lab/logan/v3"
	"context"
	"encoding/hex"
	"github.com/piotrnar/gocoin/lib/btc"
)

// TODO Pass the Horizon request signed by Withdraw Service.
// TODO Verify, that everything is fine.
func (s *Service) verifyApprove(ctx context.Context, requestID uint64, requestHash, txHexToSign string) error {
	// TODO Verify, that everything is fine.

	// Processing fine Approval
	signedTXHex, err := s.btcClient.SignAllTXInputs(txHexToSign, s.config.HotWalletScriptPubKey, &s.config.HotWalletRedeemScript, s.config.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to sing raw TX")
	}

	fields := logan.F{
		"signed_tx_hex": signedTXHex,
	}

	// Obtaining TX hash
	txBytes, err := hex.DecodeString(signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to decode signed TX hex into bytes")
	}
	signedTXHash := btc.NewSha2Hash(txBytes).String()

	fields["signed_tx_hash"] = signedTXHash

	// Putting signed raw BTC TX and its hex into Horizon - Approving Request.
	err = s.submitApproveRequest(requestID, requestHash, signedTXHash, signedTXHex)
	if err != nil {
		return errors.Wrap(err, "Failed to submit Approve Request to Horizon", fields)
	}

	// Sending BTC TX into Bitcoin blockchain.
	sentTXHash, err := s.btcClient.SendRawTX(signedTXHex)
	if err != nil {
		// Error is not returned intentionally.
		//
		// This problem should be fixed manually.
		// Transactions from approved requests not existing in the Bitcoin blockchain
		// should be submitted once more.
		// This process should probably be automated.
		s.log.WithFields(fields).WithError(err).Error("Failed to send withdraw TX into Bitcoin blockchain.")
		return nil
	}

	fields["sent_tx_hash"] = sentTXHash
	s.log.WithFields(fields).Info("Sent withdraw TX to Bitcoin blockchain successfully.")

	return nil
}

// TODO Used pre-approved data by other Service
func (s *Service) submitApproveRequest(requestID uint64, requestHash, signedTXHash, signedTXHex string) error {
	externalDetails := struct {
		TXHash string `json:"tx_hash"`
		TXHex  string `json:"tx_hex"`
	}{
		TXHash: signedTXHash,
		TXHex:  signedTXHex,
	}
	detailsBytes, err := json.Marshal(externalDetails)
	if err != nil {
		errors.Wrap(err, "Failed to marshal ExternalDetails for OpWithdrawal (containing hex and hash of BTC TX)")
	}

	err = s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.SourceKP,
	}).Op(&horizon.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: horizon.ReviewRequestOpDetails{
			Type: xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
				ExternalDetails: string(detailsBytes),
			},
		},
	}).
		Sign(s.config.SignerKP).
		Submit()

	if err != nil {
		var fields logan.F

		sErr, ok := errors.Cause(err).(horizon.SubmitError)
		if ok {
			fields = logan.F{"horizon_submit_error_response_body": string(sErr.ResponseBody())}
		}

		return errors.Wrap(err, "Failed to submit Transaction to Horizon", fields)
	}

	return nil
}
