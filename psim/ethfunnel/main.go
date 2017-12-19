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

		return NewService(app.Log(ctx), config, keystore, eth), nil
	})
}

type Service struct {
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

func NewService(log *logan.Entry, config Config, keystore *keystore.KeyStore, eth *ethclient.Client) *Service {
	return &Service{
		log:        log,
		config:     config,
		keystore:   keystore,
		eth:        eth,
		blocksCh:   make(chan uint64),
		txCh:       make(chan Transaction),
		transferCh: make(chan Transfer),
	}
}

func (s *Service) Run(ctx context.Context) chan error {
	// TODO config
	confirmations := big.NewInt(1)

	go s.processBlocks(ctx)
	go s.processTXs(ctx)
	go s.processTransfers(ctx)

	go func() {
		cursor := new(big.Int).Sub(s.currentHeight(ctx), confirmations)

		fmt.Println("looking at", cursor)

		// go through all balances before current
		for _, account := range s.keystore.Accounts() {
			fmt.Println(account.Address.String())
			balance := s.balanceAt(ctx, account.Address, cursor)

			if balance.Cmp(s.config.FunnelThreshold) == -1 {
				continue
			}

			s.transferCh <- Transfer{
				account.Address, balance,
			}
		}

		for ; ; time.Sleep(10 * time.Second) {
			head := s.currentHeight(ctx)
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

func (s *Service) currentHeight(ctx context.Context) *big.Int {
	for {
		head, err := s.eth.HeaderByNumber(ctx, nil)
		if err != nil {
			s.log.WithError(err).Error("failed to fetch head")
			continue
		}
		return head.Number
	}
}

func (s *Service) balanceAt(ctx context.Context, account common.Address, block *big.Int) *big.Int {
	for {
		balance, err := s.eth.BalanceAt(ctx, account, block)
		if err != nil {
			s.log.WithError(err).Error("failed to get balance")
			continue
		}
		return balance
	}
}

func (s *Service) processBlocks(ctx context.Context) {
	for blockNumber := range s.blocksCh {
		entry := s.log.WithField("number", blockNumber)
		entry.Debug("processing block")
		block, err := s.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
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

func (s *Service) processTransfers(ctx context.Context) {
	for transfer := range s.transferCh {
		if err := s.processTransfer(ctx, transfer); err != nil {
			s.log.WithError(err).Error("failed to process transfer")
			continue
		}
		break
	}
}

func (s *Service) processTransfer(ctx context.Context, transfer Transfer) error {
	s.log.Info("processing transfer")
	nonce, err := s.eth.NonceAt(ctx, transfer.From, nil)
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
	if err := s.eth.SendTransaction(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to submit tx")
	}
	s.log.Info("transfer processed")
	return nil
}

func (s *Service) processTXs(ctx context.Context) {
	// TODO Listen to ctx.
	for tx := range s.txCh {
		for {
			if err := s.processTX(ctx, tx); err != nil {
				s.log.WithError(err).Error("failed to process tx")
				continue
			}
			break
		}
	}
}

func (s *Service) processTX(ctx context.Context, tx Transaction) (err error) {
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

	balance := s.balanceAt(ctx, *to, tx.BlockNumber)
	if balance.Cmp(s.config.FunnelThreshold) == -1 {
		return nil
	}

	s.transferCh <- Transfer{
		*to, balance,
	}

	return nil
}
