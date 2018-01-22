package ethfunnel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet interface {
	Addresses() []common.Address
	HasAddress(common.Address) bool
	SignTX(common.Address, *types.Transaction) (*types.Transaction, error)
}
