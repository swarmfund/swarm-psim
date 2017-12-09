package workers

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/twilio"
	"gitlab.com/distributed_lab/notificator-server/types"
)

func SMS(request types.Request) (bool, error) {
	connector := twilio.NewConnector()
	payload := new(notificator.EmailRequestPayload)
	err := request.Payload.Unmarshal(payload)
	if err != nil {
		return false, err
	}

	response, err := connector.SendSMS(payload.Destination, payload.Message)

	if response.IsOK() && err == nil {
		return true, nil
	}

	log.WithField("worker", "sms").WithError(err).WithField("response", response.StatusCode).Warn("twillio failure")
	return false, err
}
