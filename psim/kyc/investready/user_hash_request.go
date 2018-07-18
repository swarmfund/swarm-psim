package investready

type userHashRequest struct {
	AccountID    string `json:"account_id"`
	KYCRequestID uint64 `json:"kyc_request_id"`
	UserHash     string `json:"user_hash"`
}

func (r userHashRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id":     r.AccountID,
		"kyc_request_id": r.KYCRequestID,
		"user_hash":      r.UserHash,
	}
}

func (r userHashRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.UserHash == "" {
		return "user_hash cannot be empty."
	}
	if r.KYCRequestID == 0 {
		return "kyc_request_id cannot be empty."
	}

	return ""
}
