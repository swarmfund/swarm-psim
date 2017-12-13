package types

type Asset struct {
	Code                 string `json:"code"`
	Token                string `json:"token"`
	Owner                string `json:"owner"`
	AvailableForIssuance Amount `json:"available_for_issuance"`
}
