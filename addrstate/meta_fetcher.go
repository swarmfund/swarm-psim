package addrstate

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"context"
)

type MetaFetcher struct {
	requester Requester
	log       *logan.Entry
}

func NewChangesProvider(log *logan.Entry, requester Requester) ChangesProvider {
	return MetaFetcher{
		requester: requester,
		log:       log,
	}.process
}

type TransactionsResponse struct {
	Embedded struct {
		Records []Transaction `json:"records"`
	} `json:"_embedded"`
}

type Transaction struct {
	ResultMetaXDR string `json:"result_meta_xdr"`
}

func (f MetaFetcher) process(ctx context.Context, ledgerSeq int64) <-chan xdr.LedgerEntryChange {
	result := make(chan xdr.LedgerEntryChange)

	go func() {
		for {
			metas, err := f.fetch(ctx, ledgerSeq)
			if err != nil {
				f.log.WithError(err).Error("fetch failed")
				continue
			}
			for _, meta := range metas {
				for _, op := range meta.MustOperations() {
					for _, change := range op.Changes {
						result <- change
					}
				}
			}
			close(result)
			return
		}
	}()

	return result
}

func (f MetaFetcher) fetch(ctx context.Context, ledgerSeq int64) (metas []xdr.TransactionMeta, err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors.FromPanic(rvr)
		}
	}()

	endpoint := "/ledgers/%d/transactions"
	var txsResponse TransactionsResponse
	err = f.requester(ctx, "GET", fmt.Sprintf(endpoint, ledgerSeq), &txsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	for _, tx := range txsResponse.Embedded.Records {
		var meta xdr.TransactionMeta
		if err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &meta); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		metas = append(metas, meta)
	}

	return metas, nil
}
