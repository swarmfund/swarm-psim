package addrstate

import (
	"time"

	"fmt"

	"context"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type LedgersResponse struct {
	Embedded struct {
		Records []Ledger `json:"records"`
	} `json:"_embedded"`
}

type Ledger struct {
	ID       string    `json:"paging_token"`
	Sequence int64     `json:"sequence"`
	ClosedAt time.Time `json:"closed_at"`
	TXCount  int64     `json:"transaction_count"`
}

type LedgersFetcher struct {
	requester Requester
	log       *logan.Entry
}

func NewLedgersProvider(log *logan.Entry, requester Requester) func(ctx context.Context) <-chan Ledger {
	fetcher := LedgersFetcher{
		requester: requester,
		log:       log,
	}
	return fetcher.Run
}

func (f *LedgersFetcher) Run(ctx context.Context) <-chan Ledger {
	result := make(chan Ledger)
	go func() {
		next := "/ledgers?limit=200"
		for {
			next = f.fetch(ctx, result, next)
		}
	}()
	return result
}

func (f *LedgersFetcher) fetch(ctx context.Context, ledgers chan<- Ledger, endpoint string) (next string) {
	defer func() {
		if next == "" {
			next = endpoint
		}
		if rvr := recover(); rvr != nil {
			f.log.WithRecover(rvr).Error("panicked")
		}
	}()

	var response LedgersResponse
	if err := f.requester(ctx, "GET", endpoint, &response); err != nil {
		panic(errors.Wrap(err, "failed to perform request"))
	}

	for _, ledger := range response.Embedded.Records {
		ledgers <- ledger
		next = fmt.Sprintf("/ledgers?cursor=%s&limit=200", ledger.ID)
	}

	return next
}
