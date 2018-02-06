package deposit

// TODO Comment
type VerifyRequest struct {
	AccountID string `json:"account_id"`
}

func (r VerifyRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"account_id": r.AccountID,
	}
}
