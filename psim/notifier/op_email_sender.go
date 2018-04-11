package notifier

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"net/http"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
	"bytes"
	"html/template"
	"gitlab.com/distributed_lab/logan/v3"
	"context"
	"gitlab.com/distributed_lab/running"
)

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

type TemplatesConnector interface {
	Get(id string) ([]byte, error)
}

type OpEmailSender struct {
	subject              string
	templateName         string
	requestType          int
	logger               *logan.Entry
	template             *template.Template
	notificatorConnector NotificatorConnector
	templatesConnector   TemplatesConnector
}

func NewOpEmailSender(
	subject, templateName string,
	requestType int,
	logger *logan.Entry,
	notificatorConnector NotificatorConnector,
	templatesConnector TemplatesConnector,
) (*OpEmailSender, error) {
	var opEmailSender OpEmailSender

	opEmailSender.subject = subject
	opEmailSender.templateName = templateName
	opEmailSender.requestType = requestType
	opEmailSender.logger = logger

	bb, err := templatesConnector.Get(templateName)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to obtain template bytes")
	}

	opEmailSender.template, err = template.New("kyc-created").Parse(string(bb))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse html.Template")
	}

	opEmailSender.notificatorConnector = notificatorConnector

	return &opEmailSender, nil
}

func (ns *OpEmailSender) SendEmail(ctx context.Context, emailAddress, emailUniqueToken string, data interface{}) error {
	var buff bytes.Buffer
	err := ns.template.Execute(&buff, data)
	if err != nil {
		return errors.Wrap(err, "failed to execute html Template", logan.F{
			"template_name": ns.template.Name(),
		})
	}

	msg := buff.String()

	payload := &notificator.EmailRequestPayload{
		Destination: emailAddress,
		Subject:     ns.subject,
		Message:     msg,
	}

	running.UntilSuccess(ctx, ns.logger, "email_sender", func(ctx context.Context) (bool, error) {
		resp, err := ns.notificatorConnector.Send(ns.requestType, emailUniqueToken, payload)

		if err != nil {
			return false, errors.Wrap(err, "failed to send email via Notificator")
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			ns.logger.Info("Email has already been sent, skipping")
			return true, nil
		}

		if !resp.IsSuccess() {
			return false, errors.From(errors.New("unsuccessful response for email sending request"), logan.F{
				"notificator_response": resp,
			})
		}

		ns.logger.Info("Notificator accepted email successfully")

		return true, nil
	}, time.Second, 5*time.Second)

	return nil
}
