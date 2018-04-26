package investready

import "net/http"

type Connector struct {
	config ConnectorConfig

	client       *http.Client
	accessToken  string
	refreshToken string
}

func NewConnector(config ConnectorConfig) *Connector {
	return &Connector{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}
