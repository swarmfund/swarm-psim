package utils

import (
	"errors"

	"gitlab.com/tokend/go/strkey"
	"gitlab.com/tokend/go/xdr"
)

var (
	ErrInvalidAddress = errors.New("invalid address")
)

func ParseBalanceID(address string) (balanceID xdr.BalanceId, err error) {
	raw, err := strkey.Decode(strkey.VersionByteBalanceID, address)
	if err != nil {
		return
	}

	if len(raw) != 32 {
		return
	}

	var ui xdr.Uint256
	copy(ui[:], raw)

	return xdr.NewBalanceId(xdr.CryptoKeyTypeKeyTypeEd25519, ui)
}
