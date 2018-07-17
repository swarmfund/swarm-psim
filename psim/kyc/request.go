package kyc

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

const (
	RequestStatePending int32 = 1
)

func IsGeneral(account horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeGeneral)
}

func IsNotVerified(account horizon.Account) bool {
	return account.AccountTypeI == int32(xdr.AccountTypeNotVerified)
}

func IsUpdateToGeneral(request horizon.Request) bool {
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeGeneral)
}

func IsUpdateToVerified(request horizon.Request) bool {
	return request.Details.KYC.AccountTypeToSet.Int == int(xdr.AccountTypeVerified)
}
