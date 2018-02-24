package airdrop

import (
	"time"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	Asset           string     `fig:"issuance_asset"`
	Amount          uint64     `fig:"issuance_amount"`
	RegisteredAfter *time.Time `fig:"registered_after"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	EmailSubject        string `json:"email_subject"`
	EmailRequestType    int    `json:"email_request_type"`
	TemplateName        string `json:"template_name"`
	TemplateRedirectURL string `json:"template_redirect_url"`
}
