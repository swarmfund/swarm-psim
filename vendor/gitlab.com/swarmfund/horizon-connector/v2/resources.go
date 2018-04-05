package horizon

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/listener"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources/operations"
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
type SaleDetails = resources.SaleDetails
type User = resources.User
type UserAttributes = resources.UserAttributes
type Balance = resources.Balance
type CheckSaleState = operations.CheckSaleState
type CheckSaleStateResponse = listener.CheckSaleStateResponse
type Blob = resources.Blob
type Document = resources.Document
type Reference = resources.Reference
type Account = resources.Account
