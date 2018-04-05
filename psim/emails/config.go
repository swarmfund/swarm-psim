package emails

import "time"

type Config struct {
	Subject               string
	Message               string
	RequestType           int
	UniquenessTokenSuffix string
	SendPeriod            time.Duration
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"subject":              c.Subject,
		"message":              c.Message,
		"request_type":         c.RequestType,
		"request_token_suffix": c.UniquenessTokenSuffix,
		"send_period":          c.SendPeriod,
	}
}
