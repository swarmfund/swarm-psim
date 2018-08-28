package kyc

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

const (
	RequestStatePending int32 = 1
)

func IsGeneral(account horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeGeneral)
}

func IsVerified(account horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeVerified)
}

func IsNotVerified(account horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeNotVerified)
}

func IsUpdateToGeneral(request regources.ReviewableRequest) bool {
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeGeneral)
}

func IsUpdateToVerified(request regources.ReviewableRequest) bool {
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeVerified)
}
