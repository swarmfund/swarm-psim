package horizon

import (
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/network"
	"gitlab.com/tokend/go/xdr"
	"github.com/pkg/errors"
)

type TransactionBuilder struct {
	Envelope   string
	Source     keypair.KP
	Salt       uint64
	TimeBounds *xdr.TimeBounds
	Signatures []string
	Operations []xdr.Operation

	pendingSigns []keypair.KP

	ops       []Operation
	connector *Connector
	err       error
}

func (t *TransactionBuilder) GetSalt() xdr.Salt {
	// TODO revert it back when JS is fixed
	//if t.Salt == 0 {
	//	return xdr.Salt(rand.Uint32())
	//}
	return xdr.Salt(t.Salt)
}

func (t *TransactionBuilder) GetTimeBounds() xdr.TimeBounds {
	if t.TimeBounds == nil {
		return t.connector.TimeBounds()
	}
	return *t.TimeBounds
}

func (t *TransactionBuilder) OperationsCount() int {
	return len(t.ops) + len(t.Operations)
}

func (t *TransactionBuilder) Op(op Operation) *TransactionBuilder {
	t.ops = append(t.ops, op)
	return t
}

func (t *TransactionBuilder) Sign(kp keypair.KP) *TransactionBuilder {
	t.pendingSigns = append(t.pendingSigns, kp)
	return t
}

func (t *TransactionBuilder) XDR() (*xdr.Transaction, error) {
	transaction := xdr.Transaction{
		Salt:       t.GetSalt(),
		Operations: t.Operations,
		TimeBounds: t.GetTimeBounds(),
	}

	err := transaction.SourceAccount.SetAddress(t.Source.Address())
	return &transaction, err
}

func (t *TransactionBuilder) hash(tx *xdr.Transaction) (*Hash, error) {
	raw, err := network.HashTransaction(tx, t.connector.info.NetworkPassphrase)
	if err != nil {
		return nil, err
	}
	return &Hash{raw: raw}, nil
}

func (t *TransactionBuilder) Hash() (*Hash, error) {
	tx, err := t.XDR()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build xdr")
	}
	return t.hash(tx)
}

func (t *TransactionBuilder) Marshal64() (*string, error) {
	if t.err != nil {
		return nil, t.err
	}

	for _, op := range t.ops {
		xdrop, err := op.XDR()
		if err != nil {
			return nil, err
		}
		t.Operations = append(t.Operations, *xdrop)
	}

	transaction, err := t.XDR()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build xdr")
	}

	hash, err := t.hash(transaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash")
	}

	signatures := []xdr.DecoratedSignature{}

	for _, kp := range t.pendingSigns {
		signature, err := kp.SignDecorated(hash.Slice())
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, signature)
	}

	for _, signature := range t.Signatures {
		s := xdr.DecoratedSignature{}
		err := xdr.SafeUnmarshalBase64(signature, &s)
		if err != nil {
			return nil, nil
		}
		signatures = append(signatures, s)
	}

	envelope, err := xdr.MarshalBase64(&xdr.TransactionEnvelope{
		Tx:         *transaction,
		Signatures: signatures,
	})

	return &envelope, err
}

func (t *TransactionBuilder) Submit() error {

	env, err := t.Marshal64()
	if err != nil {
		return err
	}
	return t.connector.SubmitTX(*env)
}

type Operation interface {
	XDR() (*xdr.Operation, error)
}

type SetRateOp struct {
	BaseAsset     string
	QuoteAsset    string
	PhysicalPrice int64
}

func (op SetRateOp) XDR() (*xdr.Operation, error) {
	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageAssetPair,
			ManageAssetPairOp: &xdr.ManageAssetPairOp{
				Action:                  xdr.ManageAssetPairActionManageAssetPairUpdatePrice,
				Base:                    xdr.AssetCode(op.BaseAsset),
				Quote:                   xdr.AssetCode(op.QuoteAsset),
				PhysicalPrice:           xdr.Int64(op.PhysicalPrice),
				PhysicalPriceCorrection: xdr.Int64(0),
				MaxPriceStep:            xdr.Int64(0),
				Policies:                xdr.Int32(0),
			},
		},
	}, nil
}

type CoinsEmissionRequestOp struct {
	Reference string
	Receiver  string
	Asset     string
	Amount    int64
}

func (op CoinsEmissionRequestOp) XDR() (*xdr.Operation, error) {
	balanceID, err := ParseBalanceID(op.Receiver)
	if err != nil {
		return nil, err
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageCoinsEmissionRequest,
			ManageCoinsEmissionRequestOp: &xdr.ManageCoinsEmissionRequestOp{
				Action:    xdr.ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate,
				RequestId: 0,
				Receiver:  balanceID,
				Asset:     xdr.AssetCode(op.Asset),
				Amount:    xdr.Int64(op.Amount),
				Reference: xdr.String64(op.Reference),
			},
		},
	}, nil
}

type ReviewPaymentRequestOp struct {
	PaymentID uint64
	Accept    bool
}

func (op ReviewPaymentRequestOp) XDR() (*xdr.Operation, error) {
	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeReviewPaymentRequest,
			ReviewPaymentRequestOp: &xdr.ReviewPaymentRequestOp{
				PaymentId: xdr.Uint64(op.PaymentID),
				Accept:    op.Accept,
			},
		},
	}, nil
}

type PaymentOp struct {
	SourceBalanceID      string
	DestinationBalanceID string
	Amount               int64
	Reference            string
	Subject              string
}

func (op PaymentOp) XDR() (*xdr.Operation, error) {
	sourceBalance, err := ParseBalanceID(op.SourceBalanceID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse source balance")
	}

	destinationBalance, err := ParseBalanceID(op.DestinationBalanceID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse destination balance")
	}
	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypePayment,
			PaymentOp: &xdr.PaymentOp{
				SourceBalanceId:      sourceBalance,
				DestinationBalanceId: destinationBalance,
				Amount:               xdr.Int64(op.Amount),
				Subject:              xdr.String256(op.Subject),
				Reference:            xdr.String64(op.Reference),
				FeeData: xdr.PaymentFeeData{
					SourceFee: xdr.FeeData{
						PaymentFee: 0,
						FixedFee:   0,
					},
					DestinationFee: xdr.FeeData{
						PaymentFee: 0,
						FixedFee:   0,
					},
				},
			},
		},
	}, nil
}
