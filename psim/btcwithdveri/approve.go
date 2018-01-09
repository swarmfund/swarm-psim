package btcwithdveri

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"encoding/hex"
	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/swarmfund/psim/psim/btcwithdraw"
	"encoding/json"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"fmt"
	"net/http"
	"gitlab.com/swarmfund/go/amount"
)

func (s *Service) processApproval(w http.ResponseWriter, r *http.Request, withdrawRequest horizonV2.Request, horizonTX *horizon.TransactionBuilder) {
	opBody := horizonTX.Operations[0].Body.ReviewRequestOp
	fields := logan.F{
		"request_id": opBody.RequestId,
		"request_hash": string(opBody.RequestHash[:]),
		"request_action_i": int32(opBody.Action),
		"request_action": opBody.Action.String(),
	}

	extDetails := opBody.RequestDetails.Withdrawal.ExternalDetails
	btcDetails := btcwithdraw.ExternalDetails{}
	err := json.Unmarshal([]byte(extDetails), &btcDetails)
	if err != nil {
		s.log.WithFields(fields).WithField("raw_details", extDetails).WithError(err).
			Warn("Failed to unmarshal Withdrawal ExternalDetails of Op into btcwithdraw.ExternalDetails.")
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
			"Cannot unmarshal Withdrawal ExternalDetails of Op into btcwithdraw.ExternalDetails: %s", extDetails)))
		return
	}

	fields = fields.Merge(logan.F{
		"tx_hex": btcDetails.TXHex,
	})

	validationErr := s.validateApproval(btcDetails.TXHex, withdrawRequest)
	if validationErr != nil {
		s.log.WithFields(fields).WithError(validationErr).Warn("Approval Request is invalid.")
		ape.RenderErr(w, r, problems.Forbidden(validationErr.Error()))
		return
	}

	err = s.processValidApproval(btcDetails.TXHex, horizonTX)
	if err != nil {
		s.log.WithFields(fields).WithError(err).Error("Failed to process valid Approval Request.")
		ape.RenderErr(w, r, problems.ServerError(err))
	}
}

var (
	ErrNoOuts = errors.New("No Outputs in the provided Transaction.")
	// If start withdrawing several requests in a single BTC Transaction - get rid of this Err.
	ErrMoreThanTwoOuts = errors.New("More than 2 Outputs in the provided Transaction.")
	ErrWrongWithdrawAddress = errors.New("Withdraw Address does not match with the one in the WithdrawRequest.")
	ErrWrongWithdrawAmount = errors.New("Withdraw amount does not match with the one in the WithdrawRequest.")
	ErrUnknownChangeAddress = errors.New("Transaction is sending change to an unknown Address.")
)

func (s *Service) validateApproval(txHex string, withdrawRequest horizonV2.Request) error {
	withdrawAddress, err := btcwithdraw.ObtainAddress(withdrawRequest)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain Address of WithdrawalRequest")
	}
	// Divide by precision of the system.
	withdrawSatoshi := withdrawRequest.Details.Withdraw.DestinationAmount * (100000000 / amount.One)

	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return errors.Wrap(err, "Failed to decode txHex into bytes")
	}

	tx, _ := btc.NewTx(txBytes)
	if len(tx.TxOut) == 0 {
		return ErrNoOuts
	}
	if len(tx.TxOut) > 2 {
		return ErrMoreThanTwoOuts
	}

	addr := btc.NewAddrFromPkScript(tx.TxOut[0].Pk_script, s.btcClient.IsTestnet()).String()
	if addr != withdrawAddress {
		return errors.From(ErrWrongWithdrawAddress, logan.F{
			"btc_address": addr,
		})
	}

	if tx.TxOut[0].Value != uint64(withdrawSatoshi) {
		return errors.From(ErrWrongWithdrawAmount, logan.F{
			"btc_amount": tx.TxOut[0].Value,
		})
	}

	if len(tx.TxOut) == 2 {
		// Have change
		changeAddr := btc.NewAddrFromPkScript(tx.TxOut[1].Pk_script, s.btcClient.IsTestnet()).String()
		if changeAddr != s.config.HotWalletAddress {
			return errors.From(ErrUnknownChangeAddress, logan.F{
				"change_address":       changeAddr,
				"known_change_address": s.config.HotWalletAddress,
			})
		}
	}

	return nil
}

func (s *Service) processValidApproval(txHexToSign string, horizonTX *horizon.TransactionBuilder) error {
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

	err = s.submitApproveRequest(horizonTX)
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

func (s *Service) submitApproveRequest(tx *horizon.TransactionBuilder) error {
	err := tx.Sign(s.config.SignerKP).Submit()

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
