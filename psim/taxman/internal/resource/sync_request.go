package resource

type SyncRequest struct {
	Ledger       int64    `json:"ledger"`
	Transactions []string `json:"transactions"`
}
