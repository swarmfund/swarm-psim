package addrstate

import (
	"net/http"

	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
)

type TransactionQ interface {
	Transactions(chan<- horizon.TransactionEvent) <-chan error
}

type Client interface {
	Get(string) (*http.Response, error)
}

type StateMutator func(change xdr.LedgerEntryChange) StateUpdate
