package regources

import (
	"time"
)

// Contract represents singe contract entry with attached invoices
type Contract struct {
	ID            string                   `json:"id"`
	PT            string                   `json:"paging_token"`
	Contractor    string                   `json:"contractor"`
	Customer      string                   `json:"customer"`
	Escrow        string                   `json:"escrow"`
	Disputer      string                   `json:"disputer,omitempty"`
	StartTime     time.Time                `json:"start_time"`
	EndTime       time.Time                `json:"end_time"`
	Details       []map[string]interface{} `json:"details"`
	Invoices      []ReviewableRequest      `json:"invoices,omitempty"`
	DisputeReason map[string]interface{}   `json:"dispute_reason,omitempty"`
	State         []Flag                   `json:"state"`
}

func (c Contract) PagingToken() string {
	return c.PT
}
