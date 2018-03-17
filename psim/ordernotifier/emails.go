package ordernotifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/templates"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"bytes"
	"gitlab.com/distributed_lab/notificator-server/client"
	"net/http"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
	"context"
)

func (s *Service) sendEmail(ctx context.Context, emailAddress, saleName, uniqueToken string) error {
	msg, err := s.buildEmailMessage(saleName)
	if err != nil {
		return errors.Wrap(err, "failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     s.config.Subject,
		Message:     msg,
	}

	app.RunUntilSuccess(ctx, s.logger, "email_sender", func(ctx context.Context) error {
		resp, err := s.emailSender.Send(s.config.RequestType, uniqueToken, payload)

		if err != nil {
			return errors.Wrap(err, "failed to send email via Notificator")
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			s.logger.Info("Email has already been sent, skipping")
			return nil
		}

		if !resp.IsSuccess() {
			return errors.From(errors.New("unsuccessful response for email sending request"), logan.F{
				"notificator_response": resp,
			})
		}

		s.logger.Info("Notificator accepted email successfully")

		return nil
	}, 5*time.Second)

	return nil
}

func (s *Service) buildEmailMessage(saleName string) (string, error) {
	fields := logan.F{
		"template_name": s.config.TemplateName,
	}

	t, err := templates.GetHtmlTemplate(s.config.TemplateName)
	if err != nil {
		return "", errors.Wrap(err, "failed to get html template", fields)
	}

	data := struct {
		Link string
		Fund string
	}{
		Link: s.config.TemplateLinkURL,
		Fund: saleName,
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute html template", fields)
	}

	return buff.String(), nil
}
