package airdrop

import (
	"bytes"

	"net/http"

	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/templates"
)

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

type EmailProcessor struct {
	log         *logan.Entry
	config      EmailsConfig
	notificator NotificatorConnector

	emails SyncSet
}

func NewEmailsProcessor(
	log *logan.Entry,
	config EmailsConfig,
	notificator NotificatorConnector) *EmailProcessor {

	return &EmailProcessor{
		log:         log.WithField("helper-runner", "emails_processor"),
		config:      config,
		notificator: notificator,

		emails: NewSyncSet(),
	}
}

// Run is clocking function, returns only when ctx cancels.
func (p *EmailProcessor) Run(ctx context.Context) {
	p.log.Info("Started emails processor.")

	app.RunOverIncrementalTimer(ctx, p.log, "emails_processor", func(ctx context.Context) error {
		emailsNumber := p.emails.Length()
		if emailsNumber == 0 {
			p.log.Debug("No emails to send - waiting for next wake up.")
			return nil
		}

		p.log.WithField("emails_number", emailsNumber).Debug("Sending emails.")

		var processedEmails []string
		p.emails.Range(ctx, func(emailAddr string) {
			logger := p.log.WithField("email_addr", emailAddr)

			emailWasSent, err := p.sendEmail(emailAddr)
			if err != nil {
				logger.WithError(err).Error("Failed to send email.")
				return
			}

			processedEmails = append(processedEmails, emailAddr)
			if emailWasSent {
				logger.Info("Notificator accepted email successfully.")
			} else {
				logger.Debug("Email has been already sent earlier - skipping.")
			}
		})

		p.emails.Delete(processedEmails)
		return nil
	}, 30*time.Second, 30*time.Second)
}

func (p *EmailProcessor) AddEmailAddress(ctx context.Context, emailAddress string) {
	p.emails.Put(ctx, emailAddress)
}

func (p *EmailProcessor) sendEmail(emailAddress string) (bool, error) {
	msg, err := buildEmailMessage(p.config.TemplateName, p.config.TemplateLinkURL)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     p.config.Subject,
		Message:     msg,
	}

	uniqueToken := emailAddress + p.config.RequestTokenSuffix
	resp, err := p.notificator.Send(p.config.RequestType, uniqueToken, payload)
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

// SendEmail returns false, nil if the emails isn't sending, because it has been sent earlier.
// DEPRECATED Use EmailProcessor instead.
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
