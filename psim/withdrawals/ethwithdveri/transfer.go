package ethwithdveri

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	Id                 *big.Int
	To                 common.Address
	Amount             *big.Int
	Token              common.Address
	TransferType       uint8
	NumberOrSignatures uint8
}

// TODO GetLoganFields
//func
