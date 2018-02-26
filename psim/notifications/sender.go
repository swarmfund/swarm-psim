package notifications

import (
	"net/http"
	"bytes"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type SlackSender struct {
	config SlackConfig

	client http.Client
}

func NewSlackSender(config SlackConfig) *SlackSender {
	return &SlackSender{
		config: config,

		client: http.Client{},
	}
}

func (s *SlackSender) Send(message string) error {
	request := SlackRequest{
		s.config.ChannelName,
		"psim_notifications",
		message,
		s.config.IconEmoji,
	}
	fields := logan.F{
		"slack_request": request,
	}

	rawRequestBody, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal request to Slack", fields)
	}

	resp, err := s.client.Post(s.config.Url, "", bytes.NewReader(rawRequestBody))
	if err != nil {
		return errors.Wrap(err, "Failed to send http POST request to Slack", fields)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.From(errors.New("Unsuccessful status code in response from Slack"), fields.Merge(logan.F{
			"status_code": resp.StatusCode,
		}))
	}

	return nil
}
