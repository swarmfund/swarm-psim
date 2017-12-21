package resources

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/xdr"
)

type Transaction struct {
	CreatedAt     time.Time `json:"created_at"`
	PagingToken   string    `json:"paging_token"`
	ResultMetaXDR string    `json:"result_meta_xdr"`
}

func (tx *Transaction) LedgerChanges() []xdr.LedgerEntryChange {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &meta); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal"))
	}
	var result []xdr.LedgerEntryChange
	for _, op := range meta.MustOperations() {
		for _, change := range op.Changes {
			result = append(result, change)
		}
	}
	return result
}
