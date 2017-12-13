package operations

type BaseEffects struct {
	*MatchEffects
}

type MatchEffects struct {
	BaseAsset  string         `json:"base_asset"`
	QuoteAsset string         `json:"quote_asset"`
	IsBuy      bool           `json:"is_buy"`
	Matches    []MatchDetails `json:"matches"`
}

type MatchDetails struct {
	BaseAmount  string `json:"base_amount"`
	QuoteAmount string `json:"quote_amount"`
	FeePaid     string `json:"fee_paid"`
	Price       string `json:"price"`
}
