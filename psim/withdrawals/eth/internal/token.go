package internal

import (
	"math/big"

	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Token struct {
	abi     abi.ABI
	address common.Address
}

func NewToken(address common.Address, tokenABI string) (*Token, error) {
	parsed, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse abi")
	}
	return &Token{
		address: address,
		abi:     parsed,
	}, nil
}

func (t *Token) Transfer(from common.Address, amount *big.Int) ([]byte, error) {
	return t.abi.Pack("transfer", []interface{}{from, amount}...)
}

func (t *Token) Address() common.Address {
	return t.address
}
