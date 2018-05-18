package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type ByHasher interface {
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
}

func EnsureHashMined(ctx context.Context, log running.Logger, getter ByHasher, hash common.Hash) {
	fields := logan.F{
		"tx_hash": hash.String(),
	}
	running.UntilSuccess(ctx, log, "ensure-hash-mined", func(ctx context.Context) (bool, error) {
		tx, isPending, err := getter.TransactionByHash(ctx, hash)
		if err != nil {
			return false, errors.Wrap(err, "failed to get tx", fields)
		}
		if tx == nil {
			return false, errors.From(errors.New("tx not found"), fields)
		}
		if isPending {
			return false, nil
		}
		return true, nil
	}, 5*time.Second, time.Minute)
}
