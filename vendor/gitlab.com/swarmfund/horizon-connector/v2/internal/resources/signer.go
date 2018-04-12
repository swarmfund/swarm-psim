package resources

type Signer struct {
	PublicKey string `json:"public_key"`
	Weight    int32  `json:"weight"`
	Type      int32  `json:"signer_type_i"`
	Identity  int32  `json:"signer_identity"`
	Name      string `json:"signer_name"`
}
