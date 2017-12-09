package txhandler

import (
	"fmt"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ledgerChangeHandler func(change xdr.LedgerEntryData) error

type ledgerChangesHandler struct {
	statable       Statable
	createHandlers map[xdr.LedgerEntryType]ledgerChangeHandler
	updateHandlers map[xdr.LedgerEntryType]ledgerChangeHandler

	log *logan.Entry
}

func newLedgerChangesHandler(statable Statable, log *logan.Entry) *ledgerChangesHandler {
	h := &ledgerChangesHandler{
		statable: statable,
		log:      log,
	}

	h.createHandlers = map[xdr.LedgerEntryType]ledgerChangeHandler{
		xdr.LedgerEntryTypeBalance: h.processBalanceCreated,
		xdr.LedgerEntryTypeAccount: h.processAccountCreated,
	}

	h.updateHandlers = map[xdr.LedgerEntryType]ledgerChangeHandler{
		xdr.LedgerEntryTypeBalance: h.processBalanceUpdated,
		xdr.LedgerEntryTypeAsset:   h.processAssetUpdated,
	}

	return h
}

func (h *ledgerChangesHandler) Handle(tx horizon.Transaction) error {
	var txMeta xdr.TransactionMeta
	err := xdr.SafeUnmarshalBase64(tx.ResultMetaXDR, &txMeta)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal transaction meta")
	}

	// xdr is invalid if operations are not there, so no need to check if it's nil
	operationChanges := txMeta.MustOperations()
	for _, opChanges := range operationChanges {
		for _, change := range opChanges.Changes {
			err = h.handleLedgerEntryChange(change)
			if err != nil {
				return errors.Wrap(err, "failed to handle ledger entry change")
			}
		}
	}

	return nil
}

func (h *ledgerChangesHandler) handleLedgerEntryChange(change xdr.LedgerEntryChange) error {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		handler, ok := h.createHandlers[change.MustCreated().Data.Type]
		if ok {
			return handler(change.MustCreated().Data)
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		handler, ok := h.updateHandlers[change.MustUpdated().Data.Type]
		if ok {
			return handler(change.MustUpdated().Data)
		}
	}

	return nil
}

func (h *ledgerChangesHandler) processBalanceCreated(entry xdr.LedgerEntryData) error {
	data := entry.MustBalance()
	accountID := state.AccountID(data.AccountId.Address())
	asset := string(data.Asset)
	balance := state.Balance{
		Account:    accountID,
		Address:    state.BalanceID(data.BalanceId.AsString()),
		Asset:      state.AssetCode(asset),
		ExchangeID: state.AccountID(data.Exchange.Address()),
	}

	if h.statable.IsSpecialAccount(accountID) {
		err := h.statable.GetSpecialAccount(accountID).AddBalance(balance)
		if err != nil {
			// we ignore error here on purpose, as special accounts balances are initialized on start based on current state of the db
			h.log.WithField("account_id", accountID).WithField("balance_id", balance.Address).WithError(err).Warn("Failed to add balance for special account - ignoring error")
		}

		return nil
	}

	err := h.statable.GetAccount(accountID).AddBalance(balance)
	if err != nil {
		// we ignore error here on purpose, as operational account balances are initialized on start based on current state of the db
		if accountID != h.statable.GetOperationalAccount() {
			return errors.Wrap(err, fmt.Sprintf("failed to add balance for account %s", accountID))
		}
	}

	return nil
}

func (h *ledgerChangesHandler) processAccountCreated(entry xdr.LedgerEntryData) error {
	data := entry.MustAccount()
	// since account might refer someone in the future
	// we need to store *all* accounts info
	account := state.Account{
		Address: state.AccountID(data.AccountId.Address()),
	}

	if data.Referrer != nil {
		account.ShareForReferrer = int64(data.ShareForReferrer)
		account.Parent = state.AccountID(data.Referrer.Address())
	}

	h.statable.AddAccount(account)
	return nil
}

func (h *ledgerChangesHandler) processBalanceUpdated(entry xdr.LedgerEntryData) error {
	data := entry.MustBalance()
	accountID := state.AccountID(data.AccountId.Address())
	// for now we are not interested in special accounts balance changes
	if h.statable.IsSpecialAccount(accountID) {
		return nil
	}

	balanceID := state.BalanceID(data.BalanceId.AsString())
	balance := h.statable.GetAccount(accountID).MustGetBalanceForBalanceID(balanceID)
	balance.Amount = int64(data.Amount)
	balance.SetFeesPaid(int64(data.FeesPaid))
	return nil
}

func (h *ledgerChangesHandler) processAssetUpdated(entry xdr.LedgerEntryData) error {
	data := entry.MustAsset()
	if data.Token == nil {
		return nil
	}

	h.statable.SetToken(state.AssetCode(data.Code), state.AssetCode(*data.Token))
	return nil
}
