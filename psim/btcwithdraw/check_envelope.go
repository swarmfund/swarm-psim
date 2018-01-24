package btcwithdraw

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

// CheckPreliminaryApproveEnvelope returns text of error, or empty string if Envelope is valid.
func checkPreliminaryApproveEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash, btcTXHex string) string {
	generalCheck := checkEnvelope(envelope, requestID, requestHash)
	if generalCheck != "" {
		return generalCheck
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if op.Action != xdr.ReviewRequestOpActionApprove {
		return fmt.Sprintf("Invalid ReviewRequestOpAction (%d) expected Approve(%d).", op.Action, xdr.ReviewRequestOpActionApprove)
	}

	extDetails := op.RequestDetails.Withdrawal.ExternalDetails
	btcDetails := withdraw.ExternalDetails{}
	err := json.Unmarshal([]byte(extDetails), &btcDetails)
	if err != nil {
		return fmt.Sprintf("Cannot unmarshal Withdrawal ExternalDetails of Op: (%s).", extDetails)
	}

	if btcDetails.TXHex != btcTXHex {
		return fmt.Sprintf("Invalid BTC TX hex in the Envelope: (%s), expected (%s).", btcDetails.TXHex, btcTXHex)
	}

	return ""
}

// CheckApproveEnvelope returns text of error, or empty string if Envelope is valid.
func checkApproveEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash,
	withdrawAddress string, withdrawAmount float64, changeAddress string, netParams *chaincfg.Params) (txHex, checkErr string) {

	generalCheck := checkEnvelope(envelope, requestID, requestHash)
	if generalCheck != "" {
		return "", generalCheck
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if op.Action != xdr.ReviewRequestOpActionApprove {
		return "", fmt.Sprintf("Invalid ReviewRequestOpAction (%d) expected Approve(%d).", op.Action, xdr.ReviewRequestOpActionApprove)
	}

	extDetails := op.RequestDetails.Withdrawal.ExternalDetails
	btcDetails := withdraw.ExternalDetails{}
	err := json.Unmarshal([]byte(extDetails), &btcDetails)
	if err != nil {
		return "", fmt.Sprintf("Cannot unmarshal Withdrawal ExternalDetails of Op: %s; extDetails: (%s).", err.Error(), extDetails)
	}

	validationErr, err := withdraw.ValidateBTCTx(btcDetails.TXHex, netParams, withdrawAddress, changeAddress, withdrawAmount)
	if err != nil {
		return "", fmt.Sprintf("Failed to validate BTC TX: %s", err.Error())
	}
	if validationErr != "" {
		return "", validationErr
	}

	return btcDetails.TXHex, ""
}

// CheckRejectEnvelope returns text of error, or empty string if Envelope is valid.
func checkRejectEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash string, rejectReason withdraw.RejectReason) string {
	generalCheck := checkEnvelope(envelope, requestID, requestHash)
	if generalCheck != "" {
		return generalCheck
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if op.Action != xdr.ReviewRequestOpActionPermanentReject {
		return fmt.Sprintf("Invalid ReviewRequestOpAction (%d) expected PermanentReject(%d).", op.Action, xdr.ReviewRequestOpActionPermanentReject)
	}

	if string(op.Reason) != string(rejectReason) {
		return fmt.Sprintf("Invalid RejectReason (%s), expected (%s).", op.Reason, rejectReason)
	}

	return ""
}

func checkEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash string) string {
	if len(envelope.Tx.Operations) != 1 {
		return "Number of Operations does not equal 1."
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if uint64(op.RequestId) != requestID {
		return fmt.Sprintf("Invalid Request ID (%d), expected (%d).", envelope.Tx.Operations[0].Body.ReviewRequestOp.RequestId, requestID)
	}

	reqHash := hex.EncodeToString(op.RequestHash[:])
	if reqHash != requestHash {
		return fmt.Sprintf("Invalid Request Hash (%s), expected (%s).", reqHash, requestHash)
	}

	if op.RequestDetails.RequestType != xdr.ReviewableRequestTypeWithdraw {
		return fmt.Sprintf("Invalid RequestType (%d), expected Withdraw(%d).", op.RequestDetails.RequestType, xdr.ReviewableRequestTypeWithdraw)
	}

	return ""
}
