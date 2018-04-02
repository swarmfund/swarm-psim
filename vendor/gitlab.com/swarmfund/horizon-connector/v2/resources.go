package horizon

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/listener"
)

// don't blame me, just make sure all exported types are really exported

type TransactionEvent = resources.TransactionEvent
type Transaction = resources.Transaction
type Request = resources.Request
type ReviewableRequestEvent = listener.ReviewableRequestEvent
type Info = resources.Info
type Signer = resources.Signer
type Asset = resources.Asset
type Amount = resources.Amount
type Sale = resources.Sale
type User = resources.User
type Balance = resources.Balance
type Blob = resources.Blob
type Document = resources.Document
type Reference = resources.Reference
type Account = resources.Account
