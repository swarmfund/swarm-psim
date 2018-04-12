package emails

import "time"

type Config struct {
	RequestType           int
	UniquenessTokenSuffix string
	SendPeriod            time.Duration
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"request_type":         c.RequestType,
		"request_token_suffix": c.UniquenessTokenSuffix,
		"send_period":          c.SendPeriod,
	}
}
