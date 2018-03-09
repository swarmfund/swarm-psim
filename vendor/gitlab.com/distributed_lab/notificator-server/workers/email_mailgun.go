package workers

import (
	"gitlab.com/distributed_lab/notificator-server/client"

	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/mailgun"
	"gitlab.com/distributed_lab/notificator-server/types"
)

// MailgunEmail send incoming request via Mailgun mailing service.
func MailgunEmail(request types.Request, cfg conf.Config) bool {
	entry := cfg.Log().WithField("worker", "mailgun_email")
	entry.WithField("request", request.ID).Info("starting")

	payload := new(notificator.EmailRequestPayload)
	err := request.Payload.Unmarshal(payload)
	if err != nil {
		entry.WithError(err).Error("failed to unmarshal email payload")
		return false
	}

	mail := cfg.Mailgun()
	resp, id, err := mailgun.SendEmail(
		payload.Destination, payload.Subject, payload.Message,
		mail.From, mail.Domain, mail.Key, mail.PublicKey,
	)
	if err != nil {
		entry.WithError(err).Error("failed to send email")
	}
	entry.WithField("response", resp).WithField("id", id).Debug("sent email")

	return err == nil
}
