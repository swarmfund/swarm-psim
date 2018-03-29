package withdraw

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/go/xdr"
)

// CheckPreliminaryApproveEnvelope returns text of error, or empty string if Envelope is valid.
func checkPreliminaryApproveEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash, offchainTXHex string) string {
	generalCheck := checkEnvelope(envelope, requestID, requestHash, xdr.ReviewableRequestTypeTwoStepWithdrawal)
	if generalCheck != "" {
		return generalCheck
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if op.Action != xdr.ReviewRequestOpActionApprove {
		return fmt.Sprintf("Invalid ReviewRequestOpAction (%d) expected Approve(%d).", op.Action, xdr.ReviewRequestOpActionApprove)
	}

	extDetails := op.RequestDetails.TwoStepWithdrawal.ExternalDetails
	offchainDetails := ExternalDetails{}
	err := json.Unmarshal([]byte(extDetails), &offchainDetails)
	if err != nil {
		return fmt.Sprintf("Cannot unmarshal Withdrawal ExternalDetails of Op: (%s).", extDetails)
	}

	if offchainDetails.TXHex != offchainTXHex {
		return fmt.Sprintf("Invalid Offchain TX hex in the Envelope: (%s), expected (%s).", offchainDetails.TXHex, offchainTXHex)
	}

	return ""
}

// CheckApproveEnvelope returns text of error, or empty string if Envelope is valid.
func (s Service) checkApproveEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash,
	withdrawAddress string, withdrawAmount int64) (txHex, checkErr string) {

	generalCheck := checkEnvelope(envelope, requestID, requestHash, xdr.ReviewableRequestTypeWithdraw)
	if generalCheck != "" {
		return "", generalCheck
	}

	op := envelope.Tx.Operations[0].Body.ReviewRequestOp

	if op.Action != xdr.ReviewRequestOpActionApprove {
		return "", fmt.Sprintf("Invalid ReviewRequestOpAction (%d) expected Approve(%d).", op.Action, xdr.ReviewRequestOpActionApprove)
	}

	extDetails := op.RequestDetails.Withdrawal.ExternalDetails
	offchainDetails := ExternalDetails{}
	err := json.Unmarshal([]byte(extDetails), &offchainDetails)
	if err != nil {
		return "", fmt.Sprintf("Cannot unmarshal Withdrawal ExternalDetails of Op: %s; extDetails: (%s).", err.Error(), extDetails)
	}

	validationErr, err := s.offchainHelper.ValidateTX(offchainDetails.TXHex, withdrawAddress, withdrawAmount)
	if err != nil {
		return "", fmt.Sprintf("Failed to validate Offchain TX: %s", err.Error())
	}
	if validationErr != "" {
		return "", validationErr
	}

	expectedHash, err := s.offchainHelper.GetHash(offchainDetails.TXHex)
	if err != nil {
		// Almost unreal, as ValidateTX method above didn't fail(so Tx is parsable), but just in case.
		return "", fmt.Sprintf("Failed to get hash of the Offchain TX: %s", err.Error())
	}
	if expectedHash != offchainDetails.TXHash {
		return "", fmt.Sprintf("Invalid offchain TXHash in ExternalDetails of Op: (%s); expected TXHash: (%s).",
			offchainDetails.TXHash, expectedHash)
	}

	return offchainDetails.TXHex, ""
}

// CheckRejectEnvelope returns text of error, or empty string if Envelope is valid.
func checkRejectEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash string, rejectReason RejectReason) string {
	generalCheck := checkEnvelope(envelope, requestID, requestHash, xdr.ReviewableRequestTypeTwoStepWithdrawal)
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

func checkEnvelope(envelope xdr.TransactionEnvelope, requestID uint64, requestHash string, requestType xdr.ReviewableRequestType) string {
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

	if op.RequestDetails.RequestType != requestType {
		return fmt.Sprintf("Invalid RequestType (%d), expected (%d).", op.RequestDetails.RequestType, requestType)
	}

	return ""
}
