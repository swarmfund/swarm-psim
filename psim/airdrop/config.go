package airdrop

import (
	"time"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	Asset           string     `fig:"issuance_asset"`
	Amount          uint64     `fig:"issuance_amount"`
	RegisteredAfter *time.Time `fig:"registered_after"`
	WhiteList       []string   `fig:"white_list"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	EmailSubject        string `fig:"email_subject"`
	EmailRequestType    int    `fig:"email_request_type"`
	TemplateName        string `fig:"template_name"`
	TemplateRedirectURL string `fig:"template_redirect_url"`
}
