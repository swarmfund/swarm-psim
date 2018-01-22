package withdraw

// WithdrawalRequest is struct for describing identifiers of withdrawal request - ID and Hash.
type WithdrawalRequest struct {
	ID   uint64 `json:"id"`
	Hash string `json:"hash"`
}

type ApproveRequest struct {
	Request WithdrawalRequest `json:"request"`

	TXHex string `json:"tx_hex"`
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

// EnvelopeResponse is used to pass TX Envelope encoded to base64 between PSIMs.
type EnvelopeResponse struct {
	Envelope string `json:"envelope"`
}
