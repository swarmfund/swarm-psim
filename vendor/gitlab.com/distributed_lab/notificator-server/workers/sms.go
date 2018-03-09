package workers

import (
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/twilio"
	"gitlab.com/distributed_lab/notificator-server/types"
)

func SMS(request types.Request, cfg conf.Config) bool {
	connector := twilio.NewConnector()
	payload := new(notificator.EmailRequestPayload)

	if err := request.Payload.Unmarshal(&payload); err != nil {
		cfg.Log().WithField("worker", "sms").WithError(err).Warn("twillio failure")
		return false
	}

	userData := cfg.Twilio()
	response, err := connector.SendSMS(payload.Destination, payload.Message, userData.FromNumber, userData.SID, userData.Token)
	if response.IsOK() && err == nil {
		return true
	}
	cfg.Log().WithField("worker", "sms").WithError(err).WithField("response", response.StatusCode).Warn("twillio failure")

	return false
}
