package system

import (
	"encoding/json"

	"gitlab.com/tokend/horizon-connector/internal"
	"gitlab.com/tokend/horizon-connector/internal/errors"
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/regources"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

func (q *Q) Info() (info *resources.Info, err error) {
	response, err := q.client.Get("/")
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if err := json.Unmarshal(response, &info); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal info")
	}
	return info, nil
}

func (q *Q) Statistics() (stats *regources.SystemStatistics, err error) {
	response, err := q.client.Get("/statistics")
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if err := json.Unmarshal(response, &stats); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal system stats")
	}
	return stats, nil
}
