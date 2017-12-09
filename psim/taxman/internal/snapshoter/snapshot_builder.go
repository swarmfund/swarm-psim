package snapshoter

import (
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"github.com/pkg/errors"
)

//go:generate mockery -case underscore -testonly -inpkg -name statable
type statable interface {
	// GetLedger - returns current ledger sequence number
	GetLedger() int64
	// GetTotalFeesToShare - returns total amount of fees to share in current payout, grouped by asset
	GetTotalFeesToShare() map[state.AssetCode]int64
	// GetChildren - provides chanel with accounts which we referred (has parent)
	GetChildren() chan *state.Account
	// GetTotalTokensAmount - returns total amount of each token
	GetTotalTokensAmount() map[state.AssetCode]int64
	// TokenBalances - returns chanel with tokens balances
	TokenBalances() chan *state.Balance
	// GetAssetByToken - returns asset by token. Panics if fails to find one
	GetAssetByToken(token state.AssetCode) state.AssetCode
	// GetOperationalAccount - returns Operational Account
	GetOperationalAccount() state.AccountID
	// GetAccount - returns account by address, panics if account does not exist
	GetAccount(address state.AccountID) *state.Account
	// GetMainBalanceForAsset - returns main balance for asset. panics if not found
	GetMainBalanceForAsset(accountID state.AccountID, asset state.AssetCode) *state.Balance
	// GetMasterAccount - returns ID of the master account
	GetMasterAccount() state.AccountID
	// GetCommissionAccount - returns Commission Account
	GetCommissionAccount() state.AccountID
	// GetSpecialAccounts - returns special accounts
	GetSpecialAccounts() *state.Accounts
	// PayoutCompleted - updates fees after payout for all balances
	PayoutCompleted()
}

// Builder - build snapshot based on state
type Builder struct {
	state                 statable
	tokenShareProvider    tokenShareProvider
	referralShareProvider referralShareProvider
	payoutBuilder         payoutBuilder
}

// NewBuilder - creates new instance of `Builder`
func NewBuilder(state statable) *Builder {
	return &Builder{
		state: state,
		tokenShareProvider: &tokenShareProviderImpl{
			state: state,
		},
		referralShareProvider: &referralShareProviderImpl{
			state: state,
		},
		payoutBuilder: &payoutBuilderImpl{
			state: state,
		},
	}
}

// Build - builds snapshot
func (b *Builder) Build() (*Snapshot, error) {
	snapshot := New()
	snapshot.Ledger = b.state.GetLedger()
	var err error
	snapshot.FeesToShareToTokenHolders, err = b.tokenShareProvider.GetFeesToShareToTokenHolders()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees to share to token holders")
	}

	snapshot.FeesToShareToParent, err = b.referralShareProvider.GetReferralSharePayout()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees to share with parents")
	}

	err = b.updatePayoutForOperationalAccount(snapshot.FeesToShareToTokenHolders[b.state.GetOperationalAccount()], snapshot.FeesToShareToParent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update payout for operational account")
	}

	err = b.payoutBuilder.BuildOperations(snapshot.SyncState, snapshot.FeesToShareToParent, payoutTypeReferral)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build operations for referral payout")
	}

	err = b.payoutBuilder.BuildOperations(snapshot.SyncState, snapshot.FeesToShareToTokenHolders, payoutTypeToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build operations for token payout")
	}

	b.state.PayoutCompleted()

	return snapshot, nil
}

// Referrers fees share is paid by operational account, so to simplify flow we'll subtract corresponding amount
// from operational account payout for holding tokens
func (b *Builder) updatePayoutForOperationalAccount(operationalAccountTokensShare map[state.AssetCode]int64, feesToShareToParent map[state.AccountID]map[state.AssetCode]int64) error {
	if operationalAccountTokensShare == nil {
		return nil
	}

	err := ensureAllNotNegative(operationalAccountTokensShare)
	if err != nil {
		return errors.Wrap(err, "invalid argument: expected operationalAccountTokensShare to have only non negative values")
	}

	for _, feeToSharePerAsset := range feesToShareToParent {
		err = ensureAllNotNegative(feeToSharePerAsset)
		if err != nil {
			return errors.Wrap(err, "invalid argument: expected feesToShareToParent to have only non negative values")
		}

		for asset, amountToShare := range feeToSharePerAsset {
			operationalAccountShareForAsset := operationalAccountTokensShare[asset]
			// overflow here is not possible - min value for operationalAccountShareForAsset is 0
			// amountToShare - is always positive
			updatedOperationAccountShareForAsset := operationalAccountShareForAsset - amountToShare

			// if this amount goes below 0, it's ok. We assume that it's responsibility of admin to monitor the state of
			// operational account, and if fees share for referrals greater than Operational account share in fee pool,
			// admin must deposit to operational account manually
			if updatedOperationAccountShareForAsset < 0 {
				updatedOperationAccountShareForAsset = 0
			}

			operationalAccountTokensShare[asset] = updatedOperationAccountShareForAsset
		}
	}

	return nil
}
