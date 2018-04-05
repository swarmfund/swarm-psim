package idmind

import "gitlab.com/tokend/keypair"

type Config struct {
	Connector     ConnectorConfig    `fig:"connector,required"`
	RejectReasons RejectReasonConfig `fig:"reject_reasons,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"connector":      c.Connector,
		"reject_reasons": c.RejectReasons,
	}
}

type RejectReasonConfig struct {
	KYCStateRejected           string `fig:"kyc_state_rejected,required"`
	FraudPolicyResultDenied    string `json:"fraud_policy_result_denied,required"`
	InvalidKYCData             string `json:"invalid_kyc_data,required"`
}

func (c RejectReasonConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"kyc_state_rejected":         c.KYCStateRejected,
		"fraud_policy_result_denied": c.FraudPolicyResultDenied,
		"invalid_kyc_data":           c.InvalidKYCData,
	}
}
