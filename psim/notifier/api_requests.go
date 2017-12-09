package notifier

import (
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/psim/psim/notifier/internal/operations"
	"gitlab.com/tokend/psim/psim/notifier/internal/types"
)

const (
	horizonPathAssets       = "/assets"
	horizonPathParticipants = "/participants"
	horizonPathPayments     = "/payments"
)

func (s *Service) getAssetsList() ([]types.Asset, error) {
	request, err := s.horizon.SignedRequest("GET", horizonPathAssets, s.Signer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewSignedRequest")
	}

	var resp []types.Asset
	err = sendRequest(request, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	s.logger.WithField("response", resp).Debug("Got assets")
	return resp, nil
}

func (s *Service) getParticipants(requestData *operations.ParticipantsRequest) (map[int64][]operations.ApiParticipant, error) {
	request, err := s.horizon.SignedRequest("POST", horizonPathParticipants, s.Signer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewSignedRequest")
	}

	err = setJSONBody(request, requestData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set request body")
	}

	var resp map[int64][]operations.ApiParticipant
	err = sendRequest(request, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	s.logger.WithField("response", resp).Debug("Got participants")
	return resp, nil
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
	return s.horizon.SignedRequest("GET", u.String(), s.Signer)
}
