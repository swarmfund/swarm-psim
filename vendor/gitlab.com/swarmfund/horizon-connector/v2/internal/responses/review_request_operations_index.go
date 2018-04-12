package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources/operations"

type ReviewRequestOperationIndex struct {
	Embedded struct {
		Records []operations.ReviewRequest `json:"records"`
	} `json:"_embedded"`
}
