package horizon

import (
	"errors"

	"gitlab.com/swarmfund/go/strkey"
	"gitlab.com/swarmfund/go/xdr"
)

func BalanceIDFromStr(addr string) (xdr.BalanceId, error) {
	raw, err := strkey.Decode(strkey.VersionByteBalanceID, addr)
	if err != nil {
		return xdr.BalanceId{}, err
	}

	if len(raw) != 32 {
		return xdr.BalanceId{}, errors.New("invalid address")
	}

	var ui xdr.Uint256
	copy(ui[:], raw)

	return xdr.NewBalanceId(xdr.CryptoKeyTypeKeyTypeEd25519, ui)
}
