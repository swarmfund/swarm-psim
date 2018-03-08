package earlybird

import (
	"time"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	Asset            string     `fig:"issuance_asset"`
	Amount           uint64     `fig:"issuance_amount"`
	RegisteredBefore *time.Time `fig:"registered_before"`
	WhiteList        []string   `fig:"white_list"`

	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	EmailSubject            string `fig:"email_subject"`
	EmailRequestType        int    `fig:"email_request_type"`
	EmailRequestTokenSuffix string `fig:"email_request_token_suffix"`
	TemplateName            string `fig:"template_name"`
	TemplateRedirectURL     string `fig:"template_redirect_url"`
}
