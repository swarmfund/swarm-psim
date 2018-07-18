package investready

type redirectedRequest struct {
	AccountID    string `json:"account_id"`
	KYCRequestID uint64 `json:"kyc_request_id"`
	OauthCode    string `json:"oauth_code"`
}

func (r redirectedRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id":     r.AccountID,
		"kyc_request_id": r.KYCRequestID,
	}
}

func (r redirectedRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.OauthCode == "" {
		return "oauth_code cannot be empty."
	}
	if r.KYCRequestID == 0 {
		return "kyc_request_id cannot be empty."
	}

	return ""
}
