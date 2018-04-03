package idmind

import "gitlab.com/tokend/keypair"

type Config struct {
	Connector     ConnectorConfig    `fig:"connector,required"`
	RejectReasons RejectReasonConfig `fig:"reject_reasons"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	WhiteList []string `fig:"white_list"`
	BlackList []string `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"connector":      c.Connector,
		"reject_reasons": c.RejectReasons,

		"white_list_len": len(c.WhiteList),
		"black_list_len": len(c.BlackList),
	}
}

type RejectReasonConfig struct {
	KYCStateRejected        string `fig:"kyc_state_rejected,required"`
	FraudPolicyResultDenied string `json:"fraud_policy_result_denied,required"`
}
