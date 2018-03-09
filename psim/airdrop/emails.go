package airdrop

import (
	"bytes"

	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/templates"
)

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

// SendEmail returns false, nil if the emails isn't sending, because it has been sent earlier.
func SendEmail(emailAddress string, config EmailsConfig, notificatorClient NotificatorConnector) (bool, error) {
	msg, err := buildEmailMessage(config.TemplateName, config.TemplateLinkURL)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     config.Subject,
		Message:     msg,
	}

	uniqueToken := emailAddress + config.RequestTokenSuffix
	resp, err := notificatorClient.Send(config.RequestType, uniqueToken, payload)
	if err != nil {
		return false, errors.Wrap(err, "Failed to send email via Notificator")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// The emails has already been sent earlier.
		return false, nil
	}

	if !resp.IsSuccess() {
		return false, errors.From(errors.New("Unsuccessful response for email sending request."), logan.F{
			"notificator_response": resp,
		})
	}

	return true, nil
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
