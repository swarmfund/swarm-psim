package notifications

type SlackRequest struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

func (r SlackRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"channel":    r.Channel,
		"username":   r.Username,
		"text":       r.Text,
		"icon_emoji": r.IconEmoji,
	}
}
