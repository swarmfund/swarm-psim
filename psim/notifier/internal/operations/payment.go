package operations

import (
	"encoding/json"
	"fmt"

	"gitlab.com/tokend/psim/psim/notifier/internal/emails"
	"gitlab.com/tokend/psim/psim/notifier/internal/types"
)

type BasePayment struct {
	From                  string       `json:"from"`
	To                    string       `json:"to"`
	FromBalance           string       `json:"from_balance"`
	ToBalance             string       `json:"to_balance"`
	UserDetails           string       `json:"user_details"`
	Asset                 string       `json:"asset"`
	Amount                types.Amount `json:"amount"`
	SourcePaymentFee      types.Amount `json:"source_payment_fee"`
	DestinationPaymentFee types.Amount `json:"destination_payment_fee"`
	SourceFixedFee        types.Amount `json:"source_fixed_fee"`
	DestinationFixedFee   types.Amount `json:"destination_fixed_fee"`
	SourcePaysForDest     bool         `json:"source_pays_for_dest"`
}

type Payment struct {
	Base
	BasePayment
	Subject   string `json:"subject"`
	Reference string `json:"reference"`
	Asset     string `json:"asset"`
}

// Populate unmarshal Payment from raw JSON, merge it with Base.
func (p *Payment) Populate(base *Base, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.Base = *base
	return nil
}

// CraftLetters returns an array of messages to notify the payment participants.
func (p *Payment) CraftLetters(project string) ([]emails.NoticeLetterI, error) {
	return []emails.NoticeLetterI{
		p.craftLetter(project, p.From, true),
		p.craftLetter(project, p.To, false),
	}, nil
}

func (p *Payment) craftLetter(project, addressee string, isSender bool) (letter *emails.PaymentNoticeLetter) {
	letter = new(emails.PaymentNoticeLetter)
	letter.Template = emails.NoticeTemplatePayment
	letter.Header = fmt.Sprintf("%s | New Transaction", project)

	letter.Amount = fmt.Sprintf("%s %s", p.Amount, p.Asset)
	letter.Date = p.LedgerCloseTime.Format(emails.TimeLayout)

	// For situation if len(p.Participants) == 0
	letter.ID = fmt.Sprintf("%d;%s;%s", p.ID, addressee, letter.Amount)
	for _, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		if participant.AccountID == addressee {
			letter.ID = fmt.Sprintf("%d;%s", p.ID, participant.BalanceID)
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

	var fee types.Amount
	if isSender {
		if p.SourcePaysForDest {
			fee = p.SourceFixedFee + p.SourcePaymentFee + p.DestinationFixedFee + p.DestinationPaymentFee
		} else {
			fee = p.SourceFixedFee + p.SourcePaymentFee
		}
		letter.Fee = fmt.Sprintf("%s %s", fee, p.Asset)
		letter.FullAmount = fmt.Sprintf("%s %s", p.Amount+fee, p.Asset)
	} else {
		if p.SourcePaysForDest {
			letter.Fee = "Sender paid"
			letter.FullAmount = letter.Amount
		} else {
			fee = p.DestinationFixedFee + p.DestinationPaymentFee
			letter.Fee = fmt.Sprintf("%s %s", fee, p.Asset)
			letter.FullAmount = fmt.Sprintf("%s %s", p.Amount-fee, p.Asset)
		}
	}

	subjectPrefix := "No message"
	if len(p.Subject) > 4 {
		subjectPrefix = p.Subject[0:4]
		letter.Reference = p.Subject[4:]
	} else {
		letter.Reference = p.Subject
	}

	switch subjectPrefix {
	case "gf: ":
		letter.Type = "Gift"
	case "in: ":
		letter.Type = "Invoice"
	case "tf: ":
		letter.Type = "Transfer"
	default:
		letter.Type = "Transfer"
	}

	if isSender {
		letter.Action = "Paid"
		letter.CounterpartyType = "Receiver"
		letter.Message = fmt.Sprintf("You sent a %s to %s", letter.Amount, letter.Counterparty)
	} else {
		letter.Action = "Received"
		letter.CounterpartyType = "Sender"
		letter.Message = fmt.Sprintf("You received a %s from %s", letter.Amount, letter.Counterparty)
	}

	return letter
}
