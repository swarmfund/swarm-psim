package operations

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

type ManageCoinsEmissionRequest struct {
	Base
	RequestID int64  `json:"request_id"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}

type ReviewCoinsEmissionRequest struct {
	Base
	RequestID  uint64  `json:"request_id"`
	Amount     string  `json:"amount"`
	Asset      string  `json:"asset"`
	IsApproved *bool   `json:"approved"`
	Issuer     string  `json:"issuer"`
	Reason     *string `json:"reason"`
}

// Populate unmarshal ReviewCoinsEmissionRequest from raw JSON and merge it with operations.Base.
func (p *ReviewCoinsEmissionRequest) Populate(base *Base, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.Base = *base
	return nil
}

// CraftLetters returns an array of messages to notify the operation participants.
func (p *ReviewCoinsEmissionRequest) CraftLetters(project string) ([]emails.NoticeLetterI, error) {
	isManage := xdr.OperationType(p.TypeI) == xdr.OperationTypeManageCoinsEmissionRequest

	return []emails.NoticeLetterI{
		p.craftLetter(project, isManage),
	}, nil
}

func (p *ReviewCoinsEmissionRequest) craftLetter(project string, isManage bool) (letter *emails.CoinsEmissionNoticeLetter) {
	letter = &emails.CoinsEmissionNoticeLetter{
		NoticeLetter: emails.NoticeLetter{
			Header:   fmt.Sprintf("%s | New Deposit", project),
			Template: emails.NoticeTemplateDeposit,
		},
		TransferNotice: emails.TransferNotice{
			Date:   p.LedgerCloseTime.Format(emails.TimeLayout),
			Amount: fmt.Sprintf("%s %s", p.Amount, p.Asset),
			Type:   "Deposit",
		},
	}

	for _, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		letter.ID = fmt.Sprintf("%d;%s", p.ID, participant.BalanceID)
		letter.Email = participant.Email
		letter.Addressee = participant.Email
		break
	}

	var m string
	if isManage {
		letter.Status = opBaseStatuses[p.State]
		var b bool
		switch p.State {
		case opBaseStatePending:
			p.IsApproved = nil
		case opBaseStateSuccess:
			b = true
			p.IsApproved = &b
		case opBaseStateRejected:
			b = false
			p.IsApproved = &b
		}
	}

	if p.IsApproved == nil {
		letter.Status = "Pending"
		m = "A new deposit of %s has been requested to your account."
	} else if *p.IsApproved {
		letter.Status = "Success"
		m = "A new deposit of %s has been credited to your account."
	} else {
		letter.Status = "Rejected"
		m = "A deposit of %s has been rejected with reason: " +
			strings.Replace(letter.RejectReason, "%", "%%", -1)
	}

	letter.Message = fmt.Sprintf(m, letter.Amount)
	if p.Reason != nil {
		letter.RejectReason = *p.Reason
	} else {
		letter.RejectReason = "Reason is not specified"
	}

	return letter

}
