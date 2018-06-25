package telegram

type UserRequest struct {
	AccountID string `json:"account_id"`
	TelegramHandle string `json:"telegram_handle"`
}

func (r UserRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
		"telegram_handle": r.TelegramHandle,
	}
}

func (r UserRequest) Validate() (validationErr string) {
	if r.AccountID == "" {
		return "account_id cannot be empty."
	}
	if r.TelegramHandle == "" {
		return "telegram_handle cannot be empty."
	}

	return ""
}
