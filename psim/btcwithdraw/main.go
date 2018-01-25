package btcwithdraw

import (
	"context"

	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/psim/withdraw"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdraw, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceBTCWithdraw)).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCWithdraw,
		})
	}

	horizonConnector := globalConfig.Horizon()

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return withdraw.New(
		conf.ServiceBTCWithdraw,
		conf.ServiceBTCWithdrawVerify,
		config.SignerKP,
		log,
		horizonConnector.Listener(),
		horizonConnector,
		builder,
		globalConfig.Discovery(),
		New(config, globalConfig.Bitcoin()),
	), nil
}

// BTCClient is interface to be implemented by Bitcoin Core client
// to parametrise the Service.
type BTCClient interface {
	// CreateAndFundRawTX sets Change position in Outputs to 1.
	CreateAndFundRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error)
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	GetNetParams() *chaincfg.Params
}

// BTCHelper is BTC specific implementation of the OffchainHelper interface from package withdraw.
type BTCHelper struct {
	config    Config
	btcClient BTCClient
}

// New is constructor for BTCHelper.
func New(config Config, btcClient BTCClient) *BTCHelper {
	return &BTCHelper{
		config:    config,
		btcClient: btcClient,
	}
}

// TODO Config
// GetAsset is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) GetAsset() string {
	// TODO Config
	return "BTC"
}

// TODO Remove me.
// GetHotWallerAddress is implementation of OffchainHelper interface from package withdraw.
//func (h BTCHelper) GetHotWallerAddress() string {
//	return h.config.HotWalletAddress
//}

// GetMinWithdrawAmount is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) GetMinWithdrawAmount() int64 {
	return h.config.MinWithdrawAmount
}

// ValidateAddress is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) ValidateAddress(addr string) error {
	_, err := btcutil.DecodeAddress(addr, h.btcClient.GetNetParams())
	return err
}

// ValidateTx is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) ValidateTx(txHex string, withdrawAddress string, withdrawAmount int64) (string, error) {
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
	txOutAddresses, err := getBtcTXAddresses(tx, h.btcClient.GetNetParams())
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

// GetWithdrawAmount is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) ConvertAmount(destinationAmount int64) int64 {
	return destinationAmount * ((10^8) / amount.One)
}

// CreateTX is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) CreateTX(addr string, amount int64) (txHex string, err error) {
	txHex, err = h.btcClient.CreateAndFundRawTX(addr, float64(amount), h.config.HotWalletAddress)
	if err != nil {
		if errors.Cause(err) == bitcoin.ErrInsufficientFunds {
			return "", errors.Wrap(err, "Could not create raw TX - not enough BTC on hot wallet")
		}

		return "", errors.Wrap(err, "Failed to create raw TX")
	}

	return txHex, nil
}

// SignTX is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) SignTX(txHex string) (string, error) {
	return h.btcClient.SignAllTXInputs(txHex, h.config.HotWalletScriptPubKey, h.config.HotWalletRedeemScript, h.config.PrivateKey)
}

// SendTX is implementation of OffchainHelper interface from package withdraw.
func (h BTCHelper) SendTX(txHex string) (txHash string, err error) {
	return h.btcClient.SendRawTX(txHex)
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
