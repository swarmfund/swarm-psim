package btcdeposit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/tokend/go/amount"
	"context"
)

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	GetBlockCount(context.Context) (uint64, error)
	GetBlock(blockIndex uint64) (*btcutil.Block, error)

	IsTestnet() bool
}

// CommonBTCHelper is BTC specific implementation of the OffchainHelper interface from package deposit.
type CommonBTCHelper struct {
	log *logan.Entry

	depositAsset     string
	minDepositAmount uint64
	fixedDepositFee  uint64
	netParams        *chaincfg.Params

	btcClient BTCClient
}

// NewBTCHelper is constructor for CommonBTCHelper.
func NewBTCHelper(
	log *logan.Entry,

	depositAsset string,
	minDepositAmount uint64,
	fixedDepositFee uint64,
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
		log: log,

		depositAsset:     depositAsset,
		minDepositAmount: minDepositAmount,
		fixedDepositFee:  fixedDepositFee,
		netParams:        netParams,

		btcClient: btcClient,
	}, nil
}

// TODO Add context as argument + into interface
func (h CommonBTCHelper) GetLastKnownBlockNumber() (uint64, error) {
	return h.btcClient.GetBlockCount(context.TODO())
}

func (h CommonBTCHelper) GetBlock(number uint64) (*deposit.Block, error) {
	block, err := h.btcClient.GetBlock(number)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block from BTCClient")
	}
	// TODO Handle absent Block

	var depositTXs []deposit.Tx

	for _, tx := range block.Transactions() {
		depositTX := h.parseTX(*tx)
		depositTXs = append(depositTXs, depositTX)
	}

	return &deposit.Block{
		Hash:      block.Hash().String(),
		Timestamp: block.MsgBlock().Header.Timestamp,
		TXs:       depositTXs,
	}, nil

}

func (h CommonBTCHelper) parseTX(tx btcutil.Tx) deposit.Tx {
	var depositOuts []deposit.Out

	for i, out := range tx.MsgTx().TxOut {
		depositOut, err := h.parseOut(*out)
		if err != nil {
			// Don't return any errors, as there can be strange Outputs in Bitcoin blockchain - it's OK.
			h.log.WithFields(logan.F{
				"tx_hash": tx.Hash().String(),
				"out_i":   i,
			}).WithError(err).Debug("Failed to parse TX Output.")
		}

		// Indexes of outputs must be strict, so never loose any Outputs, even though they're malformed.

		if depositOut == nil {
			// Just an empty Output.
			depositOut = &deposit.Out{}
		}

		depositOuts = append(depositOuts, *depositOut)
	}

	return deposit.Tx{
		Hash: tx.Hash().String(),
		Outs: depositOuts,
	}
}

func (h CommonBTCHelper) parseOut(out wire.TxOut) (*deposit.Out, error) {
	scriptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, h.netParams)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to extract PK script Addresses from TX Output")
	}

	if scriptClass != txscript.PubKeyHashTy {
		// Output, which pays not to a pub-key-hash Address - just ignoring.
		// We only accept deposits to our Addresses which are all actually pay-to-pub-key-hash addresses.
		return nil, nil
	}

	addr58 := addrs[0].String()

	return &deposit.Out{
		Address: addr58,
		Value:   uint64(out.Value),
	}, nil
}

func (h CommonBTCHelper) GetMinDepositAmount() uint64 {
	return h.minDepositAmount
}

func (h CommonBTCHelper) GetFixedDepositFee() uint64 {
	return h.fixedDepositFee
}

func (h CommonBTCHelper) ConvertToSystem(offchainAmount uint64) (systemAmount uint64) {
	return uint64(float64(offchainAmount) * (float64(amount.One) / 100000000.0))
}

func (h CommonBTCHelper) GetAsset() string {
	return h.depositAsset
}

func (h CommonBTCHelper) BuildReference(blockNumber uint64, txHash, offchainAddress string, outIndex uint, maxLen int) string {
	var firstNewStyleBlock uint64

	if h.btcClient.IsTestnet() {
		firstNewStyleBlock = 1287181
	} else {
		// Last deprecated deposit we processed was in the 507222 BTC Block.
		firstNewStyleBlock = 507223
	}

	if blockNumber < firstNewStyleBlock {
		// Deprecated approach
		reference := txHash + string(outIndex)

		if len(reference) > maxLen {
			reference = reference[len(reference)-maxLen:]
		}

		return reference
	} else {
		// New approach
		base := fmt.Sprintf("%s:%d", txHash, outIndex)
		hash := sha256.Sum256([]byte(base))
		return hex.EncodeToString(hash[:])
	}
}

func (h CommonBTCHelper) GetAddressSynonyms(address string) []string {
	// No Address synonyms are considered in Bitcoin, base 58 is used for encoding Addresses into string.
	return []string{address}
}

