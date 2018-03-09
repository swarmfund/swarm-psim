package airdrop

import (
	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

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

var EmailsHooks = figure.Hooks{
	"airdrop.EmailsConfig": func(raw interface{}) (reflect.Value, error) {
		rawEmails, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
		}

		var emails EmailsConfig
		err = figure.
			Out(&emails).
			From(rawEmails).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to figure out EmailsConfig")
		}

		return reflect.ValueOf(emails), nil
	},
}
