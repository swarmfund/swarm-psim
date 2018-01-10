package bitcoin

import (
	"encoding/hex"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Client uses Connector to request some Bitcoin Node
// and transforms raw responses to gocoin structures.
type Client struct {
	connector Connector
}

// NewClient simply creates a new Client using provided Connector.
func NewClient(connector Connector) *Client {
	return &Client{
		connector: connector,
	}
}

// GetBlockCount returns index of the last known Block.
func (c Client) GetBlockCount() (uint64, error) {
	return c.connector.GetBlockCount()
}

// GetBlock gets Block hash by index via Connector,
// gets raw Block(hex) by hash from Connector
// and tries to parse raw Block into gocoin Block structure.
func (c Client) GetBlock(blockIndex uint64) (*btc.Block, error) {
	hash, err := c.connector.GetBlockHash(blockIndex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block hash")
	}

	return c.GetBlockByHash(hash)
}

func (c Client) GetBlockByHash(blockHash string) (*btc.Block, error) {
	blockHex, err := c.connector.GetBlock(blockHash)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block", logan.F{"block_hash": blockHash})
	}

	block, err := c.parseBlock(blockHex)
	if err != nil {
		hexToLog := blockHex[:10] + ".." + blockHex[len(blockHex)-10:]
		return nil, errors.Wrap(err, "FAiled to parse Block hex", logan.F{"block_hex": hexToLog})
	}

	err = block.BuildTxList()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build TX list of Block")
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

// CreateRawTX creates TX, which pays provided amount
// to the provided goalAddress, passing change to the provided changeAddress.
// Node decides, which UTXOs use for inputs for the TX during the FundRawTX request.
//
// The returned Transaction is not submitted into the network,
// it is not even signed yet.
// However, UTXOs used as inputs in this TX has been locked.
//
// If there is not enough unlocked BTC to fulfil the TX -
// error with cause ErrInsufficientFunds is returned.
func (c Client) CreateRawTX(goalAddress string, amount float64, changeAddress string) (resultTXHex string, err error) {
	txHex, err := c.connector.CreateRawTX(goalAddress, amount)
	if err != nil {
		return "", errors.Wrap(err, "Failed to CreateRawTX")
	}

	txHex, err = c.connector.FundRawTX(txHex, changeAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to FundRawTX", logan.F{
			"created_tx_hex": txHex,
		})
	}

	return txHex, nil
}

// SignRawTX signs provided TX with the provided privateKey.
func (c Client) SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error) {
	tx, err := c.parseTX(txHex)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse provided txHex into btc.Tx")
	}

	var outputs []Out
	for _, in := range tx.TxIn {
		outputs = append(outputs, Out{
			TXHash: hex.EncodeToString(in.Input.Hash[:]),
			Vout:   in.Input.Vout,
			ScriptPubKey: scriptPubKey,
			RedeemScript: redeemScript,
		})
	}

	return c.connector.SignRawTX(txHex, outputs, privateKey)
}

// SendRawTX submits TX into the blockchain.
func (c Client) SendRawTX(txHex string) (txHash string, err error) {
	return c.connector.SendRawTX(txHex)
}

func (c Client) parseBlock(blockHex string) (*btc.Block, error) {
	bb, err := hex.DecodeString(blockHex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode Block hex string to bytes")
	}

	block, err := btc.NewBlock(bb)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build new Block from bytes")
	}

	return block, nil
}

func (c Client) parseTX(txHex string) (*btc.Tx, error) {
	bb, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode TX hex string to bytes")
	}

	tx, _ := btc.NewTx(bb)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build new TX from bytes")
	}

	return tx, nil
}

func (c Client) IsTestnet() bool {
	return c.connector.IsTestnet()
}

func BuildCoinEmissionRequestReference(txHash string, outIndex uint32) string {
	return txHash + string(outIndex)
}
