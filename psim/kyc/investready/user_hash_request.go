package investready

type userHashRequest struct {
	AccountID string `json:"account_id"`
	UserHash string `json:"user_hash"`
}

func (r userHashRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
		"user_hash": r.UserHash,
	}
}

func (r userHashRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.UserHash == "" {
		return "user_hash cannot be empty."
	}

	return ""
}
