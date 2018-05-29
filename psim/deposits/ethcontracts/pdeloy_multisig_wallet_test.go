package ethcontracts

import (
	"testing"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

func TestMultisigWalletDeploy(_ *testing.T) {
	signers := []string{
		"Ec4750c1A9E93D182af2d5004A31c891cD126Cca",
		"75517081f281f1570645Dd0162541a8b045a2652",
	}

	addr, err := s.deployMultisigWalletContract(ctx, signers, 2, 900000)
	if err != nil {
		s.log.WithError(err).Error("Failed to deploy MultisigWallet contract")
	}

	fmt.Println(addr.String())
}

// This method is not used by service, but it's needed to run once manually to deploy MultisigWallet contract.
// Signers - addresses in hex.
func deployMultisigWalletContract(
	ctx context.Context,
	signers []string,
	requiredSignatures uint8,
	gasLimit uint64) (*common.Address, error) {

	var signersAddresses []common.Address
	for _, signer := range signers {
		signersAddresses = append(signersAddresses, common.HexToAddress(signer))
	}

	_, tx, _, err := DeployMultisigWallet(&bind.TransactOpts{
		From:  s.keypair.Address(),
		Nonce: nil,
		Signer: func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return s.keypair.SignTX(tx)
		},
		Value:    big.NewInt(0),
		GasPrice: eth.FromGwei(s.config.GasPrice),
		GasLimit: gasLimit,
		Context:  ctx,
	}, s.eth, signersAddresses, requiredSignatures)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to submit contract TX")
	}
	fields := logan.F{
		"tx_hash": tx.Hash().String(),
	}

	eth.EnsureHashMined(ctx, s.log, s.eth, tx.Hash())

	receipt, err := s.eth.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get TX receipt", fields)
	}

	// TODO check transaction state/status to see if contract actually was deployed
	// TODO panic if we are not sure if contract is valid

	return &receipt.ContractAddress, nil
}
