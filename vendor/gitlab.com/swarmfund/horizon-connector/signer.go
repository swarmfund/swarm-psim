package horizon

type Signer struct {
	AccountID  string `json:"public_key"`
	Weight     int32  `json:"weight"`
	SignerType int32  `json:"signer_type_i"`
	Identity   int32  `json:"signer_identity"`
	Name       string `json:"signer_name"`
}
