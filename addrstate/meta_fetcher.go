package addrstate

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

type Request func(method, endpoint string, dest interface{}) error

type MetaFetcher struct {
	requester Requester
	log       *logan.Entry
}

func NewChangesProvider(log *logan.Entry, requester Requester) func(ledgerID string) <-chan xdr.LedgerEntryChange {
	return MetaFetcher{
		requester: requester,
		log:       log,
	}.Run
}

type TransactionsResponse struct {
	Embedded struct {
		Records []Transaction `json:"records"`
	} `json:"_embedded"`
}

type Transaction struct {
	ResultMetaXDR string `json:"result_meta_xdr"`
}

func (f MetaFetcher) Run(ledgerID string) <-chan xdr.LedgerEntryChange {
	result := make(chan xdr.LedgerEntryChange)

	go func() {
		for {
			if err := f.fetch(result, ledgerID); err != nil {
				f.log.WithError(err).Error("fetch failed")
			}
		}
	}()

	return result
}

func (f MetaFetcher) fetch(changes chan<- xdr.LedgerEntryChange, ledgerID string) (err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors.FromPanic(rvr)
		}
	}()
	endpoint := "/ledgers/%s/transactions"
	var txsResponse TransactionsResponse
	err = f.requester("GET", fmt.Sprintf(endpoint, ledgerID), &txsResponse)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	for _, tx := range txsResponse.Embedded.Records {
		var txMeta xdr.TransactionMeta
		if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &txMeta); err != nil {
			return errors.Wrap(err, "failed to unmarshal")
		}
		for _, op := range txMeta.MustOperations() {
			for _, change := range op.Changes {
				changes <- change
			}
		}
	}
	return nil
}
