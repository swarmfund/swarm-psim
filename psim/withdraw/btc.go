package withdraw

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

const (
	BTCAsset = "BTC"
)

// TODO Comment
func ValidateBTCTx(txHex string, netParams *chaincfg.Params, withdrawAddress, changeAddress string, withdrawAmount float64) (string, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode txHex into bytes")
	}

	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create BTC TX from hex")
	}

	if len(tx.MsgTx().TxOut) == 0 {
		return "No Outputs in the TX.", nil
	}
	// If start withdrawing several requests in a single BTC Transaction - get rid of this check.
	if len(tx.MsgTx().TxOut) > 2 {
		return fmt.Sprintf("More than 2 Outputs in the TX (%d).", len(tx.MsgTx().TxOut)), nil
	}

	// TODO Move to separate method
	// Addresses of TX Outputs
	txOutAddresses, err := GetBtcTXAddresses(tx, netParams)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Address from Outputs of the BTC TX")
	}
	if len(txOutAddresses) != len(tx.MsgTx().TxOut) {
		return "", fmt.Errorf("number of got Address from Outputs of the BTC TX (%d) doesn't match with the number of TX Outputs (%d)", len(txOutAddresses), len(tx.MsgTx().TxOut))
	}

	// Withdraw Address
	if txOutAddresses[0] != withdrawAddress {
		return fmt.Sprintf("Wrong BTC Address in the first Output of the TX - WithdrawAddress (%s), expected (%s).", txOutAddresses[0], withdrawAddress), nil
	}

	// Change Address
	if len(txOutAddresses) == 2 {
		// Have change
		if txOutAddresses[1] != changeAddress {
			return fmt.Sprintf("Wrong BTC Address in the second Output of the TX - CahngeAddress (%s), expected (%s).", txOutAddresses[1], changeAddress), nil
		}
	}

	// Amount
	// TODO Take into account BTC fee set by user, when it appears in the Core (if it happens)
	if (float64(tx.MsgTx().TxOut[0].Value) / 100000000) > withdrawAmount {
		return fmt.Sprintf("Wrong BTC amount in the first Output of the TX - WithdrawAmount (%d), expected not more than (%.8f).", tx.MsgTx().TxOut[0].Value, withdrawAmount), nil
	}

	return "", nil
}

func GetBtcTXAddresses(tx *btcutil.Tx, netParams *chaincfg.Params) ([]string, error) {
	var result []string

	for i, out := range tx.MsgTx().TxOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, netParams)
		if err != nil {
			return result, errors.Wrap(err, "Failed to extract Addresses from the first TX Out (the withdraw Address)", logan.F{"out_index": i})
		}
		if len(addrs) == 0 {
			return result, errors.Wrap(err, "Extracted empty Addresses array from the first TX Out (the withdraw Address)", logan.F{"out_index": i})
		}

		result = append(result, addrs[0].String())
	}

	return result, nil
}

// ValidateBTCAddress decodes the string encoding of an Address and returns
// nil if addr is a valid encoding for a known Address type and error otherwise.
func ValidateBTCAddress(addr string, defaultNet *chaincfg.Params) error {
	_, err := btcutil.DecodeAddress(addr, defaultNet)
	return err
}

// ProvePendingBTCRequest returns empty string if the Request is:
// - of Withdraw type;
// - in pending state;
// - its DestinationAsset is BTC.
func ProvePendingBTCRequest(request horizon.Request, neededRequestType int32) string {
	if request.State != RequestStatePending {
		// State is not pending
		return fmt.Sprintf("Invalid Request State (%d) expected Pending(%d).", request.State, RequestStatePending)
	}

	if request.Details.RequestType != neededRequestType {
		// not a withdraw request
		return fmt.Sprintf("Invalid RequestType (%d) expected (%d).", request.Details.RequestType, neededRequestType)
	}

	if request.Details.Withdraw.DestinationAsset != BTCAsset {
		// Withdraw not to BTC.
		return fmt.Sprintf("Wrong DestintationAsset (%s) expected BTC(%s).", request.Details.Withdraw.DestinationAsset, BTCAsset)
	}

	return ""
}
