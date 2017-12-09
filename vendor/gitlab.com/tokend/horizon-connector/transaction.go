package horizon

import "time"

type Transaction struct {
	ID            string    `json:"paging_token"`
	Ledger        int64     `json:"ledger"`
	CreatedAt     time.Time `json:"created_at"`
	ResultMetaXDR string    `json:"result_meta_xdr"`
	EnvelopeXDR   string    `json:"envelope_xdr"`
}
