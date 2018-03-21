package ordernotifier

import (
	"bytes"
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/templates"
	"net/http"
	"time"
)

// EmailUnit is structure that consists all required information
// for method "Send" of interface "NotificatorConnector"
type EmailUnit struct {
	Payload     notificator.EmailRequestPayload
	UniqueToken string
}

func (s *Service) emailSenderBuilder(ctx context.Context, emailUnit EmailUnit) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		resp, err := s.emailSender.Send(s.config.RequestType, emailUnit.UniqueToken, emailUnit.Payload)

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
	}
}

func (s *Service) sendEmail(ctx context.Context, emailUnit EmailUnit) error {
	app.RunUntilSuccess(ctx, s.logger, "email_sender", s.emailSenderBuilder(ctx, emailUnit), 5*time.Second)
	return nil
}

func (s *Service) craftEmailUnit(ctx context.Context, emailAddress, saleName, uniqueToken string) (*EmailUnit, error) {
	msg, err := s.buildEmailMessage(saleName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get email message")
	}

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     s.config.Subject,
		Message:     msg,
	}

	emailUnit := &EmailUnit{
		Payload:     *payload,
		UniqueToken: uniqueToken,
	}

	return emailUnit, nil
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
