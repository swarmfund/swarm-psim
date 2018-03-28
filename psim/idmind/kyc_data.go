package idmind

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// KYCData describes the structure of KYC blob retrieved form Horizon.
type KYCData struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Address    KYCAddress `json:"address"`
	ETHAddress string     `json:"eth_address"`
	Documents  Documents  `json:"documents"`
	Sequence   string     `json:"sequence"`
}

func (d KYCData) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"first_name":  d.FirstName,
		"last_name":   d.LastName,
		"address":     d.Address,
		"eth_address": d.ETHAddress,
		"documents":   d.Documents,
		"sequence":    d.Sequence,
	}
}

// KYCAddress is only a nested structure in KYCData structure.
type KYCAddress struct {
	Line1      string `json:"line_1"`
	Line2      string `json:"line_2"`
	City       string `json:"city"` // Use Detroit on Sandbox to receive failed result
	Country    string `json:"country"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
}

func (a KYCAddress) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"line_1":      a.Line1,
		"line_2":      a.Line2,
		"city":        a.City,
		"country":     a.Country,
		"state":       a.State,
		"postal_code": a.PostalCode,
	}
}

type Documents struct {
	KYCIdDocument     string `json:"kyc_id_document"`
	KYCProofOfAddress string `json:"kyc_poa"`
}

func (d Documents) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"kyc_id":               d.KYCIdDocument,
		"kyc_proof_of_address": d.KYCProofOfAddress,
	}
}

func parseKYCData(data string) (*KYCData, error) {
	var kycData KYCData
	err := json.Unmarshal([]byte(data), &kycData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal data bytes into KYCData structure")
	}

	return &kycData, nil
}
