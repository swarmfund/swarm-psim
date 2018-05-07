package contractfunnel

import (
	"context"

	"crypto/ecdsa"

	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type ETHClient bind.ContractBackend

type Service struct {
	log    *logan.Entry
	config Config

	ethClient  ETHClient
	privateKey *ecdsa.PrivateKey

	erc20Contract *ERC20
	contracts     map[common.Address]*Contract
}

func NewService(
	log *logan.Entry,
	config Config,
	ethClient ETHClient,
	//ethWallet ETHWallet,
) (*Service, error) {

	privKey, err := crypto.HexToECDSA(config.ETHPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create PrivateKey from eth_private_key string from config")
	}

	erc20Contract, err := NewERC20(config.TokenToFunnelContractAddress, ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new ERC20Contract")
	}

	contracts := make(map[common.Address]*Contract)
	for _, addr := range config.ContractsAddresses {
		contract, err := NewContract(addr, ethClient)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create new Contract")
		}

		contracts[addr] = contract
	}

	return &Service{
		log:    log,
		config: config,

		ethClient:  ethClient,
		privateKey: privKey,

		erc20Contract: erc20Contract,
		contracts:     contracts,
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Started.")

	if s.config.OnlyViewBalances {
		s.printBalancesReport()
		return
	}

	running.WithBackOff(ctx, s.log, "all_contracts_funnel_iteration", func(ctx context.Context) error {
		s.log.Info("Started iteration of funnelling tokens from all the Contracts.")

		for addr, contract := range s.contracts {
			if running.IsCancelled(ctx) {
				return nil
			}

			_, err := s.withdrawAllTokens(contract, addr)
			if err != nil {
				return errors.Wrap(err, "Failed to withdraw all tokens", logan.F{
					"contract_address": addr,
				})
			}
		}

		return nil
	}, s.config.FunnelPeriod, 10*time.Second, time.Hour)

	s.log.Info("Received ctx cancel - service stopped cleanly.")
}

// WithdrawAllTokens can return empty txHash with nil error, if no tokens to funnel (thus no tx)
func (s *Service) withdrawAllTokens(contract *Contract, contractAddress common.Address) (txHash string, err error) {
	logger := s.log.WithField("contract_address", contractAddress.String())

	contractBalance, err := s.erc20Contract.BalanceOf(nil, contractAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get balance of the Contract")
	}

	if contractBalance.Cmp(big.NewInt(0)) == 0 {
		logger.Info("No tokens to funnel from this Contract.")
		return "", nil
	}

	tx, err := contract.WithdrawAllTokens(bind.NewKeyedTransactor(s.privateKey), s.config.TokenReceiverAddress, s.config.TokenToFunnelContractAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to Withdraw all tokens")
	}

	logger.WithFields(logan.F{
		"tx_hash":          txHash,
		"funnelled_amount": contractBalance.String(),
	}).Info("Funneled tokens from the Contract successfully.")
	return tx.Hash().String(), nil
}
