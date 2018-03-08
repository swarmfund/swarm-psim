package airdrop

import (
	"bytes"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/templates"
)

type Notificator interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

func SendEmail(emailAddress string, config EmailsConfig, notificatorClient Notificator) error {
	msg, err := buildEmailMessage(config.TemplateName, config.TemplateRedirectURL)
	if err != nil {
		return errors.Wrap(err, "Failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     config.EmailSubject,
		Message:     msg,
	}

	uniqueToken := emailAddress + config.EmailRequestTokenSuffix
	resp, err := notificatorClient.Send(config.EmailRequestType, uniqueToken, payload)
	if err != nil {
		return errors.Wrap(err, "Failed to send email via Notificator")
	}

	// TODO Check 429 statusCode and return nil

	if !resp.IsSuccess() {
		return errors.New("Unsuccessful response for email sending request.")
	}

	return nil
}

// TODO Cache the html template once and reuse it
func buildEmailMessage(templateName, link string) (string, error) {
	fields := logan.F{
		"template_name": templateName,
	}

	t, err := templates.GetHtmlTemplate(templateName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get html Template", fields)
	}

	data := struct {
		Link string
	}{
		Link: link,
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to execute html Template", fields)
	}

	return buff.String(), nil
}
