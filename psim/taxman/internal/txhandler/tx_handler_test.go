package txhandler

import (
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/txhandler/mocks"
)

func TestTxHandler(t *testing.T) {
	Convey("Given valid txHandler", t, func() {
		statable := &mocks.Statable{}
		defer statable.AssertExpectations(t)
		handler := newTxHandler(statable, logan.New())
		Convey("Failed to unmarshal xdr", func() {
			err := handler.Handle(horizon.Transaction{})
			So(err, ShouldNotBeNil)
		})
		Convey("Updated payout period", func() {
			expectedPayout1 := xdr.Int64(120)
			expectedPayout2 := xdr.Int64(60)
			expectedPayout3 := xdr.Int64(math.MaxInt64)
			txEnvelope := xdr.TransactionEnvelope{
				Tx: xdr.Transaction{
					Operations: []xdr.Operation{
						{
							Body: xdr.OperationBody{
								Type:         xdr.OperationTypeSetOptions,
								SetOptionsOp: &xdr.SetOptionsOp{},
							},
						},
						{
							Body: xdr.OperationBody{
								Type: xdr.OperationTypeSetFees,
								SetFeesOp: &xdr.SetFeesOp{
									PayoutsPeriod: &expectedPayout1,
								},
							},
						},
						{
							Body: xdr.OperationBody{
								Type:      xdr.OperationTypeSetFees,
								SetFeesOp: &xdr.SetFeesOp{},
							},
						},
						{
							Body: xdr.OperationBody{
								Type: xdr.OperationTypeSetFees,
								SetFeesOp: &xdr.SetFeesOp{
									PayoutsPeriod: &expectedPayout2,
								},
							},
						},
						{
							Body: xdr.OperationBody{
								Type: xdr.OperationTypeSetFees,
								SetFeesOp: &xdr.SetFeesOp{
									PayoutsPeriod: &expectedPayout3,
								},
							},
						},
					},
				},
			}

			So(txEnvelope.Tx.SourceAccount.SetAddress("GDEAXXME3D62T4T4A4PKO3PTKVKU36EBRPKJXWAAP4FS4QSFLB5GXHX3"), ShouldBeNil)

			base64TxEnvelope, err := xdr.MarshalBase64(txEnvelope)
			So(err, ShouldBeNil)

			var actualPayoutPeriod *time.Duration
			statable.On("SetPayoutPeriod", mock.Anything).Run(func(args mock.Arguments) {
				actualPayoutPeriod = args.Get(0).(*time.Duration)
				So(*actualPayoutPeriod, ShouldEqual, time.Duration(expectedPayout1)*time.Second)
			}).Once()
			statable.On("SetPayoutPeriod", mock.Anything).Run(func(args mock.Arguments) {
				actualPayoutPeriod = args.Get(0).(*time.Duration)
				So(*actualPayoutPeriod, ShouldEqual, time.Duration(expectedPayout2)*time.Second)
			}).Once()
			statable.On("SetPayoutPeriod", mock.Anything).Run(func(args mock.Arguments) {
				actualPayoutPeriod = args.Get(0).(*time.Duration)
				So(actualPayoutPeriod, ShouldBeNil)
			}).Once()

			err = handler.Handle(horizon.Transaction{
				EnvelopeXDR: string(base64TxEnvelope),
			})
			So(err, ShouldBeNil)
		})
	})
}
