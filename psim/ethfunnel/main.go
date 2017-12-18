package ethfunnel

import (
	"context"

	"math/big"

	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type Config struct {
	Keystore        string
	FunnelThreshold *big.Int
	Destination     string
	GasPrice        *big.Int
}

func init() {
	app.RegisterService(conf.ServiceETHFunnel, func(ctx context.Context) (utils.Service, error) {
		config := Config{
			FunnelThreshold: big.NewInt(1),
			GasPrice:        big.NewInt(1000000000),
		}
		err := figure.
			Out(&config).
			From(app.Config(ctx).Get(conf.ServiceETHFunnel)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		keystore := keystore.NewKeyStore(
			config.Keystore, keystore.LightScryptN, keystore.LightScryptP,
		)
		for _, account := range keystore.Accounts() {
			err := keystore.Unlock(account, "foobar")
			if err != nil {
				return nil, errors.Wrap(err, "failed to unlock")
			}
		}

		eth := app.Config(ctx).Ethereum()

		return NewService(ctx, app.Log(ctx), config, keystore, eth), nil
	})
}

type Service struct {
	ctx      context.Context
	keystore *keystore.KeyStore
	eth      *ethclient.Client
	log      *logan.Entry
	config   Config

	blocksCh   chan uint64
	txCh       chan Transaction
	transferCh chan Transfer
}

type Transfer struct {
	From  common.Address
	Value *big.Int
}

func NewService(ctx context.Context, log *logan.Entry, config Config, keystore *keystore.KeyStore, eth *ethclient.Client) *Service {
	return &Service{
		ctx:        ctx,
		log:        log,
		config:     config,
		keystore:   keystore,
		eth:        eth,
		blocksCh:   make(chan uint64),
		txCh:       make(chan Transaction),
		transferCh: make(chan Transfer),
	}
}

func (s *Service) Run() chan error {
	// TODO config
	confirmations := big.NewInt(1)

	go s.processBlocks()
	go s.processTXs()
	go s.processTransfers()

	go func() {
		cursor := new(big.Int).Sub(s.currentHeight(), confirmations)

		fmt.Println("looking at", cursor)

		// go through all balances before current
		for _, account := range s.keystore.Accounts() {
			fmt.Println(account.Address.String())
			balance := s.balanceAt(account.Address, cursor)

			if balance.Cmp(s.config.FunnelThreshold) == -1 {
				continue
			}

			s.transferCh <- Transfer{
				account.Address, balance,
			}
		}

		for ; ; time.Sleep(10 * time.Second) {
			head := s.currentHeight()
			fmt.Println("head is", head)
			for new(big.Int).Sub(head, confirmations).Cmp(cursor) == 1 {
				fmt.Println("adding block")
				s.blocksCh <- cursor.Uint64()
				cursor.Add(cursor, big.NewInt(1))
			}
		}
	}()
	return make(chan error)
}

func (s *Service) currentHeight() *big.Int {
	for {
		head, err := s.eth.HeaderByNumber(s.ctx, nil)
		if err != nil {
			s.log.WithError(err).Error("failed to fetch head")
			continue
		}
		return head.Number
	}
}

func (s *Service) balanceAt(account common.Address, block *big.Int) *big.Int {
	for {
		balance, err := s.eth.BalanceAt(s.ctx, account, block)
		if err != nil {
			s.log.WithError(err).Error("failed to get balance")
			continue
		}
		return balance
	}
}

func (s *Service) processBlocks() {
	for blockNumber := range s.blocksCh {
		entry := s.log.WithField("number", blockNumber)
		entry.Debug("processing block")
		block, err := s.eth.BlockByNumber(s.ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			entry.Error("failed to get block")
			s.blocksCh <- blockNumber
			continue
		}

		if block == nil {
			entry.Error("missing block")
			continue
		}

		for _, tx := range block.Transactions() {
			if tx == nil {
				continue
			}
			s.txCh <- Transaction{
				BlockNumber: block.Number(),
				TX:          *tx,
			}
		}
	}
}

type Transaction struct {
	BlockNumber *big.Int
	TX          types.Transaction
}

func (s *Service) processTransfers() {
	for transfer := range s.transferCh {
		if err := s.processTransfer(transfer); err != nil {
			s.log.WithError(err).Error("failed to process transfer")
			continue
		}
		break
	}
}

func (s *Service) processTransfer(transfer Transfer) error {
	s.log.Info("processing transfer")
	nonce, err := s.eth.NonceAt(s.ctx, transfer.From, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get nonce")
	}
	txgas := big.NewInt(21000)
	tx, err := s.keystore.SignTx(
		accounts.Account{
			Address: transfer.From,
		},
		types.NewTransaction(
			nonce,
			common.HexToAddress(s.config.Destination),
			new(big.Int).Sub(transfer.Value, new(big.Int).Mul(txgas, s.config.GasPrice)),
			txgas,
			s.config.GasPrice,
			nil),
		nil)
	fmt.Println(tx.String())
	if err != nil {
		return errors.Wrap(err, "failed to sign tx")
	}
	if err := s.eth.SendTransaction(s.ctx, tx); err != nil {
		return errors.Wrap(err, "failed to submit tx")
	}
	s.log.Info("transfer processed")
	return nil
}

func (s *Service) processTXs() {
	for tx := range s.txCh {
		for {
			if err := s.processTX(tx); err != nil {
				s.log.WithError(err).Error("failed to process tx")
				continue
			}
			break
		}
	}
}

func (s *Service) processTX(tx Transaction) (err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors.FromPanic(rvr)
		}
	}()

	// TODO also possible to monitor from and report suspicious behaviour

	to := tx.TX.To()
	if to == nil {
		return nil
	}

	if !s.keystore.HasAddress(*to) {
		return nil
	}

	balance := s.balanceAt(*to, tx.BlockNumber)
	if balance.Cmp(s.config.FunnelThreshold) == -1 {
		return nil
	}

	s.transferCh <- Transfer{
		*to, balance,
	}

	return nil
}
