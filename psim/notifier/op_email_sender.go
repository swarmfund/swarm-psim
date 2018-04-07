package notifier

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"net/http"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"time"
	"gitlab.com/swarmfund/psim/psim/templates"
	"bytes"
	"html/template"
	"gitlab.com/distributed_lab/logan/v3"
	"context"
)

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

type OpEmailSender struct {
	subject              string
	templateName         string
	requestType          int
	logger               *logan.Entry
	template             *template.Template
	notificatorConnector NotificatorConnector
}

func NewOpEmailSender(subject, templateName string, requestType int, logger *logan.Entry, notificatorConnector NotificatorConnector) (*OpEmailSender, error) {
	var opEmailSender OpEmailSender

	opEmailSender.subject = subject
	opEmailSender.templateName = templateName
	opEmailSender.requestType = requestType
	opEmailSender.logger = logger

	var err error
	opEmailSender.template, err = templates.GetHtmlTemplate(templateName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get html Template", logan.F{
			"template_name": templateName,
		})
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

	app.RunUntilSuccess(
		ctx,
		ns.logger,
		"email_sender",
		func(ctx context.Context) error {
			resp, err := ns.notificatorConnector.Send(ns.requestType, emailUniqueToken, payload)

			if err != nil {
				return errors.Wrap(err, "failed to send email via Notificator")
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				ns.logger.Info("Email has already been sent, skipping")
				return nil
			}

			if !resp.IsSuccess() {
				return errors.From(errors.New("unsuccessful response for email sending request"), logan.F{
					"notificator_response": resp,
				})
			}

			ns.logger.Info("Notificator accepted email successfully")

			return nil
		},
		5*time.Second)

	return nil
}
