package horizon

type CoinEmissionRequest struct {
	Amount       string `json:"amount"`
	Approved     *bool  `json:"approved"`
	Asset        string `json:"asset"`
	CreatedAt    string `json:"created_at"`
	ExchangeName string `json:"exchange_name"`
	Issuer       string `json:"issuer"`
	ID           string `json:"paging_token"`
	Receiver     string `json:"receiver"`
	Reference    string `json:"reference"`
}

type CoinEmissionRequestsParams struct {
	Reference string
	Exchange  string
}

type CoinEmissionRequestsResponse struct {
	Embedded struct {
		Records []CoinEmissionRequest `json:"records"`
	} `json:"_embedded"`
}
