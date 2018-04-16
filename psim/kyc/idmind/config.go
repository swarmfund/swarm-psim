package idmind

import (
	"time"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	Connector               ConnectorConfig    `fig:"connector,required"`
	RejectReasons           RejectReasonConfig `fig:"reject_reasons,required"`
	AdminNotifyEmailsConfig EmailConfig        `fig:"emails,required"`
	TemplateLinkURL         string             `fig:"template_link_url,required"`
	AdminEmailsToNotify     []string           `fig:"emails_to_notify,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"connector":            c.Connector,
		"reject_reasons":       c.RejectReasons,
		"emails_config":        c.AdminNotifyEmailsConfig,
		"emails_to_notify_len": len(c.AdminEmailsToNotify),
	}
}

type RejectReasonConfig struct {
	KYCStateRejected        string `fig:"kyc_state_rejected,required"`
	FraudPolicyResultDenied string `json:"fraud_policy_result_denied,required"`
	InvalidKYCData          string `json:"invalid_kyc_data,required"`
}

func (c RejectReasonConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"kyc_state_rejected":         c.KYCStateRejected,
		"fraud_policy_result_denied": c.FraudPolicyResultDenied,
		"invalid_kyc_data":           c.InvalidKYCData,
	}
}

type EmailConfig struct {
	Subject     string        `fig:"subject,required"`
	RequestType int           `fig:"request_type,required"`
	Message     string        `fig:"message,required"`
	SendPeriod  time.Duration `fig:"send_period,required"`
}

func (c EmailConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"subject":      c.Subject,
		"request_type": c.RequestType,
		"message":      c.Message,
		"send_period":  c.SendPeriod,
	}
}
