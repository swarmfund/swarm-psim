package operations

import (
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/psim/psim/notifier/internal/emails"
)

type Offer struct {
	Base
}

func (of *Offer) Populate(base *Base, _ []byte) error {
	of.Base = *base
	return nil
}

// CraftLetters returns an array of messages to notify the operation participants.
func (of *Offer) CraftLetters(project string) ([]emails.NoticeLetterI, error) {
	letters := make([]emails.NoticeLetterI, 0)
	for _, p := range of.Participants {
		if p.Email == "" {
			continue
		}

		localLetters, err := p.craftOfferLetters(project, of.ID, of.LedgerCloseTime)
		if err != nil {
			return nil, errors.Wrap(err, "canton create offer letter", logan.F{
				"participant_email":   p.Email,
				"participant_balance": p.BalanceID,
			})
		}

		letters = append(letters, localLetters...)
	}
	return letters, nil
}

func (p *Participant) craftOfferLetters(project string, opId int64, date time.Time) (lts []emails.NoticeLetterI, err error) {
	if p.Effects.MatchEffects == nil {
		return nil, errors.New("match effects empty")
	}
	var letterBase = emails.NoticeLetter{
		Email:    p.Email,
		Header:   fmt.Sprintf("%s | New Match", project),
		Template: emails.NoticeTemplateOffer,
	}

	effects := p.Effects.MatchEffects
	letters := make([]emails.NoticeLetterI, 0, len(effects.Matches))

	for i, m := range effects.Matches {
		letter := &emails.OfferNoticeLetter{
			NoticeLetter: letterBase,
			TransferNotice: emails.TransferNotice{
				Addressee: p.Email,
				Date:      date.Format(emails.TimeLayout),
				Type:      "Trade",
			},
			Price:       fmt.Sprintf("%s %s", m.Price, effects.QuoteAsset),
			QuoteAmount: fmt.Sprintf("%s %s", m.QuoteAmount, effects.QuoteAsset),
			Fee:         fmt.Sprintf("%s %s", m.FeePaid, effects.QuoteAsset),
		}

		quoteAmount, err := amount.Parse(m.QuoteAmount)
		if err != nil {
			return lts, err
		}
		fee, err := amount.Parse(m.FeePaid)
		if err != nil {
			return lts, err
		}

		var action string
		var orderPrice int64
		if effects.IsBuy {
			action = "buy"
			orderPrice = quoteAmount + fee
		} else {
			action = "sell"
			orderPrice = quoteAmount - fee
		}

		letter.ID = fmt.Sprintf("%d;%s;%d", opId, p.BalanceID, i)
		letter.Amount = fmt.Sprintf("%s %s", m.BaseAmount, effects.BaseAsset)
		letter.OrderPrice = fmt.Sprintf("%s %s", amount.String(orderPrice), effects.QuoteAsset)
		letter.Message = fmt.Sprintf("Your order for %s %s at a price of %s is fulfilled.",
			action, letter.Amount, letter.Price)

		letters = append(letters, letter)
	}

	return letters, nil
}
