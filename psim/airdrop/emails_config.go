package airdrop

import (
	"reflect"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/spf13/cast"
)

type EmailsConfig struct {
	EmailSubject            string `fig:"email_subject"`
	EmailRequestType        int    `fig:"email_request_type"`
	EmailRequestTokenSuffix string `fig:"email_request_token_suffix"`
	TemplateName            string `fig:"template_name"`
	TemplateRedirectURL     string `fig:"template_redirect_url"`
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

