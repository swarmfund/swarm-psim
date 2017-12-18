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
		return nil, errors.Wrap(err, "Failed to get Block", logan.Field("block_hash", blockHash))
	}

	block, err := c.parseBlock(blockHex)
	if err != nil {
		hexToLog := blockHex[:10] + ".." + blockHex[len(blockHex)-10:]
		return nil, errors.Wrap(err, "FAiled to parse Block hex", logan.Field("block_hex", hexToLog))
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
	balance, err := c.connector.GetBalance()
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
			logan.Field("confirmed_wallet_balance", balance))
	}

	return resultTXHash, nil
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

func (c Client) IsTestnet() bool {
	return c.connector.IsTestnet()
}

func BuildCoinEmissionRequestReference(txHash string, outIndex int) string {
	return txHash + string(outIndex)
}
