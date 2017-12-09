package operations

import (
	"encoding/json"
	"fmt"

	"strings"

	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

type ManageInvoice struct {
	Base
	Amount          string  `json:"amount"`
	Asset           string  `json:"asset"`
	ReceiverBalance string  `json:"receiver_balance"`
	Sender          string  `json:"sender"`
	InvoiceID       int64   `json:"invoice_id"`
	RejectReason    *string `json:"reject_reason"`
}

const (
	InvoiceStatePending   int32 = 1
	InvoiceStateSuccess         = 2
	InvoiceStateRejected        = 3
	InvoiceStateCancelled       = 4
	InvoiceStateFailed          = 5
)

var InvoiceStatuses = map[int32]string{
	InvoiceStatePending:   "Pending",
	InvoiceStateSuccess:   "Success",
	InvoiceStateRejected:  "Rejected",
	InvoiceStateCancelled: "Cancelled",
	InvoiceStateFailed:    "Failed",
}

func (p *ManageInvoice) Populate(base *Base, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.Base = *base
	return nil
}

// CraftLetters returns an array of messages to notify the operation participants.
func (p *ManageInvoice) CraftLetters(project string) ([]emails.NoticeLetterI, error) {
	return []emails.NoticeLetterI{
		p.craftLetter(project, p.SourceAccount, false),
		p.craftLetter(project, p.Sender, true),
	}, nil
}

func (p *ManageInvoice) craftLetter(project, addressee string, isSender bool) (letter *emails.InvoiceNoticeLetter) {
	letter = &emails.InvoiceNoticeLetter{
		NoticeLetter: emails.NoticeLetter{
			Header:   fmt.Sprintf("%s | New Request", project),
			Template: emails.NoticeTemplateInvoice,
			// For situation if len(p.Participants) == 0
			ID: fmt.Sprintf("%d;%s;%s", p.ID, addressee, letter.Amount),
		},
		TransferNotice: emails.TransferNotice{
			Type:   "Payment Request",
			Amount: fmt.Sprintf("%s %s", p.Amount, p.Asset),
			Date:   p.LedgerCloseTime.Format(emails.TimeLayout),
		},
		Status: InvoiceStatuses[p.State],
	}

	for i, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		if participant.AccountID == addressee {
			letter.ID = fmt.Sprintf("%d;%s;%d", p.ID, participant.BalanceID, i)
			letter.Email = participant.Email
			letter.Addressee = participant.Email
		} else {
			if participant.Email != "" {
				letter.Counterparty = participant.Email
			} else {
				letter.Counterparty = participant.AccountID
			}
		}
	}

	var m string
	if isSender {
		letter.CounterpartyType = "Applicant"
		switch p.State {
		case InvoiceStatePending:
			m = "You have received a request for payment from %s."
		case InvoiceStateSuccess:
			m = "A request for payment from %s fulfilled!"
		case InvoiceStateRejected:
			m = "A request for payment from %s decline with reason: " + letter.RejectReason
		case InvoiceStateCancelled:
			m = "%s cancel his request for payment."
		case InvoiceStateFailed:
			m = "A request for payment from %s failed."
		default:
			m = "Request for payment from %s"
		}
	} else {
		letter.CounterpartyType = "Sender"
		switch p.State {
		case InvoiceStatePending:
			m = "You sent a payment request to %s."
		case InvoiceStateSuccess:
			m = "Your payment request to %s fulfilled!"
		case InvoiceStateRejected:
			m = "%s decline your payment request with reason: " +
				strings.Replace(letter.RejectReason, "%", "%%", -1)
		case InvoiceStateCancelled:
			m = "%s canceled."
		case InvoiceStateFailed:
			m = "Your payment request to %s was failed"
		default:
			m = "Request for payment to %s"
		}
	}

	letter.Message = fmt.Sprintf(m, letter.Counterparty)
	if p.RejectReason != nil {
		letter.RejectReason = *p.RejectReason
	} else {
		letter.RejectReason = "Reject Reason is not specified"
	}
	return letter
}
