package snapshoter

import (
	"math/rand"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

func TestSnapshotBuilder(t *testing.T) {
	Convey("Given valid builder", t, func() {
		statable := &mockStatable{}
		defer statable.AssertExpectations(t)

		builder := NewBuilder(statable)

		tokenShareProvider := &mockTokenShareProvider{}
		defer tokenShareProvider.AssertExpectations(t)
		builder.tokenShareProvider = tokenShareProvider

		referralShareProvider := &mockReferralShareProvider{}
		defer referralShareProvider.AssertExpectations(t)
		builder.referralShareProvider = referralShareProvider

		payoutBuilder := &mockPayoutBuilder{}
		defer statable.AssertExpectations(t)
		builder.payoutBuilder = payoutBuilder

		ledger := rand.Int63()
		statable.On("GetLedger").Return(ledger).Once()
		Convey("Failed to create fees to share to token holders", func() {
			expectedError := errors.New("failed to create share to token holders")
			tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(nil, expectedError).Once()
			result, err := builder.Build()
			So(err, ShouldNotBeNil)
			So(errors.Cause(err), ShouldEqual, expectedError)
			So(result, ShouldBeNil)
		})
		Convey("Failed to create fees to share to parent", func() {
			tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(map[state.AccountID]map[state.AssetCode]int64{}, nil).Once()
			expectedError := errors.New("failed to create share to referrals")
			referralShareProvider.On("GetReferralSharePayout").Return(nil, expectedError).Once()
			result, err := builder.Build()
			So(err, ShouldNotBeNil)
			So(errors.Cause(err), ShouldEqual, expectedError)
			So(result, ShouldBeNil)
		})
		Convey("updatePayoutForOperationalAccount", func() {
			operationalAccountID := state.AccountID("operational_account_id")
			statable.On("GetOperationalAccount").Return(operationalAccountID).Once()
			Convey("Invalid operationalAccountTokensShare", func() {
				tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(map[state.AccountID]map[state.AssetCode]int64{
					operationalAccountID: {
						state.AssetCode(""):     0,
						state.AssetCode("XAAU"): -100,
						state.AssetCode("XAAG"): 0,
					},
				}, nil).Once()
				referralShareProvider.On("GetReferralSharePayout").Return(map[state.AccountID]map[state.AssetCode]int64{}, nil).Once()
				result, err := builder.Build()
				So(err, ShouldNotBeNil)
				So(result, ShouldBeNil)
			})
			Convey("Invalid feesToShareToParent", func() {
				tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(map[state.AccountID]map[state.AssetCode]int64{
					operationalAccountID: {
						state.AssetCode(""):     0,
						state.AssetCode("XAAU"): 100,
						state.AssetCode("XAAG"): 0,
					},
				}, nil).Once()
				referralShareProvider.On("GetReferralSharePayout").Return(map[state.AccountID]map[state.AssetCode]int64{
					operationalAccountID: {
						state.AssetCode(""):     0,
						state.AssetCode("XAAU"): -100,
						state.AssetCode("XAAG"): 0,
					},
				}, nil).Once()
				result, err := builder.Build()
				So(err, ShouldNotBeNil)
				So(result, ShouldBeNil)
			})
		})
		Convey("Given valid referral and token holders share", func() {
			operationalAccountID := state.AccountID("operational_account_id")
			statable.On("GetOperationalAccount").Return(operationalAccountID).Once()
			tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(map[state.AccountID]map[state.AssetCode]int64{}, nil).Once()
			referralShareProvider.On("GetReferralSharePayout").Return(map[state.AccountID]map[state.AssetCode]int64{}, nil).Once()
			Convey("BuildOperations for referral", func() {
				expectedError := errors.New("failed to create payout ops for referrals")
				payoutBuilder.On("BuildOperations", mock.Anything, mock.Anything, payoutTypeReferral).Return(expectedError).Once()
				result, err := builder.Build()
				So(errors.Cause(err), ShouldEqual, expectedError)
				So(result, ShouldBeNil)
			})
			Convey("BuildOperations for tokens holders", func() {
				payoutBuilder.On("BuildOperations", mock.Anything, mock.Anything, payoutTypeReferral).Return(nil).Once()
				expectedError := errors.New("failed to create payout ops for token holders")
				payoutBuilder.On("BuildOperations", mock.Anything, mock.Anything, payoutTypeToken).Return(expectedError).Once()
				result, err := builder.Build()
				So(errors.Cause(err), ShouldEqual, expectedError)
				So(result, ShouldBeNil)
			})
		})
		Convey("Happy path", func() {
			operationalAccountID := state.AccountID("operational_account_id")
			statable.On("GetOperationalAccount").Return(operationalAccountID).Once()
			tokenShareProvider.On("GetFeesToShareToTokenHolders").Return(map[state.AccountID]map[state.AssetCode]int64{
				operationalAccountID: {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 3000,
					state.AssetCode("XAAG"): 0,
				},
				state.AccountID("random_account_1"): {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 100,
					state.AssetCode("XAAG"): 200,
				},
			}, nil).Once()

			expectedFeesToShareToTokenHolders := map[state.AccountID]map[state.AssetCode]int64{
				operationalAccountID: {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 2850,
					state.AssetCode("XAAG"): 0,
				},
				state.AccountID("random_account_1"): {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 100,
					state.AssetCode("XAAG"): 200,
				},
			}

			referralSharePayout := map[state.AccountID]map[state.AssetCode]int64{
				state.AccountID("random_account_1"): {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 100,
					state.AssetCode("XAAG"): 200,
				},
				state.AccountID("random_account_2"): {
					state.AssetCode(""):     0,
					state.AssetCode("XAAU"): 50,
					state.AssetCode("XAAG"): 500,
				},
			}
			referralShareProvider.On("GetReferralSharePayout").Return(referralSharePayout, nil).Once()
			payoutBuilder.On("BuildOperations", mock.Anything, referralSharePayout, payoutTypeReferral).Return(nil).Once()
			payoutBuilder.On("BuildOperations", mock.Anything, expectedFeesToShareToTokenHolders, payoutTypeToken).Return(nil).Once()
			statable.On("PayoutCompleted").Once()
			result, err := builder.Build()
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
		})
	})
}
