package bitcoin

import (
	"encoding/hex"

	"context"

	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Client uses Connector to request some Bitcoin Node
// and transforms raw responses to btcutil structures.
type Client struct {
	connector Connector
}

// NewClient simply creates a new Client using provided Connector.
func NewClient(connector Connector) *Client {
	return &Client{
		connector: connector,
	}
}

// TODO Handle absent Block
// GetBlockCount returns the number of the last known Block.
func (c Client) GetBlockCount(ctx context.Context) (uint64, error) {
	return c.connector.GetBlockCount(ctx)
}

// GetBlock gets Block hash by provided blockNumber via Connector,
// gets raw Block(in hex) by the hash from Connector
// and tries to parse raw Block into btcutil.Block structure.
func (c Client) GetBlock(blockNumber uint64) (*btcutil.Block, error) {
	hash, err := c.connector.GetBlockHash(blockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block hash")
	}
	// TODO Handle absent Block

	block, err := c.GetBlockByHash(hash)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block by its hash", logan.F{"block_hash": hash})
	}
	// TODO Handle absent Block

	return block, nil
}

// GetBlockByHash obtains raw Block(hex) by the provided blockHash from the connector
// and parses the raw Block into btcutil.Block structure.
func (c Client) GetBlockByHash(blockHash string) (*btcutil.Block, error) {
	blockHex, err := c.connector.GetBlock(blockHash)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block from connector")
	}
	// TODO Handle absent Block

	block, err := c.parseBlock(blockHex)
	if err != nil {
		// TODO Make sure no overflow will happen
		hexToLog := blockHex[:10] + ".." + blockHex[len(blockHex)-10:]
		return nil, errors.Wrap(err, "Failed to parse Block hex", logan.F{"shortened_block_hex": hexToLog})
	}

	return block, nil
}

// TransferAllWalletMoney gets current confirmed balance of the Wallet
// and sends all those BTCs to the provided goalAddress.
func (c Client) TransferAllWalletMoney(goalAddress string) (resultTXHash string, err error) {
	balance, err := c.connector.GetBalance(false)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Wallet balance")
	}

	if balance == 0 {
		return "", nil
	}

	// Balance is not 0 - having some BTC on our Wallet, let's transfer them all.
	resultTXHash, err = c.connector.SendToAddress(goalAddress, balance)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send BTC amount to provided Address",
			logan.F{"confirmed_wallet_balance": balance})
	}

	return resultTXHash, nil
}

// GetWalletBalance returns current confirmed balance of the Wallet.
func (c Client) GetWalletBalance(includeWatchOnly bool) (float64, error) {
	balance, err := c.connector.GetBalance(includeWatchOnly)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// CreateAndFundRawTX creates TX, which pays provided amount
// to the provided goalAddress, passing change to the provided changeAddress.
// Node decides, which UTXOs to use for inputs for the TX during the FundRawTX request.
//
// The returned Transaction is not submitted into the network,
// it is not even signed yet.
// However, UTXOs used as inputs in this TX has been locked in the Node.
//
// If there is not enough unlocked BTC to fulfil the TX -
// error with cause ErrInsufficientFunds is returned.
//
// Change position in Outputs is set to 1.
//
// Provided feeRate can be nil - in this case the Wallet of the Node determines the fee.
func (c Client) CreateAndFundRawTX(goalAddress string, amount float64, changeAddress string, feeRate *float64) (resultTXHex string, err error) {
	txHex, err := c.connector.CreateRawTX(nil, map[string]float64{
		goalAddress: amount,
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to CreateAndFundRawTX")
	}

	// Fill TX with inputs - UTXOs.
	fundResult, err := c.connector.FundRawTX(txHex, changeAddress, true, feeRate)
	if err != nil {
		return "", errors.Wrap(err, "Failed to FundRawTX", logan.F{
			"created_tx_hex": txHex,
		})
	}

	return fundResult.Hex, nil
}

func (c Client) CreateRawTX(inputUTXOs []Out, addrToAmount map[string]float64) (resultTXHex string, err error) {
	return c.connector.CreateRawTX(inputUTXOs, addrToAmount)
}

func (c Client) FundRawTX(initialTXHex, changeAddress string, includeWatching bool, feeRate *float64) (result *FundResult, err error) {
	return c.connector.FundRawTX(initialTXHex, changeAddress, includeWatching, feeRate)
}

// SignAllTXInputs signs the inputs of the provided TX with the provided privateKey.
// If the provided privateKey is nil - the TX will be tried to sign by Node, using
// the private keys Node owns.
func (c Client) SignAllTXInputs(initialTXHex, scriptPubKey string, redeemScript string, privateKey string) (resultTXHex string, err error) {
	tx, err := c.parseTX(initialTXHex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse provided initialTXHex into btc.Tx")
	}

	if len(tx.MsgTx().TxIn) == 0 {
		return "", errors.New("No TX Inputs to sign")
	}

	var inputUTXOs []InputUTXO
	for _, in := range tx.MsgTx().TxIn {
		inputUTXOs = append(inputUTXOs, InputUTXO{
			Out: Out{
				TXHash: in.PreviousOutPoint.Hash.String(),
				Vout:   in.PreviousOutPoint.Index,
			},
			ScriptPubKey: scriptPubKey,
			RedeemScript: &redeemScript,
		})
	}

	return c.connector.SignRawTX(initialTXHex, inputUTXOs, []string{privateKey})
}

func (c Client) SignRawTX(initialTXHex string, inputUTXOs []InputUTXO, privateKeys []string) (resultTXHex string, err error) {
	return c.connector.SignRawTX(initialTXHex, inputUTXOs, privateKeys)
}

// SendRawTX submits TX into the blockchain.
func (c Client) SendRawTX(txHex string) (txHash string, err error) {
	return c.connector.SendRawTX(txHex)
}

// SendToAddress sends provided amount of BTC to the provided goalAddress.
// Amount in BTC.
func (c Client) SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error) {
	resultTXHash, err = c.connector.SendToAddress(goalAddress, amount)
	if err != nil {
		return "", err
	}

	return resultTXHash, nil
}

func (c Client) SendMany(addrToAmount map[string]float64) (resultTXHash string, err error) {
	resultTXHash, err = c.connector.SendMany(addrToAmount)
	if err != nil {
		return "", err
	}

	return resultTXHash, nil
}

func (c Client) parseBlock(blockHex string) (*btcutil.Block, error) {
	bb, err := hex.DecodeString(blockHex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode Block hex string to bytes")
	}

	block, err := btcutil.NewBlockFromBytes(bb)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build new Block from bytes")
	}

	return block, nil
}

func (c Client) parseTX(txHex string) (*btcutil.Tx, error) {
	bb, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode TX hex string into bytes")
	}

	tx, _ := btcutil.NewTxFromBytes(bb)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build new TX from bytes")
	}

	return tx, nil
}

func (c Client) IsTestnet() bool {
	return c.connector.IsTestnet()
}

func (c Client) GetTxUTXO(txHash string, outNumber uint32) (*UTXO, error) {
	return c.connector.GetTxUTXO(txHash, outNumber, false)
}

// EstimateFee receives blocks to be a TX included in (from 2 to 25) and
// returns which fee(in BTC) should be payed for each KB.
func (c Client) EstimateFee(blocksToBeIncluded uint) (float64, error) {
	return c.connector.EstimateFee(blocksToBeIncluded)
}

// GetAddrUTXOs returns the list of UTXOs of the provided Address.
func (c Client) GetAddrUTXOs(address string) ([]WalletUTXO, error) {
	return c.connector.ListUnspent(1, 9999999, []string{address})
}
