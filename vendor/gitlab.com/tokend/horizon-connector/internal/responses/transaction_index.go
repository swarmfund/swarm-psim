package responses

import "gitlab.com/tokend/horizon-connector/internal/resources"

type TransactionIndex struct {
	Embedded struct {
		Meta    resources.PageMeta      `json:"meta"`
		Records []resources.Transaction `json:"records"`
	} `json:"_embedded"`
}
