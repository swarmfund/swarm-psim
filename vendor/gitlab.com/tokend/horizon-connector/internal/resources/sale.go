package resources

import (
	"encoding/json"
	"strconv"
)

type Sale struct {
	ID      uint64      `json:"id"`
	Details SaleDetails `json:"details"`
}

type SaleDetails struct {
	Name string `json:"name"`
}

func (s *Sale) Name() string {
	return s.Details.Name
}

func (s *Sale) UnmarshalJSON(data []byte) error {
	type Alias Sale
	aux := &struct {
		ID string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	var err error
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.ID, err = strconv.ParseUint(aux.ID, 10, 64)
	if err != nil {
		return err
	}
	return nil
}
