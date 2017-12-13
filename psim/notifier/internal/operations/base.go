package operations

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

const (
	opBaseStatePending  int32 = 1
	opBaseStateSuccess        = 2
	opBaseStateRejected       = 3
	opBaseStateFailed         = 4
)

var opBaseStatuses = map[int32]string{
	opBaseStatePending:  "Pending",
	opBaseStateSuccess:  "Success",
	opBaseStateRejected: "Rejected",
	opBaseStateFailed:   "Failed",
}

// Base is a common structure for operations.
type Base struct {
	ID              int64         `json:"id,string"`
	PT              string        `json:"paging_token"`
	SourceAccount   string        `json:"source_account"`
	Type            string        `json:"type"`
	TypeI           int32         `json:"type_i"`
	State           int32         `json:"state"`
	Identifier      string        `json:"identifier"`
	LedgerCloseTime time.Time     `json:"ledger_close_time"`
	Participants    []Participant `json:"participants,omitempty"`
}

// ParticipantsRequest is the request structure for the participant data.
type ParticipantsRequest struct {
	ForAccount   string                     `json:"for_account"`
	Participants map[int64][]ApiParticipant `json:"participants"`
}

// UpdateParticipants is merge present participants with ApiParticipant.
func (b *Base) UpdateParticipants(pMap []ApiParticipant) {
	for i, participant := range pMap {
		b.Participants[i].fromApiParticipant(&participant)
	}
}

// ParticipantsRequest returns ParticipantsRequest.
func (b *Base) ParticipantsRequest() *ParticipantsRequest {
	var ap = make([]ApiParticipant, len(b.Participants))
	for i, val := range b.Participants {
		ap[i] = ApiParticipant{
			AccountID: val.AccountID,
			BalanceID: val.BalanceID,
		}
	}

	return &ParticipantsRequest{
		ForAccount:   b.SourceAccount,
		Participants: map[int64][]ApiParticipant{b.ID: ap},
	}
}

func (b *Base) LogFields() logan.F {
	return logan.F{
		"op_type":       b.Type,
		"op_id":         b.ID,
		"op_close_time": b.LedgerCloseTime.Format(emails.TimeLayout),
	}
}
