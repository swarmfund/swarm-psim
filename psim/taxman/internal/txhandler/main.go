// Package txhandler provides methods to handle blockchain event which might have effect on the state of the taxman
package txhandler

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//go:generate mockery -case underscore -name Statable

// Statable - provides methods to store changes of the taxman state
type Statable interface {
	// SetPayoutPeriod - sets payout period
	SetPayoutPeriod(d *time.Duration)
	// IsSpecialAccount - returns true if account requires special handing
	IsSpecialAccount(accountID state.AccountID) bool
	// GetAccount - returns account by accountID, panics if account not found
	GetAccount(accountID state.AccountID) *state.Account
	// GetSpecialAccount - returns special account by accountID, panics if account not found
	GetSpecialAccount(accountID state.AccountID) *state.Account
	// AddAccount - adds account. If account already exists - panic
	AddAccount(account state.Account)
	// SetToken - sets token for the asset
	SetToken(asset, token state.AssetCode)
	// SetLedger - set's current ledger
	SetLedger(ledger int64)
	// GetOperationalAccount - returns id of operational account
	GetOperationalAccount() state.AccountID
}

//go:generate mockery -case underscore -name HorizonTxHandler
// HorizonTxHandler - handles horizon tx
type HorizonTxHandler interface {
	Handle(transaction horizon.Transaction) error
}

// Handler - handles blockchain events and changes state.State accordingly
type Handler struct {
	handlers  []HorizonTxHandler
	txsToSkip map[string]bool
	statable  Statable

	log *logan.Entry
}

func NewHandler(statable Statable, txsToSkip map[string]bool, log *logan.Entry) *Handler {
	return &Handler{
		statable: statable,
		handlers: []HorizonTxHandler{
			newTxHandler(statable, log),
			newLedgerChangesHandler(statable, log),
		},
		txsToSkip: txsToSkip,
		log:       log,
	}
}

// Handle - handles the sse event. Returns identifier of the processed event
func (h *Handler) Handle(tx horizon.Transaction) error {
	log := h.log.WithField("id", tx.ID).WithField("ledger", tx.Ledger)
	log.Info("processing tx")
	if h.txsToSkip[tx.ID] {
		log.Warn("Skipping")
		return nil
	}

	for i := range h.handlers {
		err := h.handlers[i].Handle(tx)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to handle tx %s", tx.ID))
		}
	}

	h.statable.SetLedger(tx.Ledger)

	return nil
}
