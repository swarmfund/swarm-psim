package resource

import "time"

type PayoutState struct {
	Successful bool      `json:"successful"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
}
