package bitcoin

import (
	"encoding/hex"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Client uses Connector to request blockchain info from some Bitcoin Node
// and transforms raw requests to gocoin structures
type Client struct {
	Connector Connector
}

// NewClient simply creates a new Client using provided Connector.
func NewClient(connector Connector) *Client {
	return &Client{
		Connector: connector,
	}
}

// GetBlockCount returns index of the last known Block.
func (c Client) GetBlockCount() (uint64, error) {
	return c.Connector.GetBlockCount()
}

// GetBlock gets Block hash by index via Connector,
// gets raw Block(hex) by hash from Connector
// and tries to parse raw Block into gocoin Block structure.
func (c Client) GetBlock(blockIndex uint64) (*btc.Block, error) {
	hash, err := c.Connector.GetBlockHash(blockIndex)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block hash")
	}

	return c.GetBlockByHash(hash)
}

func (c Client) GetBlockByHash(blockHash string) (*btc.Block, error) {
	blockHex, err := c.Connector.GetBlock(blockHash)
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
	return c.Connector.IsTestnet()
}

func BuildCoinEmissionRequestReference(txHash string, outIndex int) string {
	return txHash + string(outIndex)
}
