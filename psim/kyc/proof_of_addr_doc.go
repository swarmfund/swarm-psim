package kyc

type ProofOfAddrDoc struct {
	FaceFile DocFile `json:"front"`
}

func (d ProofOfAddrDoc) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"face": d.FaceFile.ID,
	}
}
