package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources/operations"

type CheckSaleStateOperationsIndex struct {
	Embedded struct {
		Records []operations.CheckSaleState `json:"records"`
	} `json:"_embedded"`
}
