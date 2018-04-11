package ethfunnel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"context"
)

type Wallet interface {
	Addresses(ctx context.Context) []common.Address
	HasAddress(common.Address) bool
	SignTX(common.Address, *types.Transaction) (*types.Transaction, error)
}
