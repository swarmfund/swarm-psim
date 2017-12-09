package resources

type SubmitResponse struct {
	Type   string `json:"type"`
	Extras struct {
		EnvelopeXDR string `json:"envelope_xdr"`
		ResultCodes struct {
			Transaction string   `json:"transaction"`
			Operations  []string `json:"operations"`
		} `json:"result_codes"`
		ResultXDR string `json:"result_xdr"`
	} `json:"extras"`
}
