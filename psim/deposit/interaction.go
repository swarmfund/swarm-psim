package deposit

// TODO Comment
type VerifyRequest struct {
	AccountID string `json:"account_id"`
	Envelope  string `json:"envelope"`
}

func (r VerifyRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"account_id": r.AccountID,
		"envelope":   r.Envelope,
	}
}
