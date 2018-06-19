package btcwithdraw

import (
	"encoding/hex"
	"fmt"

	"crypto/sha256"

	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/tokend/go/amount"
)

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	// CreateAndFundRawTX sets Change position in Outputs to 1.
	CreateAndFundRawTX(goalAddress string, amount float64, changeAddress string, feeRate *float64) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
}

// CommonBTCHelper is BTC specific implementation of the OffchainHelper interface from package withdraw.
type CommonBTCHelper struct {
	log *logan.Entry

	tokendAsset string
	minWithdrawAmount     int64
	hotWalletAddress      string
	hotWalletScriptPubKey string
	hotWalletRedeemScript string
	privateKey            string
	netParams             *chaincfg.Params

	btcClient BTCClient
}

// NewBTCHelper is constructor for CommonBTCHelper.
func NewBTCHelper(
	log *logan.Entry,
	tokendAsset string,
	minWithdrawAmount int64,
	hotWalletAddress,
	hotWalletScriptPubKey,
	hotWalletRedeemScript string,
	privateKey string,
	currency, blockchain string,
	btcClient BTCClient) (*CommonBTCHelper, error) {

	netParams, err := bitcoin.GetNetParams(currency, blockchain)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build NetParams by currency and blockchain", logan.F{
			"currency":   currency,
			"blockchain": blockchain,
		})
	}

	return &CommonBTCHelper{
		// TODO Not actually a helper, but if you suggest a better name - tell me.
		log: log.WithField("service", "btc_helper"),

		tokendAsset: tokendAsset,
		minWithdrawAmount:     minWithdrawAmount,
		hotWalletAddress:      hotWalletAddress,
		hotWalletScriptPubKey: hotWalletScriptPubKey,
		hotWalletRedeemScript: hotWalletRedeemScript,
		privateKey:            privateKey,
		netParams:             netParams,

		btcClient: btcClient,
	}, nil
}

// Run method is needed for withdraw.OffchainHelper interface.
func (h CommonBTCHelper) Run(context.Context) {
	return
}

// GetAsset is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) GetAsset() string {
	return h.tokendAsset
}

// GetMinWithdrawAmount is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) GetMinWithdrawAmount() int64 {
	return h.minWithdrawAmount
}

// ValidateAddress is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) ValidateAddress(addr string) error {
	_, err := btcutil.DecodeAddress(addr, h.netParams)
	return err
}

// ValidateTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) ValidateTX(txHex string, withdrawAddress string, withdrawAmount int64) (string, error) {
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
	txOutAddresses, err := getBtcTXAddresses(tx, h.netParams)
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
		if txOutAddresses[1] != h.hotWalletAddress {
			return fmt.Sprintf("Wrong BTC Address in the second Output of the TX - CahngeAddress (%s), expected (%s).", txOutAddresses[1], h.hotWalletAddress), nil
		}
	}

	// Amount
	// TODO
	// TODO
	// TODO Check that Out.Value + fee <= withdrawAmount
	// TODO
	// TODO
	// TODO Take into account BTC fee set by user, when it appears in the Core (if it happens)
	if tx.MsgTx().TxOut[0].Value > withdrawAmount {
		// TODO Add fee to log.
		return fmt.Sprintf("Wrong BTC amount in the first Output of the TX - WithdrawAmount (%d), expected not more than (%d).", tx.MsgTx().TxOut[0].Value, withdrawAmount), nil
	}

	return "", nil
}

func getBtcTXAddresses(tx *btcutil.Tx, netParams *chaincfg.Params) ([]string, error) {
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

// GetWithdrawAmount is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) ConvertAmount(destinationAmount int64) int64 {
	return destinationAmount * ((100000000) / amount.One)
}

// CreateTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) CreateTX(_ context.Context, addr string, amount int64) (txHex string, err error) {
	floatAmount := float64(amount) / 100000000

	txHex, err = h.btcClient.CreateAndFundRawTX(addr, floatAmount, h.hotWalletAddress, nil)
	if err != nil {
		if errors.Cause(err) == bitcoin.ErrInsufficientFunds {
			return "", errors.Wrap(err, "Could not create raw TX - not enough BTC on hot wallet", logan.F{
				"float_amount": floatAmount,
			})
		}

		return "", errors.Wrap(err, "Failed to create or fund raw TX", logan.F{
			"float_amount": floatAmount,
		})
	}

	return txHex, nil
}

// SignTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) SignTX(txHex string) (string, error) {
	return h.btcClient.SignAllTXInputs(txHex, h.hotWalletScriptPubKey, h.hotWalletRedeemScript, h.privateKey)
}

// SendTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonBTCHelper) SendTX(ctx context.Context, txHex string) (txHash string, err error) {
	txHash, err = h.btcClient.SendRawTX(txHex)

	if errors.Cause(err) == bitcoin.ErrAlreadyInChain {
		h.log.WithFields(logan.F{
			"tx_trying_to_send": txHex,
		}).WithError(err).Warn("Was asked to send TX, got response that it's already in chain.")

		// Already in chain, we have the Hash - just don't tell anyone about this error. Warning is logged though.
		bb, _ := hex.DecodeString(txHex)
		if bb == nil {
			return txHash, err
		}

		tx, _ := btcutil.NewTxFromBytes(bb)
		if tx == nil {
			return txHash, err
		}

		return tx.Hash().String(), nil
	}

	return txHash, err
}

func (h CommonBTCHelper) GetHash(txHex string) (string, error) {
	fmt.Println(txHex)

	txBB, err := hex.DecodeString(txHex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode TX bytes from string")
	}

	t := sha256.Sum256(txBB)
	t = sha256.Sum256(t[:])
	hashBB := t[:]
	reverse(hashBB)

	return hex.EncodeToString(hashBB), nil
}

func reverse(bb []byte) {
	last := len(bb) - 1
	for i := 0; i < len(bb)/2; i++ {
		bb[i], bb[last-i] = bb[last-i], bb[i]
	}
}
