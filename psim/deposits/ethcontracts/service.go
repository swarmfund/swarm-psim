package ethcontracts

import (
	"context"
	"math/big"

	"time"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

type ETHKeypair interface {
	Address() common.Address
	SignTX(*types.Transaction) (*types.Transaction, error)
}

type EntityCountGetter func(systemType string) (uint64, error)

type Service struct {
	log         *logan.Entry
	config      *Config
	eth         *ethclient.Client
	builder     *xdrbuild.Builder
	horizon     *horizon.Connector
	keypair     ETHKeypair
	entityCount EntityCountGetter
	deployerID  uint64
}

func NewService(
	log *logan.Entry,
	config *Config,
	builder *xdrbuild.Builder,
	horizon *horizon.Connector,
	keypair ETHKeypair,
	deployerID uint64,
	entityCount EntityCountGetter,
	eth *ethclient.Client,
) *Service {

	return &Service{
		log:         log,
		config:      config,
		builder:     builder,
		keypair:     keypair,
		deployerID:  deployerID,
		entityCount: entityCount,
		horizon:     horizon,
		eth:         eth,
	}
}

func (s *Service) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	do := func(ctx context.Context) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				// we are spending actual ether here,
				// so in case of emergency abandon the operations completely
				cancel()
				err = errors.Wrap(errors.FromPanic(rvr), "service panicked")
			}
		}()
		for _, systemType := range s.config.ExternalTypes {
			current, err := s.entityCount(systemType)
			if err != nil {
				return errors.Wrap(err, "failed to get current entity count")
			}

			for current <= s.config.TargetCount {
				fmt.Println(current, s.config.TargetCount)
				if running.IsCancelled(ctx) {
					return nil
				}
				fields := logan.F{}
				contract, err := s.deployContract()
				if err != nil {
					return errors.Wrap(err, "failed to deploy contract")
				}
				fields["contract"] = contract.Hex()
				s.log.WithFields(fields).Info("contract deployed")
				// critical section. contract has been deployed, we need to create entity at any cost
				running.UntilSuccess(context.Background(), s.log, "create-pool-entity", func(i context.Context) (bool, error) {
					if err := s.createPoolEntities(contract.Hex()); err != nil {
						return false, err
					}
					return true, nil
				}, 1*time.Second, 1*time.Minute)

				s.log.WithFields(fields).Info("entities created")

				current += 1
			}
		}
		s.log.Info("all good")
		return nil
	}

	running.WithBackOff(ctx, s.log, "deploy-iteration", do, 10*time.Second, 10*time.Second, 1*time.Hour)

	s.log.WithField("state", "deadlock").Error("abnormal execution")

	// FIXME will not let normal termination go
	<-make(chan struct{})
}

func (s *Service) createPoolEntities(data string) error {
	tx := s.builder.Transaction(s.config.Source)
	for _, systemType := range s.config.ExternalTypes {
		tx = tx.Op(xdrbuild.CreateExternalPoolEntry(cast.ToInt32(systemType), data, s.deployerID))
	}
	tx = tx.Sign(s.config.Signer)
	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal tx")
	}

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		return errors.Wrap(result.Err, "failed to submit tx", logan.F{
			"tx_result": result.GetLoganFields(),
		})
	}
	return nil
}

func (s *Service) deployContract() (*common.Address, error) {
	_, tx, _, err := DeployContract(&bind.TransactOpts{
		From:  s.keypair.Address(),
		Nonce: nil,
		Signer: func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return s.keypair.SignTX(tx)
		},
		Value:    big.NewInt(0),
		GasPrice: eth.FromGwei(s.config.GasPrice),
		GasLimit: eth.FromGwei(s.config.GasLimit).Uint64(),
		Context:  context.TODO(),
	}, s.eth, s.config.ContractOwner)

	if err != nil {
		return nil, errors.Wrap(err, "failed to submit contract tx")
	}

	eth.EnsureHashMined(context.Background(), s.log.WithField("tx_hash", tx.Hash().Hex()), s.eth, tx.Hash())

	receipt, err := s.eth.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
			"tx_hash": tx.Hash().String(),
		})
	}

	// TODO check transaction state/status to see if contract actually was deployed
	// TODO panic if we are not sure if contract is valid

	return &receipt.ContractAddress, nil
}
