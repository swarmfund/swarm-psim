package kyc

type Documents struct {
	IDDocument  IDDocument     `json:"kyc_id_document"`
	ProofOfAddr ProofOfAddrDoc `json:"kyc_poa"`
}

func (d Documents) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id_doc":           d.IDDocument,
		"proof_of_address": d.ProofOfAddr,
	}
}

type DocFile struct {
	Name string `json:"name"`
	ID   string `json:"key"`
}

func (f DocFile) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id": f.ID,
	}
}
