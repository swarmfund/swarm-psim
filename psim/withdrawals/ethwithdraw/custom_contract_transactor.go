package ethwithdraw

import (
	"github.com/ethereum/go-ethereum/common"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ContractTransactor struct {
	client bind.ContractBackend
}

func NewContractTransactor(client bind.ContractBackend) *ContractTransactor {
	return &ContractTransactor {
		client: client,
	}
}

func (t ContractTransactor) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return t.client.PendingCodeAt(ctx, account)
}

func (t ContractTransactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return t.client.PendingNonceAt(ctx, account)
}

func (t ContractTransactor) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return t.client.SuggestGasPrice(ctx)
}

func (t ContractTransactor) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return t.client.EstimateGas(ctx, call)
}
// SendTransaction is the only method that differs from bind.ContractBackend,
// it just returns nil, so that no TX is sent.
func (t ContractTransactor) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}
