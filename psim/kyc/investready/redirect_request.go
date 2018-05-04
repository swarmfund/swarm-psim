package investready

type redirectedRequest struct {
	AccountID string `json:"account_id"`
	OauthCode string `json:"oauth_code"`
}

func (r redirectedRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
	}
}

func (r redirectedRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.OauthCode == "" {
		return "oauth_code cannot be empty."
	}

	return ""
}
