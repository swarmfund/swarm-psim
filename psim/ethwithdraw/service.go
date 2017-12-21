package ethwithdraw

import (
	"context"
	"fmt"
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
)

type Service struct {
	ctx        context.Context
	log        *logan.Entry
	horizonV2  *horizonV2.Connector
	horizon    *horizon.Connector
	config     Config
	keystore   *keystore.KeyStore
	eth        *ethclient.Client
	withdrawCh chan Withdraw
}

func NewService(
	log *logan.Entry, config Config, horizonV2 *horizonV2.Connector,
	keystore *keystore.KeyStore, eth *ethclient.Client, horizon *horizon.Connector,
) *Service {
	return &Service{
		log:       log,
		config:    config,
		horizonV2: horizonV2,
		horizon:   horizon,
		keystore:  keystore,
		eth:       eth,
		// make sure channel is buffered
		withdrawCh: make(chan Withdraw, 1),
	}
}

func (s *Service) listenRequests() {
	requestCh := make(chan horizonV2.Request)
	errs := s.horizonV2.Listener().Requests(requestCh)
	for {
		select {
		case request := <-requestCh:
			s.queueRequest(request)
		case err := <-errs:
			s.log.WithError(err).Error("failed to get request")
		}
	}
}

func (s *Service) queueRequest(request horizonV2.Request) {
	entry := s.log.WithField("request", request.ID)
	if request.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		entry.Debug("not a withdraw, skipping")
		return
	}

	if request.Details.Withdraw.DestinationAsset != s.config.Asset {
		entry.Debug("different asset, skipping")
		return
	}

	if request.State != requestStatePending {
		entry.Debug("not pending, skipping")
		return
	}

	// request passes filters
	// let's get to processing
	s.withdrawCh <- Withdraw{
		Request: request,
	}
}

type Withdraw struct {
	Request  horizonV2.Request
	ETH      *types.Transaction
	Approved bool
}

func (s *Service) processRequests() {
	do := func(log *logan.Entry, withdraw *Withdraw) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()

		// craft ethereum transaction
		withdraw.ETH, err = s.craftETH(withdraw)
		if err != nil {
			return errors.Wrap(err, "failed craft eth tx")
		}

		// submit stellar tx
		if !withdraw.Approved {
			if err := s.approveWithdraw(withdraw); err != nil {
				return errors.Wrap(err, "failed to approve request")
			}
			log.Info("request approved")
			withdraw.Approved = true
		}

		// submit eth tx
		s.submitETH(withdraw.ETH)

		// wait while tx is mined
		s.ensureMined(withdraw.ETH.Hash())

		return nil
	}

	for withdraw := range s.withdrawCh {
		entry := s.log.WithField("request", withdraw.Request.ID)
		if err := do(entry, &withdraw); err != nil {
			entry.WithError(err).Error("processing failed")
			go func() { s.withdrawCh <- withdraw }()
			time.Sleep(10 * time.Second)
			continue
		}
		entry.Info("processed")
	}
}

func (s *Service) ensureMined(hash common.Hash) {
	do := func(hash common.Hash) (ok bool, err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()
		tx, pending, err := s.eth.TransactionByHash(s.ctx, hash)
		if err != nil {
			return false, errors.Wrap(err, "failed to get tx")
		}
		return tx != nil && !pending, nil
	}
	entry := s.log.WithField("hash", hash.Hex())
	for ; ; time.Sleep(10 * time.Second) {
		ok, err := do(hash)
		if err != nil {
			entry.WithError(err).Error("failed to get tx")
			continue
		}
		if ok {
			entry.Info("mined")
			return
		}
	}
}

func (s *Service) submitETH(tx *types.Transaction) {
	do := func(tx *types.Transaction) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()
		return s.eth.SendTransaction(s.ctx, tx)
	}
	entry := s.log.WithField("hash", tx.Hash().Hex())
	for {
		if err := do(tx); err != nil {
			entry.WithError(err).Error("failed to submit eth tx")
			// TODO incremental backoff
			time.Sleep(10 * time.Second)
			continue
		}
		entry.Info("submitted eth tx")
		return
	}
}

func (s *Service) craftETH(withdraw *Withdraw) (*types.Transaction, error) {
	txFee := new(big.Int).Mul(txGas, s.config.GasPrice)

	nonce, err := s.eth.PendingNonceAt(s.ctx, s.config.From)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	value := new(big.Int).Sub(
		new(big.Int).Mul(
			big.NewInt(int64(withdraw.Request.Details.Withdraw.Amount)),
			ethPrecision,
		),
		txFee,
	)

	tx, err := s.keystore.SignTx(
		accounts.Account{
			Address: s.config.From,
		},
		types.NewTransaction(
			nonce,
			common.HexToAddress(withdraw.Request.Details.Withdraw.ExternalDetails),
			value,
			txGas,
			s.config.GasPrice,
			nil,
		),
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign tx")
	}

	return tx, nil
}

func (s *Service) approveWithdraw(withdraw *Withdraw) error {
	txraw, err := rlp.EncodeToBytes(withdraw.ETH.Data())
	if err != nil {
		return errors.Wrap(err, "failed to encode eth tx")
	}

	err = s.horizon.Transaction(&horizon.TransactionBuilder{
		Source: s.config.Source,
	}).Op(&horizon.ReviewRequestOp{
		ID:     withdraw.Request.ID,
		Hash:   withdraw.Request.Hash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: horizon.ReviewRequestOpDetails{
			Type: xdr.ReviewableRequestTypeWithdraw,
			Withdrawal: &horizon.ReviewRequestOpWithdrawalDetails{
				ExternalDetails: fmt.Sprintf("%x", txraw),
			},
		},
	}).
		Sign(s.config.Signer).
		Submit()
	if err != nil {
		return errors.Wrap(err, "failed to submit tx")
	}
	return nil
}

func (s *Service) Run(ctx context.Context) chan error {
	s.ctx = ctx

	// TODO check there is no pending transactions in the pool
	// TODO check all approved withdraw requests are really approved

	go s.listenRequests()
	go s.processRequests()

	return make(chan error)
}
