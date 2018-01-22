package notifier

import (
	"net/http"
	"net/url"
	"time"

	"bytes"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/signcontrol"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/operations"
)

const (
	horizonPathParticipants = "/participants"
	horizonPathPayments     = "/payments"
)

func (s *Service) getParticipants(requestData *operations.ParticipantsRequest) (map[int64][]operations.ApiParticipant, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&requestData); err != nil {
		return nil, errors.Wrap(err, "failed to marshal")
	}

	response, err := s.horizon.WithSigner(s.Signer).Client().Post(horizonPathParticipants, &buf)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	var result map[int64][]operations.ApiParticipant
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	s.logger.Debug("got participants")
	return result, nil
}

func (s *Service) paymentsRequest() (*http.Request, error) {
	u, err := url.Parse(horizonPathPayments)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("order", "asc")
	q.Set("cursor", s.Operations.Cursor)

	if s.Operations.Cursor == "" {
		duration, err := time.ParseDuration(s.Operations.IgnoreOlderThan)
		if err != nil {
			return nil, err
		}
		q.Set("since", time.Now().UTC().Add(-duration).Format(time.RFC3339))
	}
	u.RawQuery = q.Encode()

	s.logger.WithField("cursor", s.Operations.Cursor).WithField("path", u.String()).Warn("Remake request")
	request, err := http.NewRequest("GET", u.String(), nil)
	if err := signcontrol.SignRequest(request, s.Signer); err != nil {
		return nil, errors.Wrap(err, "failed to sign request")
	}
	return request, nil
}
