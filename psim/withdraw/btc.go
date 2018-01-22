package withdraw

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	BTCAsset = "BTC"
)

// ValidateBTCAddress decodes the string encoding of an Address and returns
// nil if addr is a valid encoding for a known Address type and error otherwise.
func ValidateBTCAddress(addr string, defaultNet *chaincfg.Params) error {
	_, err := btcutil.DecodeAddress(addr, defaultNet)
	return err
}

// IsPendingBTCWithdraw returns true if the Request is of Withdraw type,
// is in pending state
// and its DestinationAsset is BTC.
func IsPendingBTCWithdraw(request horizon.Request) bool {
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// not a withdraw request
		return false
	}

	if request.State != RequestStatePending {
		// State is not pending
		return false
	}

	if request.Details.Withdraw.DestinationAsset != BTCAsset {
		// Withdraw not to BTC.
		return false
	}

	return true
}
