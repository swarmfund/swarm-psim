package horizon

import (
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/horizon-connector/internal/listener"
	"gitlab.com/tokend/horizon-connector/internal/resources/operations"
)

// don't blame me, just make sure all exported types are really exported

type Transaction = resources.Transaction
type TransactionEvent = resources.TransactionEvent
type TXPacket = listener.TXPacket

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
type Wallet = resources.Wallet
type CreateKYCRequestOp = operations.CreateKYCRequest
type CreateKYCRequestOpResponse = listener.CreateKYCRequestOpResponse
type ReviewRequestOp = operations.ReviewRequest
type ReviewRequestOpResponse = listener.ReviewRequestOpResponse
