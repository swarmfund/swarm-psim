package ethfunnel

import (
	"context"
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	txgas = big.NewInt(21000)
)

type Service struct {
	ctx    context.Context
	wallet Wallet
	eth    *ethclient.Client
	log    *logan.Entry
	config Config

	blocksCh   chan *big.Int
	txCh       chan types.Transaction
	accountCh  chan common.Address
	transferCh chan Transfer
}

type Transfer struct {
	From  common.Address
	Value *big.Int
}

func NewService(
	ctx context.Context, log *logan.Entry, config Config, wallet Wallet, eth *ethclient.Client,
) *Service {
	return &Service{
		ctx:    ctx,
		log:    log,
		config: config,
		wallet: wallet,
		eth:    eth,
		// most workers will submit task back in case of error,
		// make sure channels are buffered
		blocksCh:   make(chan *big.Int, 1),
		txCh:       make(chan types.Transaction, 1),
		accountCh:  make(chan common.Address, 1),
		transferCh: make(chan Transfer, 1),
	}
}

func (s *Service) Run(ctx context.Context) {
	go s.watchHeight()
	go s.processAccounts()
	go s.fetchBlocks()
	go s.accountBacklog()
	go s.consumeTXs()
}

func (s *Service) accountBacklog() {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.log.WithError(errors.FromPanic(rvr)).Error("backlog panicked")
		}
	}()

	s.log.Debug("processing backlog")
	for _, addr := range s.wallet.Addresses() {
		s.accountCh <- addr
	}
	s.log.Info("backlog processed")
}

func (s *Service) watchHeight() *big.Int {
	do := func(cursor *big.Int) (_ *big.Int, err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()

		head, err := s.eth.HeaderByNumber(s.ctx, nil)
		if err != nil {
			return cursor, errors.Wrap(err, "failed to get head")
		}

		if cursor == nil {
			cursor = new(big.Int).Sub(head.Number, s.config.Confirmations)
		}

		for new(big.Int).Sub(head.Number, s.config.Confirmations).Cmp(cursor) == 1 {
			s.blocksCh <- new(big.Int).Set(cursor)
			cursor.Add(cursor, big.NewInt(1))
		}

		return cursor, nil
	}

	var cursor *big.Int
	for ; ; time.Sleep(10 * time.Second) {
		cursorUpdate, err := do(cursor)
		if err != nil {
			s.log.WithError(err).Error("failed to update height")
			continue
		}
		cursor = cursorUpdate
	}

}

func (s *Service) fetchBlocks() {
	do := func(number *big.Int) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.Wrap(err, "block fetch panicked")
			}
		}()
		block, err := s.eth.BlockByNumber(s.ctx, number)
		if err != nil {
			return errors.Wrap(err, "failed to get block")
		}

		if block == nil {
			return errors.New("block missing")
		}

		for _, tx := range block.Transactions() {
			if tx == nil {
				continue
			}
			s.txCh <- *tx
		}
		return nil
	}

	for number := range s.blocksCh {
		entry := s.log.WithField("number", number)
		entry.Debug("fetching block")
		if err := do(number); err != nil {
			s.log.WithError(err).Error("failed to fetch block")
			s.blocksCh <- number
		}
		entry.Debug("block fetched")
	}
}

func (s *Service) processAccounts() {
	do := func(entry *logan.Entry, address common.Address) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()

		balance, err := s.eth.BalanceAt(s.ctx, address, nil)
		if err != nil {
			return errors.Wrap(err, "failed to get balance")
		}

		if balance.Cmp(s.config.Threshold) == -1 {
			entry.WithFields(logan.F{
				"balance":   balance,
				"threshold": s.config.Threshold,
			}).Debug("skipping account")
			return nil
		}

		nonce, err := s.eth.NonceAt(s.ctx, address, nil)
		if err != nil {
			return errors.Wrap(err, "failed to get nonce")
		}

		tx, err := s.wallet.SignTX(
			address,
			types.NewTransaction(
				nonce,
				s.config.Destination,
				new(big.Int).Sub(balance, new(big.Int).Mul(txgas, s.config.GasPrice)),
				txgas,
				s.config.GasPrice,
				nil),
		)
		if err != nil {
			return errors.Wrap(err, "failed to sign tx")
		}

		if err := s.eth.SendTransaction(s.ctx, tx); err != nil {
			return errors.Wrap(err, "failed to submit tx")
		}

		entry.WithField("hash", tx.Hash().Hex()).Info("tx submitted")

		return nil
	}

	for account := range s.accountCh {
		entry := s.log.WithField("account", account.Hex())
		entry.Debug("processing account")
		if err := do(entry, account); err != nil {
			entry.WithError(err).Error("failed to process account")
			s.accountCh <- account
			continue
		}
		s.log.Debug("account processed")
	}
}

func (s *Service) consumeTXs() {
	do := func(tx types.Transaction) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()

		to := tx.To()
		if to != nil && s.wallet.HasAddress(*to) {
			s.log.WithField("addr", to.Hex()).Info("found managed address")
			s.accountCh <- *to
		}

		return nil
	}

	for tx := range s.txCh {
		entry := s.log.WithField("hash", tx.Hash().Hex())
		entry.Debug("processing tx")
		if err := do(tx); err != nil {
			s.log.WithError(err).Error("failed to process tx")
			s.txCh <- tx
			continue
		}
		entry.Debug("tx processed")
	}
}
