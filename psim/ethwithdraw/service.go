package ethwithdraw

import (
	"context"
	"math/big"

	"time"

	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
)

type Service struct {
	log        *logan.Entry
	horizon    *horizon.Connector
	config     Config
	eth        *ethclient.Client
	withdrawCh chan Withdraw

	address common.Address
	wallet  Wallet
}

func NewService(
	log *logan.Entry, config Config, horizon *horizon.Connector,
	wallet Wallet, eth *ethclient.Client, address common.Address,
) *Service {
	return &Service{
		log:     log,
		config:  config,
		horizon: horizon,
		wallet:  wallet,
		address: address,
		eth:     eth,
		// make sure channel is buffered
		withdrawCh: make(chan Withdraw, 1),
	}
}

func (s *Service) listenRequests() {
	requestCh := make(chan horizon.Request)
	errs := s.horizon.Listener().WithdrawalRequests(requestCh)
	for {
		select {
		case request := <-requestCh:
			s.queueRequest(request)
		case err := <-errs:
			s.log.WithError(err).Error("failed to get request")
		}
	}
}

func (s *Service) queueRequest(request horizon.Request) {
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

	// checking if request is well-formed
	// TODO reject otherwise
	if request.Details.Withdraw.ExternalDetails == nil {
		entry.Warn("missing external details, skipping")
		return
	}

	if _, ok := request.Details.Withdraw.ExternalDetails["address"]; !ok {
		entry.Warn("missing external address, skipping")
		return
	}

	// request passes filters
	// let's get to processing
	s.withdrawCh <- Withdraw{
		Request: request,
	}
}

type Withdraw struct {
	Request  horizon.Request
	ETH      *types.Transaction
	Approved bool
}

func (s *Service) processRequests(ctx context.Context) {
	do := func(log *logan.Entry, withdraw *Withdraw) (err error) {
		defer func() {
			//if rvr := recover(); rvr != nil {
			//	err = errors.FromPanic(rvr)
			//}
		}()

		// craft ethereum transaction
		withdraw.ETH, err = s.craftETH(ctx, withdraw)
		if err != nil {
			return errors.Wrap(err, "failed craft eth tx")
		}

		// submit stellar tx
		if !withdraw.Approved {
			if err := s.approveWithdraw(ctx, withdraw); err != nil {
				return errors.Wrap(err, "failed to approve request")
			}
			log.Info("request approved")
			withdraw.Approved = true
		}

		// submit eth tx
		s.submitETH(ctx, withdraw.ETH)

		// wait while tx is mined
		s.ensureMined(ctx, withdraw.ETH.Hash())

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

func (s *Service) ensureMined(ctx context.Context, hash common.Hash) {
	do := func(hash common.Hash) (ok bool, err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()
		tx, pending, err := s.eth.TransactionByHash(ctx, hash)
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

func (s *Service) submitETH(ctx context.Context, tx *types.Transaction) {
	do := func(tx *types.Transaction) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(rvr)
			}
		}()
		return s.eth.SendTransaction(ctx, tx)
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

func (s *Service) craftETH(ctx context.Context, withdraw *Withdraw) (*types.Transaction, error) {
	txFee := new(big.Int).Mul(txGas, s.config.GasPrice)

	nonce, err := s.eth.PendingNonceAt(ctx, s.address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	value := new(big.Int).Sub(
		new(big.Int).Mul(
			big.NewInt(int64(withdraw.Request.Details.Withdraw.DestinationAmount)),
			ethPrecision,
		),
		txFee,
	)

	fmt.Println(value.String())

	tx, err := s.wallet.SignTX(
		s.address,
		types.NewTransaction(
			nonce,
			// FIXME
			common.HexToAddress(withdraw.Request.Details.Withdraw.ExternalDetails["address"].(string)),
			value,
			txGas,
			s.config.GasPrice,
			nil,
		),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign tx")
	}

	return tx, nil
}

func (s *Service) approveWithdraw(ctx context.Context, withdraw *Withdraw) error {
	var builder *xdrbuild.Builder
	{
		info, err := s.horizon.Info()
		if err != nil {
			return errors.Wrap(err, "failed to get horizon info")
		}

		builder = xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	}

	envelope, err := builder.Transaction(s.config.Source).
		Op(xdrbuild.ReviewRequestOp{
			Hash: withdraw.Request.Hash,
			ID:   withdraw.Request.ID,
			Details: xdrbuild.WithdrawalDetails{
				// TODO Set Hash of the ETH TX
				ExternalDetails: "",
			},
			Action: xdr.ReviewRequestOpActionApprove,
		}).
		Sign(s.config.Signer).
		Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to build review op")
	}

	result := s.horizon.Submitter().Submit(ctx, envelope)

	if result.Err != nil {
		return errors.Wrap(result.Err, "failed to submit review op", logan.F{
			"submit_response_raw":      string(result.RawResponse),
			"submit_response_tx_code":  result.TXCode,
			"submit_response_op_codes": result.OpCodes,
		})
	}

	return nil
}

func (s *Service) Run(ctx context.Context) {
	// TODO check there is no pending transactions in the pool
	// TODO check all approved withdraw requests are really approved

	go s.listenRequests()
	go s.processRequests(ctx)

	<-ctx.Done()
}
