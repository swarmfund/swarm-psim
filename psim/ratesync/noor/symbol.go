package noor

import (
	"encoding/json"
	"time"
)

type Symbol struct {
	Ask        float64   `json:"Ask"`
	Bid        float64   `json:"Bid"`
	UpdateTime time.Time `json:"UpdateTime"`
	Symbol     string    `json:"Symbol"`
	Type       string    `json:"Type"`
}

func (s *Symbol) UnmarshalJSON(data []byte) error {
	type Alias Symbol
	a := &struct {
		UpdateTime string `json:"UpdateTime"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	if a.UpdateTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05", a.UpdateTime)
		if err != nil {
			return err
		}
		s.UpdateTime = t.UTC()
	}

	return nil
}
