package workers

import (
	notificator "gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/postage"
	"gitlab.com/distributed_lab/notificator-server/types"
)

// PostageEmail send incoming request via Postage mailing service.
func PostageEmail(request types.Request) bool {
	entry := log.WithField("worker", "postage_email")

	entry.WithField("request", request.ID).Info("starting")

	payload := new(notificator.EmailRequestPayload)
	err := request.Payload.Unmarshal(payload)
	if err != nil {
		entry.WithError(err).Error("failed to unmarshal email payload")
		return false
	}

	err = postage.SendEmail(payload.Destination, payload.Subject, payload.Message)
	if err != nil {
		entry.WithError(err).Error("failed to send email")
	}

	return err == nil
}
