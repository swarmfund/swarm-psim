package horizon

type Asset struct {
	Code          string `json:"code"`
	CurrentPrice  string `json:"current_price"`
	PhysicalPrice string `json:"physical_price"`
}
