package bitcoin

// ConnectorConfig is structure to parse config for NodeConnector into.
type ConnectorConfig struct {
	// TODO Create Node object inside
	NodeIP         string `fig:"node_host,required"`
	NodePort       int    `fig:"node_port,required"`
	NodeAuthKey    string `fig:"node_auth_key,required"`
	Testnet        bool   `fig:"testnet"`
	RequestTimeout int    `fig:"request_timeout_s"`
}
