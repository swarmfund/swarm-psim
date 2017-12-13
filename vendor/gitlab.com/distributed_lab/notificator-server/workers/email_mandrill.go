package workers

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/mandrill"
	"gitlab.com/distributed_lab/notificator-server/types"
)

// MandrillEmail send incoming request via Mandrill mailing service.
func MandrillEmail(request types.Request) bool {
	entry := log.WithField("worker", "mandrill_email")
	entry.WithField("request", request.ID).Info("starting")

	connector := mandrill.NewConnector()
	payload := new(notificator.EmailRequestPayload)
	err := request.Payload.Unmarshal(payload)
	if err != nil {
		entry.WithError(err).Warn("failed to send email")
		return false
	}

	receiver := mandrill.NewReceiver(payload.Destination)

	err = connector.SendEmail(receiver, payload.Subject, payload.Message)
	if err != nil {
		entry.WithError(err).Warn("failed to send email")
	}
	return err == nil
}
