package airdrop

import "gitlab.com/distributed_lab/logan/v3/errors"

type EmailsConfig struct {
	Disabled           bool   `fig:"disabled"`
	Subject            string `fig:"subject,required"`
	RequestType        int    `fig:"request_type,required"`
	RequestTokenSuffix string `fig:"request_token_suffix,required"`
	TemplateName       string `fig:"template_name,required"`
	TemplateLinkURL    string `fig:"template_link_url,required"`
}

func (c EmailsConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"subject":              c.Subject,
		"request_type":         c.RequestType,
		"request_token_suffix": c.RequestTokenSuffix,
		"template_name":        c.TemplateName,
		"template_link_url":    c.TemplateLinkURL,
	}
}

func (c EmailsConfig) Validate() error {
	if c.Disabled {
		return nil
	}
	if len(c.RequestTokenSuffix) == 0 {
		return errors.New("'email_request_token_suffix' in config must not be empty")
	}
	return nil
}
