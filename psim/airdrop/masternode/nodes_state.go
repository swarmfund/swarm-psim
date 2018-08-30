package masternode

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NodesState struct {
	eth           *ethclient.Client
	token         common.Address
	rq            RewardQueue
	furnace       common.Address
	firstBlock    time.Time
	blockDuration time.Duration
	currentBlock  uint64
}

type PayoutMeta struct {
	Block     uint64
	BlockTime time.Time
}

func (s *NodesState) Payout(ctx context.Context) (common.Address, *PayoutMeta, error) {
	blockTime := s.firstBlock.Add(time.Duration(s.currentBlock+1) * s.blockDuration)
	lastBlock := blockTime.Add(-1 * s.blockDuration)

	ethlowerblock, err := s.blockByTime(ctx, lastBlock)
	if err != nil {
		panic(err)
	}
	ethupperblock, err := s.blockByTime(ctx, blockTime)
	if err != nil {
		panic(err)
	}

	txs, err := s.interestingTXs(ctx, ethlowerblock, ethupperblock)
	if err != nil {
		panic(err)
	}

	for _, txinfo := range txs {
		tx, isPending, err := s.eth.TransactionByHash(ctx, txinfo.Hash)
		if err != nil {
			panic(err)
		}
		if isPending {
			continue
		}
		sender, err := s.eth.TransactionSender(ctx, tx, txinfo.BlockHash, txinfo.Index)
		if err != nil {
			panic(err)
		}

		s.rq.Add(sender.Hex())
	}

	// TODO go through whole reward queue and remove nodes w/o enough balance

	address := s.rq.Next()

	meta := PayoutMeta{
		Block:     s.currentBlock,
		BlockTime: blockTime,
	}
	s.currentBlock += 1
	if address == nil {
		return s.furnace, &meta, nil
	}
	return common.HexToAddress(*address), &meta, nil
}

func (s *NodesState) blockByTime(ctx context.Context, ts time.Time) (uint64, error) {
	height, err := s.eth.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	// check if highest block reached desired timestamp
	if !ts.Before(time.Unix(height.Time().Int64(), 0)) {
		// wait until we are there
		ticker := time.NewTicker(5 * time.Second)
		cursor := height.NumberU64()
		defer ticker.Stop()
		for ; ; <-ticker.C {
			block, err := s.eth.BlockByNumber(ctx, big.NewInt(int64(cursor+1)))
			if err == ethereum.NotFound {
				continue
			}
			if err != nil {
				return 0, err
			}
			cursor = block.NumberU64()
			if ts.Before(time.Unix(block.Time().Int64(), 0)) {
				return block.NumberU64(), nil
			}
		}
	}
	var upper, lower uint64 = height.NumberU64(), 0
	i := 0
	for lower < upper {
		middle := (lower + upper) / 2
		block, err := s.eth.BlockByNumber(ctx, big.NewInt(int64(middle)))
		if err != nil {
			return 0, err
		}
		blockts := time.Unix(block.Time().Int64(), 0)

		if blockts.Equal(ts) {
			return block.NumberU64(), nil
		}

		if blockts.Before(ts) {
			lower = middle + 1
		}
		if blockts.After(ts) {
			upper = middle
		}
		i += 1
	}
	return lower, nil
}

func (s *NodesState) sieveBlock(ctx context.Context, blockNumber uint64) (txs []common.Hash, err error) {
	panic("not implemented")
}

type TXInfo struct {
	Hash      common.Hash
	BlockHash common.Hash
	Index     uint
}

func (s *NodesState) interestingTXs(ctx context.Context, from, to uint64) (txs []TXInfo, err error) {
	logs, err := s.eth.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: []common.Address{s.token},
		Topics: [][]common.Hash{
			{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		// no interesting outputs
		return nil, nil
	}

	for _, log := range logs {
		if len(log.Topics) != 3 {
			// TODO log invalid log
			continue
		}
		// third indexed topic is 20 bytes receiver address packed in 40 bytes, big-endian layout
		receiver := common.BytesToAddress(log.Topics[2][len(log.Topics[2])-20:])
		if receiver != s.furnace {
			continue
		}
		txs = append(txs, TXInfo{
			Hash:      log.TxHash,
			BlockHash: log.BlockHash,
			Index:     log.TxIndex,
		})
	}
	return txs, nil
}
