package btcwithdraw

import (
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// TODO Consider moving to so common, as this logic is common for BTC and ETH.

// TODO Comment
func ObtainAddress(request horizonV2.Request) (string, error) {
	addrValue, ok := request.Details.Withdraw.ExternalDetails["address"]
	if !ok {
		return "", ErrMissingAddress
	}

	addr, ok := addrValue.(string)
	if !ok {
		return "", errors.From(ErrAddressNotAString, logan.F{"raw_address_value": addrValue})
	}

	return addr, nil
}

// TODO Comment
func ValidateBTCAddress(addr string, defaultNet *chaincfg.Params) error {
	_, err := btcutil.DecodeAddress(addr, defaultNet)
	return err
}
