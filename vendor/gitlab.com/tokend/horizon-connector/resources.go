package horizon

import (
	goresources "gitlab.com/tokend/go/resources"
	"gitlab.com/tokend/horizon-connector/internal/listener"
	"gitlab.com/tokend/horizon-connector/internal/operation"
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/horizon-connector/internal/resources/operations"
	"gitlab.com/tokend/regources"
)

// don't blame me, just make sure all exported types are really exported

type TransactionEvent = resources.TransactionEvent
type TXPacket = listener.TXPacket
type Request = resources.Request
type WithdrawRequest = resources.RequestWithdrawDetails
type KYCRequest = resources.RequestKYCDetails
type ReviewableRequestEvent = listener.ReviewableRequestEvent
type WithdrawalRequestStreamingOpts = listener.WithdrawalRequestStreamingOpts
type Info = resources.Info
type Signer = goresources.Signer
type Sale = resources.Sale
type CoreSale = resources.CoreSale
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
type RequestKYCDetails = resources.RequestKYCDetails

type ReviewableRequestType = operation.ReviewableRequestType
type KYCData = resources.KYCData
type KeyValue = resources.KeyValue

// DEPRECATED: use regources directly
type Asset = regources.Asset

// DEPRECATED: use regources directly
type Amount = regources.Amount
