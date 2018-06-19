package dashwithdraw

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

const (
	txTemplateSize = 20
	inSize         = 260
	outSize        = 21
)

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	GetBlockCount(ctx context.Context) (uint64, error)
	GetBlock(blockNumber uint64) (*btcutil.Block, error)

	CreateRawTX(inputUTXOs []bitcoin.Out, addrToAmount map[string]float64) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)

	EstimateFee(blocksToBeIncluded uint) (float64, error)
}

type CoinSelector interface {
	// AddUTXO must ignore duplications without any exceptions.
	AddUTXO(UTXO)
	// Fund must deactivate (not delete) the UTXOs used for funding.
	// Fund must wait until all existing Blocks are fetched.
	Fund(amount int64) (utxos []bitcoin.Out, change int64, err error)
	TryRemoveUTXO(bitcoin.Out) bool
}

// CommonDashHelper is BTC specific implementation of the OffchainHelper interface from package withdraw.
type CommonDashHelper struct {
	log *logan.Entry

	config    Config
	netParams *chaincfg.Params

	utxoFetched  chan struct{}
	btcClient    BTCClient
	coinSelector CoinSelector
}

// NewDashHelper is constructor for CommonDashHelper.
func NewDashHelper(
	log *logan.Entry,
	config Config,
	btcClient BTCClient,
	coinSelector CoinSelector) (*CommonDashHelper, error) {

	netParams, err := bitcoin.GetNetParams(config.OffchainCurrency, config.OffchainBlockchain)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build NetParams by currency and blockchain", logan.F{
			"currency":   config.OffchainCurrency,
			"blockchain": config.OffchainBlockchain,
		})
	}

	return &CommonDashHelper{
		log: log.WithField("helper", "dash_offchain_helper"),

		config:    config,
		netParams: netParams,

		utxoFetched:  make(chan struct{}),
		btcClient:    btcClient,
		coinSelector: coinSelector,
	}, nil
}

func (h CommonDashHelper) Run(ctx context.Context) {
	h.fetchUTXOsInfinitely(ctx, h.config.FetchUTXOFrom)
}

// GetAsset is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) GetAsset() string {
	return h.config.OffchainCurrency
}

// GetMinWithdrawAmount is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) GetMinWithdrawAmount() int64 {
	return h.config.MinWithdrawAmount
}

// ValidateAddress is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) ValidateAddress(addr string) error {
	_, err := btcutil.DecodeAddress(addr, h.netParams)
	return err
}

// ValidateTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) ValidateTX(txHex string, withdrawAddress string, withdrawAmount int64) (string, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode txHex into bytes")
	}

	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create BTC TX from bytes")
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
	txOutAddresses, err := getBtcTXOutAddresses(tx, h.netParams)
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
		if txOutAddresses[1] != h.config.HotWalletAddress {
			return fmt.Sprintf("Wrong BTC Address in the second Output of the TX - CahngeAddress (%s), expected (%s).", txOutAddresses[1], h.config.HotWalletAddress), nil
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

func getBtcTXOutAddresses(tx *btcutil.Tx, netParams *chaincfg.Params) ([]string, error) {
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
func (h CommonDashHelper) ConvertAmount(destinationAmount int64) int64 {
	return destinationAmount * ((100000000) / amount.One)
}

// CreateTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) CreateTX(ctx context.Context, addr string, amount int64) (txHex string, err error) {
	feePerKB, err := h.btcClient.EstimateFee(h.config.BlocksToBeIncluded)
	if err != nil {
		return "", errors.Wrap(err, "Failed to EstimateFee")
	}
	fields := logan.F{"fee_per_kb": feePerKB}

	if feePerKB > h.config.MaxFeePerKB {
		return "", errors.From(NewErrTooBigFeePerKB(feePerKB, h.config.MaxFeePerKB), fields.Merge(logan.F{
			"config_max_fee_per_kb": h.config.MaxFeePerKB,
		}))
	}

	select {
	case <-ctx.Done():
		return "", nil
	case <-h.utxoFetched:
		// All existing Blocks are already processed - can continue
		break
	}

	inputUTXOs, changeAmount, err := h.coinSelector.Fund(amount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to fund requested amount by CoinSelector")
	}

	txSizeBytes := txTemplateSize + inSize*len(inputUTXOs) + outSize*2 // Always count that there are 2 outputs - just to be simpler
	txFee := (feePerKB / 1000) * float64(txSizeBytes)

	floatAmount := (float64(amount) / 100000000) - txFee

	addrToAmount := map[string]float64{addr: floatAmount}
	if changeAmount > h.config.DustThreshold {
		addrToAmount[h.config.HotWalletAddress] = float64(changeAmount) / 100000000
	}

	txHex, err = h.btcClient.CreateRawTX(inputUTXOs, addrToAmount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create raw TX", logan.F{
			"withdraw_address":      addr,
			"withdraw_float_amount": floatAmount,
			"tx_outputs":            addrToAmount,
			"input_utxos":           inputUTXOs,
			"fee_per_kb":            feePerKB,
			"tx_size":               txSizeBytes,
			"fee":                   txFee,
		})
	}

	if (len(txHex) / 2) > txSizeBytes {
		panic(errors.Errorf("TX size calculation is broken: precalculated (%d), real TX size (%d).",
			txSizeBytes, len(txHex)/2))
	}

	return txHex, nil
}

// SignTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) SignTX(txHex string) (string, error) {
	return h.btcClient.SignAllTXInputs(txHex, h.config.HotWalletScriptPubKey, h.config.HotWalletRedeemScript, h.config.PrivateKey)
}

// SendTX is implementation of OffchainHelper interface from package withdraw.
func (h CommonDashHelper) SendTX(ctx context.Context, txHex string) (txHash string, err error) {
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

func (h CommonDashHelper) GetHash(txHex string) (string, error) {
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
