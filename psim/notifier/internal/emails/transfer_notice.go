package emails

type TransferNotice struct {
	Addressee string
	Amount    string
	Asset     string
	Date      string
	Link      string
	Type      string
}

type CoinsEmissionNoticeLetter struct {
	NoticeLetter
	TransferNotice
	Status       string
	RejectReason string
}

type ForfeitNoticeLetter struct {
	NoticeLetter
	TransferNotice
	Status string
}

type InvoiceNoticeLetter struct {
	NoticeLetter
	TransferNotice
	Counterparty     string
	CounterpartyType string
	Status           string
	RejectReason     string
}

type OfferNoticeLetter struct {
	NoticeLetter
	TransferNotice
	Price       string
	OrderPrice  string
	QuoteAmount string
	Fee         string
	Direction   string
}

type PaymentNoticeLetter struct {
	NoticeLetter
	TransferNotice
	Action           string
	Counterparty     string
	CounterpartyType string
	Fee              string
	FullAmount       string
	Reference        string
}
