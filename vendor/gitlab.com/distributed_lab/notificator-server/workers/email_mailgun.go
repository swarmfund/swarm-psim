package workers

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/notificator-server/log"

	"gitlab.com/distributed_lab/notificator-server/mailgun"
	"gitlab.com/distributed_lab/notificator-server/types"
)

// MailgunEmail send incoming request via Mailgun mailing service.
func MailgunEmail(request types.Request) bool {
	entry := log.WithField("worker", "mailgun_email")
	entry.WithField("request", request.ID).Info("starting")

	payload := new(notificator.EmailRequestPayload)
	err := request.Payload.Unmarshal(payload)
	if err != nil {
		entry.WithError(err).Error("failed to unmarshal email payload")
		return false
	}

	resp, id, err := mailgun.SendEmail(payload.Destination, payload.Subject, payload.Message)
	if err != nil {
		entry.WithError(err).Error("failed to send email")
	}
	entry.WithField("response", resp).WithField("id", id).Debug("sent email")

	return err == nil
}
