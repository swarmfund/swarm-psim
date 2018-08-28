package internal

import (
	"context"
	"math"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/swarmfund/psim/psim/internal/eth"

	"strings"

	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/tokend/go/amount"
)

type ETHHelper struct {
	DepositAsset        string
	MinDepositAmount    uint64
	FixedDepositFee     uint64
	BlocksToSearchForTX uint64
	Client              *ethclient.Client
}

func (e ETHHelper) GetAsset() string {
	return e.DepositAsset
}

func (e ETHHelper) GetMinDepositAmount() uint64 {
	return e.MinDepositAmount
}

func (e ETHHelper) GetFixedDepositFee() uint64 {
	return e.FixedDepositFee
}

func (ETHHelper) ConvertToSystem(offchain uint64) uint64 {
	// FIXME copy-paste from ERC20 helper
	// expected offchain to be in gwei precision (10^9)
	var gwei int64 = 1000000000
	if offchain > math.MaxInt64 {
		panic("overflow")
	}
	result, overflow := amount.BigDivide(amount.One, int64(offchain), gwei, amount.ROUND_DOWN)
	if overflow {
		panic("overflow")
	}
	return uint64(result)
}

// BuildReference legacy bit, only direct eth transfers expected
func (ETHHelper) BuildReference(blockNumber uint64, txHash, offchainAddress string, outIndex uint, maxLen int) string {
	reference := txHash
	// yoba eth hex trimming
	if len(reference) > 64 {
		reference = reference[len(reference)-64:]
	}
	return reference
}

func (ETHHelper) GetAddressSynonyms(address string) []string {
	return []string{address, strings.ToLower(address)}
}

func (e ETHHelper) GetLastKnownBlockNumber() (uint64, error) {
	block, err := e.Client.BlockByNumber(context.TODO(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get block")
	}
	return block.NumberU64(), nil
}

func (e ETHHelper) GetBlock(number uint64) (*deposit.Block, error) {
	block, err := e.Client.BlockByNumber(context.TODO(), big.NewInt(int64(number)))

	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	result := deposit.Block{
		Hash:      block.Hash().Hex(),
		Timestamp: time.Unix(block.Time().Int64(), 0),
	}

	for _, tx := range block.Transactions() {
		result.TXs = append(result.TXs, e.transpileTX(tx))
	}

	return &result, nil
}

func (e ETHHelper) FindTX(
	ctx context.Context, blockNumber uint64, txHash string,
) (deposit.TXFindMeta, *deposit.Tx, error) {
	for i := blockNumber; i < (blockNumber + e.BlocksToSearchForTX); i++ {
		if running.IsCancelled(ctx) {
			return deposit.TXFindMeta{}, nil, nil
		}

		block, err := e.Client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err == ethereum.NotFound {
			return deposit.TXFindMeta{
				StopWaiting: false,
			}, nil, nil
		}
		if err != nil {
			return deposit.TXFindMeta{}, nil, errors.Wrap(err, "failed to get block", logan.F{"block_number": i})
		}

		for _, tx := range block.Transactions() {
			if tx.Hash().Hex() == txHash {
				deposittx := e.transpileTX(tx)
				return deposit.TXFindMeta{
					BlockWhereFound: i,
					BlockTime:       time.Unix(block.Time().Int64(), 0),
				}, &deposittx, nil
			}
		}
	}

	return deposit.TXFindMeta{
		StopWaiting: true,
	}, nil, nil
}

func (e ETHHelper) transpileTX(tx *types.Transaction) deposit.Tx {
	result := deposit.Tx{
		Hash: tx.Hash().Hex(),
	}
	if tx.To() == nil {
		return result
	}
	value := tx.Value()
	if value == nil {
		return result
	}
	value = eth.ToGwei(value)
	if !value.IsUint64() {
		panic("overflow")
	}
	result.Outs = []deposit.Out{
		{Address: tx.To().Hex(), Value: value.Uint64()},
	}
	return result
}
