package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources/operations"

type CreateKYCRequestOperationIndex struct {
	Embedded struct {
		Records []operations.CreateKYCRequest `json:"records"`
	} `json:"_embedded"`
}
