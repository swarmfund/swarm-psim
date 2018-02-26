package withdraw

const (
	VerifyPreliminaryApproveURLSuffix = "/preliminary_approve"
	VerifyApproveURLSuffix            = "/approve"
	VerifyRejectURLSuffix             = "/reject"
)

// WithdrawalRequest is struct for describing identifiers of withdrawal request - ID and Hash.
type WithdrawalRequest struct {
	ID   uint64 `json:"id"`
	Hash string `json:"hash"`
}

func (r WithdrawalRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id":   r.ID,
		"hash": r.Hash,
	}
}

type ApproveRequest struct {
	Request WithdrawalRequest `json:"request"`

	TXHex string `json:"tx_hex"`
}

func (r ApproveRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"request": r.Request,
		"tx_hex":  r.TXHex,
	}
}

// NewApprove is constructor for ApproveRequest.
func NewApprove(requestID uint64, requestHash, btcTXHex string) *ApproveRequest {
	return &ApproveRequest{
		Request: WithdrawalRequest{
			ID:   requestID,
			Hash: requestHash,
		},
		TXHex: btcTXHex,
	}
}

type RejectRequest struct {
	Request WithdrawalRequest `json:"request"`

	RejectReason RejectReason `json:"reject_reason"`
}

func (r RejectRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"request":       r.Request,
		"reject_reason": r.RejectReason,
	}
}

// NewReject is constructor for RejectRequest.
func NewReject(requestID uint64, requestHash string, rejectReason RejectReason) *RejectRequest {
	return &RejectRequest{
		Request: WithdrawalRequest{
			ID:   requestID,
			Hash: requestHash,
		},
		RejectReason: rejectReason,
	}
}

// ExternalDetails is used to marshal and unmarshal external
// details of Withdrawal Details for ReviewRequest Operation
// during approve.
type ExternalDetails struct {
	TXHex   string `json:"tx_hex,omitempty"`
}
