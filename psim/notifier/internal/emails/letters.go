package emails

const (
	NoticeTemplatePayment NoticeTemplateType = 1 + iota
	NoticeTemplateInvoice
	NoticeTemplateOffer
	NoticeTemplateForfeit
	NoticeTemplateDeposit
	NoticeTemplateLowIssuance
)

var NoticeTemplate = map[NoticeTemplateType]string{
	NoticeTemplatePayment:     "payment",
	NoticeTemplateInvoice:     "manage_offer",
	NoticeTemplateOffer:       "manage_offer",
	NoticeTemplateForfeit:     "manage_offer",
	NoticeTemplateDeposit:     "review_coins_emission",
	NoticeTemplateLowIssuance: "base",
}
