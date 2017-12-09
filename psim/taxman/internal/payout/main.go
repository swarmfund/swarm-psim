package payout

import (
	"time"

	"encoding/json"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/resource"
	"gitlab.com/tokend/psim/psim/utils"
	"github.com/pkg/errors"
	"gitlab.com/tokend/psim/psim/taxman/internal/snapshoter"
)

const (
	opsPerTx = 50
)

type TransactionBuilder func() *horizon.TransactionBuilder
type Ask func(*resource.SyncRequest) (utils.AskResult, *utils.AskMeta)

type Payout struct {
	snapshot  *snapshoter.Snapshot
	txBuilder TransactionBuilder
	signer    keypair.KP
	ask       Ask
}

func New(
	snapshot *snapshoter.Snapshot, txBuilder TransactionBuilder, signer keypair.KP,
	ask Ask,
) *Payout {
	return &Payout{
		snapshot:  snapshot,
		txBuilder: txBuilder,
		signer:    signer,
		ask:       ask,
	}
}

func (p *Payout) Ensure(lock *discovery.Lock) (err error) {

	err = acquireLock(lock)
	if err != nil {
		return errors.Wrap(err, "failed to acquire lock")
	}

	defer func() {
		if err = lock.Unlock(); err != nil {
			err = errors.Wrap(err, "failed to unlock state key")
		}
	}()

	ok, err := checkPayoutSuccessful(lock)
	if err != nil {
		return errors.Wrap(err, "failed to check payout state")
	}
	if ok {
		return nil
	}

	ok, err = p.sync()
	if err != nil {
		return errors.Wrap(err, "failed to sync")
	}
	if !ok {
		return errors.New("foobar")
	}

	if err = lock.Set(resource.PayoutState{
		Successful: true,
		UpdatedAt:  time.Now().UTC(),
	}); err != nil {
		return errors.Wrap(err, "failed to set payout state")
	}

	return nil
}

func acquireLock(lock *discovery.Lock) error {
	locked, err := lock.Lock()
	if err != nil {
		return errors.Wrap(err, "failed to lock payout state")
	}
	if locked == nil {
		return errors.New("not locked")
	}
	return nil
}

func checkPayoutSuccessful(lock *discovery.Lock) (bool, error) {
	bytes, err := lock.Get()
	if err != nil {
		return false, errors.Wrap(err, "failed to get payout state")
	}
	if bytes == nil {
		return false, nil
	}
	state := resource.PayoutState{}
	if err = json.Unmarshal(bytes, &state); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal")
	}
	return state.Successful, nil
}

func (p *Payout) prepareSync() (*resource.SyncRequest, error) {
	var syncRequest resource.SyncRequest
	payoutTXs := []*horizon.TransactionBuilder{
		p.txBuilder(),
	}

	count := 0
	for _, op := range p.snapshot.SyncState {
		count += 1
		if count%opsPerTx == 0 {
			payoutTXs = append(payoutTXs, p.txBuilder())
		}

		tx := payoutTXs[len(payoutTXs)-1]
		tx.Op(horizon.PaymentOp{
			Reference:            op.Reference,
			SourceBalanceID:      string(op.SourceBalanceID),
			DestinationBalanceID: string(op.DestinationBalanceID),
			Amount:               op.Amount,
		})
	}

	for _, payoutTX := range payoutTXs {
		if payoutTX.OperationsCount() == 0 {
			continue
		}
		envelope, err := payoutTX.Sign(p.signer).Marshal64()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal tx")
		}

		syncRequest.Transactions = append(syncRequest.Transactions, *envelope)
	}

	return &syncRequest, nil
}

func (p *Payout) sync() (ok bool, err error) {
	syncRequest, err := p.prepareSync()
	if err != nil {
		return true, err
	}

	syncRequest.Ledger = p.snapshot.Ledger

	switch result, meta := p.ask(syncRequest); result {
	case utils.AskResultSuccess:
		return true, nil
	case utils.AskResultFailure:
		if meta.NeighborsAsked == 0 {
			return false, errors.New("no neighbors found")
		}
		return false, errors.New("failed to verify")
	case utils.AskResultPermanentFailure:
		return false, errors.Wrap(err, "verification failed")
	}
	panic("you should not saw this")
}
