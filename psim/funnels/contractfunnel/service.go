package contractfunnel

import (
	"context"

	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

type ETHClient bind.ContractBackend

type Opts struct {
	Log             *logan.Entry
	Config          Config
	Eth             *ethclient.Client
	Keypair         *eth.Keypair
	AddressProvider *addrstate.Watcher
	ExternalSystems []int32
	Tokens          []common.Address
	HotWallet       common.Address
	Threshold       *big.Int
	GasPrice        *big.Int
}

type Service struct {
	Opts
	// instantiated ERC20 tokens
	tokens []ERC20
	// instantiated deposit contracts
	contracts map[string]Contract
}

func (s *Service) initTokens() error {
	for _, address := range s.Tokens {
		token, err := NewERC20(address, s.Eth)
		if err != nil {
			return errors.Wrap(err, "failed to init token contract", logan.F{
				"token_address": address,
			})
		}
		s.tokens = append(s.tokens, *token)
	}
	return nil
}

// return deposit contract instance by address, doing some checks.
// not safe for concurrent use
func (s *Service) getContract(address common.Address) (*Contract, error) {
	if s.contracts == nil {
		s.contracts = map[string]Contract{}
	}

	if contract, ok := s.contracts[address.Hex()]; ok {
		return &contract, nil
	}

	contract, err := NewContract(address, s.Eth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init contract", logan.F{
			"contract_addr": address.Hex(),
		})
	}

	s.contracts[address.Hex()] = *contract

	return contract, nil
}

func (s *Service) isOwner(contract *Contract) (bool, error) {
	owner, err := contract.Owner(nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to get contract owner")
	}
	return owner.Hex() == s.Keypair.Address().Hex(), nil
}

func (s *Service) Run(ctx context.Context) {
	if err := s.initTokens(); err != nil {
		panic(errors.Wrap(err, "failed to init tokens"))
	}
	do := func(ctx context.Context) error {
		for _, system := range s.ExternalSystems {
			fields := logan.F{
				"external_system": system,
				"eth_signer":      s.Keypair.Address().Hex(),
			}
			entities := s.AddressProvider.BindedExternalSystemEntities(ctx, system)
			for _, entity := range entities {
				address := common.HexToAddress(entity)
				fields["contract_address"] = address.Hex()
				contract, err := s.getContract(address)
				if err != nil {
					return errors.Wrap(err, "failed to get contract", fields)
				}
				ok, err := s.isOwner(contract)
				if err != nil {
					return errors.Wrap(err, "failed to check owner", fields)
				}
				if !ok {
					s.Log.WithFields(fields).Warn("not an owner")
					continue
				}
				for i, token := range s.tokens {
					tokenaddr := s.Tokens[i]
					fields["tokend_addr"] = tokenaddr.Hex()
					balance, err := token.BalanceOf(nil, address)
					if err != nil {
						return errors.Wrap(err, "failed to get balance of", fields)
					}
					fields["balance"] = balance.String()
					if balance.Cmp(eth.FromGwei(s.Threshold)) == -1 {
						s.Log.WithFields(fields).Info("lower than threshold")
						continue
					}
					tx, err := contract.WithdrawAllTokens(&bind.TransactOpts{
						GasPrice: eth.FromGwei(s.GasPrice),
						GasLimit: 200000,
						Signer: func(signer types.Signer, addresses common.Address, transaction *types.Transaction) (*types.Transaction, error) {
							return s.Keypair.SignTX(transaction)
						},
					}, s.HotWallet, tokenaddr)
					if err != nil {
						return errors.Wrap(err, "failed to withdraw token", fields)
					}
					fields["tx_hash"] = tx.Hash()

					eth.EnsureHashMined(ctx, s.Log.WithFields(fields), s.Eth, tx.Hash())
				}
			}
		}
		return nil
	}

	running.WithBackOff(ctx, s.Log, "funnel-iteration", do, 10*time.Second, 10*time.Second, 1*time.Hour)
}
