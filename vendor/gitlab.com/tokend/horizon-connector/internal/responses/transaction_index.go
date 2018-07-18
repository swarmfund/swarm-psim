package responses

import (
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/regources"
)

type TransactionIndex struct {
	Embedded struct {
		Meta    resources.PageMeta      `json:"meta"`
		Records []resources.Transaction `json:"records"`
	} `json:"_embedded"`
}

type TransactionV2Index struct {
	Embedded struct {
		Meta    regources.PageMeta        `json:"meta"`
		Records []regources.TransactionV2 `json:"records"`
	} `json:"_embedded"`
}