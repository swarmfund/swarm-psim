package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

var ErrorUnsupportedOpType = errors.New("unknown operation type")

type OperationI interface {
	Populate(*Base, []byte) error
	CraftLetters(project string) ([]emails.NoticeLetterI, error)
	ParticipantsRequest() *ParticipantsRequest
	UpdateParticipants([]ApiParticipant)
}

func ParseOperation(base *Base, rawOperation []byte) (OperationI, error) {
	var op OperationI
	switch xdr.OperationType(base.TypeI) {
	case xdr.OperationTypePayment:
		op = &Payment{}
	case xdr.OperationTypeManageInvoice:
		op = &ManageInvoice{}
	case xdr.OperationTypeManageOffer:
		op = &Offer{}
	default:
		return nil, errors.Wrap(ErrorUnsupportedOpType, base.Type)
	}

	err := op.Populate(base, rawOperation)
	return op, err
}
