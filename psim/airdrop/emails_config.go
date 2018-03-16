package airdrop

type EmailsConfig struct {
	Subject            string `fig:"subject"`
	RequestType        int    `fig:"request_type"`
	RequestTokenSuffix string `fig:"request_token_suffix"`
	TemplateName       string `fig:"template_name"`
	TemplateLinkURL    string `fig:"template_link_url"`
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
