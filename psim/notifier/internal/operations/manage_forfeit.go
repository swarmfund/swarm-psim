package operations

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

type ManageForfeitRequest struct {
	Base
	Action      int32  `json:"action"`
	RequestID   uint64 `json:"request_id"`
	Amount      string `json:"amount"`
	Asset       string `json:"asset"`
	UserDetails string `json:"user_details"`
}

type ReviewForfeitRequest struct {
	Base
	RequestID uint64 `json:"request_id"`
	Balance   string `json:"balance"`
	Accept    bool   `json:"accept"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}

func (p *ManageForfeitRequest) Populate(base *Base, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.Base = *base
	return nil
}

// CraftLetters returns an array of messages to notify the operation participants.
func (p *ManageForfeitRequest) CraftLetters(project string) ([]emails.NoticeLetterI, error) {
	return []emails.NoticeLetterI{
		p.craftLetter(project),
	}, nil
}

func (p *ManageForfeitRequest) craftLetter(project string) (letter *emails.ForfeitNoticeLetter) {
	var letterType string
	if p.Asset == "USD" {
		letterType = "Withdraw"
	} else {
		letterType = "Redemption"
	}

	letter = &emails.ForfeitNoticeLetter{
		NoticeLetter: emails.NoticeLetter{
			Header:   fmt.Sprintf("%s | New %s", project, letterType),
			Template: emails.NoticeTemplateForfeit,
		},
		TransferNotice: emails.TransferNotice{
			Amount: fmt.Sprintf("%s %s", p.Amount, p.Asset),
			Date:   p.LedgerCloseTime.Format(emails.TimeLayout),
			Type:   letterType,
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

	var message string
	switch p.State {
	case opBaseStatePending:
		message = "%s request of %s created."
	case opBaseStateSuccess:
		message = "%s request of %s successfully done."
	case opBaseStateFailed:
		message = "%s request of %s failed."
	case opBaseStateRejected:
		message = "%s request of %s was declined."
	default:
		message = "%s request of %s."
	}

	letter.Status = opBaseStatuses[p.State]
	letter.Message = fmt.Sprintf(message, letter.Type, letter.Amount)
	return letter

}
