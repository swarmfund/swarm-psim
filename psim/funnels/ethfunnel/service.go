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
	// Seems like constant value
	gasPerTX = uint64(21000)
)

type Service struct {
	wallet Wallet
	eth    *ethclient.Client
	log    *logan.Entry
	config Config

	blocksCh    chan *big.Int
	txCh        chan types.Transaction
	addressesCh chan common.Address
	transferCh  chan Transfer
}

type Transfer struct {
	From  common.Address
	Value *big.Int
}

func NewService(
	log *logan.Entry,
	config Config,
	wallet Wallet,
	eth *ethclient.Client,
) *Service {
	return &Service{
		log:    log,
		config: config,
		wallet: wallet,
		eth:    eth,
		// most workers will submit task back in case of error,
		// make sure channels are buffered
		blocksCh:    make(chan *big.Int, 1),
		txCh:        make(chan types.Transaction, 1),
		addressesCh: make(chan common.Address, 1),
		transferCh:  make(chan Transfer, 1),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Started.")

	go s.watchHeightStreamBlockNumbers(ctx)
	go s.processAddresses(ctx)
	go s.fetchBlocksStreamAllTXs(ctx)
	go s.streamWalletAddresses()
	go s.consumeTXs()

	<-ctx.Done()
}

func (s *Service) streamWalletAddresses() {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.log.WithError(errors.FromPanic(rvr)).Error("backlog panicked")
		}
	}()

	s.log.Info("Started streaming wallet Addresses.")

	for _, addr := range s.wallet.Addresses(context.TODO()) {
		s.addressesCh <- addr
	}

	s.log.Info("Finished streaming wallet Addresses.")
}

func (s *Service) watchHeightStreamBlockNumbers(ctx context.Context) *big.Int {
	do := func(cursor *big.Int) (_ *big.Int, err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()

		head, err := s.eth.HeaderByNumber(ctx, nil)
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

func (s *Service) fetchBlocksStreamAllTXs(ctx context.Context) {
	do := func(number *big.Int) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.Wrap(err, "Block fetch panicked")
			}
		}()

		block, err := s.eth.BlockByNumber(ctx, number)
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

	for blockNumber := range s.blocksCh {
		entry := s.log.WithField("block_number", blockNumber)

		if err := do(blockNumber); err != nil {
			entry.WithError(err).Error("Failed to fetch Block.")
			s.blocksCh <- blockNumber
			continue
		}

		entry.Debug("Block fetched successfully")
	}
}

func (s *Service) processAddresses(ctx context.Context) {
	for addr := range s.addressesCh {
		entry := s.log.WithField("addr", addr.Hex())

		if err := s.processAddr(ctx, addr); err != nil {
			entry.WithError(err).Error("Failed to process Address with ETH, putting Address back to channel.")
			//s.addressesCh <- addr
			continue
		}
	}
}

func (s *Service) processAddr(ctx context.Context, address common.Address) (err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors.FromPanic(rvr)
		}
	}()

	s.log.WithField("addr", address.String()).Debug("Processing Address.")

	balance, err := s.eth.BalanceAt(ctx, address, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to get ETH balance of the Address from node")
	}

	if balance.Cmp(s.config.Threshold) == -1 {
		if balance.Cmp(big.NewInt(0)) == 1 {
			// Balance is bigger than 0
			s.log.WithFields(logan.F{
				"addr":      address.String(),
				"balance":   balance,
				"threshold": s.config.Threshold,
			}).Debug("Balance on ETH Address is more than 0, but less the threshold - skipping.")
		}
		return nil
	}

	nonce, err := s.eth.NonceAt(ctx, address, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to get nonce of the Address form node")
	}

	tx, err := s.wallet.SignTX(
		address,
		types.NewTransaction(
			nonce,
			s.config.Destination,
			new(big.Int).Sub(balance, new(big.Int).Mul(big.NewInt(int64(gasPerTX)), s.config.GasPrice)),
			gasPerTX,
			s.config.GasPrice,
			nil),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to sign ETH TX")
	}

	if err := s.eth.SendTransaction(ctx, tx); err != nil {
		return errors.Wrap(err, "Failed to submit TX")
	}

	s.log.WithFields(logan.F{
		"tx_hash": tx.Hash().String(),
		"addr":    address.String(),
	}).Info("TX submitted.")

	return nil
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
			s.log.WithFields(logan.F{
				"tx_hash": tx.Hash().String(),
				"addr":    to.Hex(),
			}).Info("Found TX to managed Address, streaming Address into channel.")

			s.addressesCh <- *to
		}

		return nil
	}

	for tx := range s.txCh {
		entry := s.log.WithField("tx_hash", tx.Hash().String())

		if err := do(tx); err != nil {
			entry.WithError(err).Error("Failed to process tx, putting the TX back to channel.")
			// FIXME This can lock the execution
			s.txCh <- tx
			continue
		}
	}
}
