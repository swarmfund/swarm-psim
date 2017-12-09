package snapshoter

import (
	"gitlab.com/tokend/go/hash"
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"encoding/base64"
	"fmt"
)

type payoutType int

const (
	_                             = iota
	payoutTypeReferral payoutType = iota
	payoutTypeToken
)

// reference - creates unique identified of the payout
func reference(
	payoutType payoutType, balance state.BalanceID, ledger int64) string {
	hashed := hash.Hash(
		[]byte(fmt.Sprintf("%d:%s:%d", payoutType, balance, ledger)))
	return base64.StdEncoding.EncodeToString(hashed[:])
}
