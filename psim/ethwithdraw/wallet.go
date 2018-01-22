package ethwithdraw

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet interface {
	SignTX(common.Address, *types.Transaction) (*types.Transaction, error)
}
