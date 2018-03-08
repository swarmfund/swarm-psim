package earlybird

import (
	"context"
	"time"

	"bytes"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	notificator "gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/templates"
)

func (s *Service) processEmails(ctx context.Context) {
	s.log.Info("Started processing emails.")

	app.RunOverIncrementalTimer(ctx, s.log, "pending_generals_processor", func(ctx context.Context) error {
		emailsNumber := s.emails.Length()
		if emailsNumber == 0 {
			return nil
		}

		s.log.WithField("emails_number", emailsNumber).Debug("Sending emails.")

		var processedEmails []string
		s.emails.Range(ctx, func(email string) {
			logger := s.log.WithField("email", email)

			err := s.sendEmail(email)
			if err != nil {
				logger.WithError(err).Error("Failed to send email.")
				return
			}

			processedEmails = append(processedEmails, email)
			logger.Info("Sent email successfully.")
		})

		s.emails.Delete(processedEmails)
		return nil
	}, 30*time.Second, 30*time.Second)
}

func (s *Service) sendEmail(emailAddress string) error {
	msg, err := s.buildEmailMessage()
	if err != nil {
		return errors.Wrap(err, "Failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     s.config.EmailSubject,
		Message:     msg,
	}

	uniqueToken := emailAddress + s.config.EmailRequestTokenSuffix
	resp, err := s.notificator.Send(s.config.EmailRequestType, uniqueToken, payload)
	if err != nil {
		return errors.Wrap(err, "Failed to send email via Notificator")
	}

	if !resp.IsSuccess() {
		return errors.New("Unsuccessful response for email sending request.")
	}

	return nil
}

// TODO Cache the html template once and reuse it
func (s *Service) buildEmailMessage() (string, error) {
	fields := logan.F{
		"template_name": s.config.TemplateName,
	}

	t, err := templates.GetHtmlTemplate(s.config.TemplateName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get html Template", fields)
	}

	data := struct {
		Link string
	}{
		Link: s.config.TemplateRedirectURL,
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to execute html Template", fields)
	}

	return buff.String(), nil
}
