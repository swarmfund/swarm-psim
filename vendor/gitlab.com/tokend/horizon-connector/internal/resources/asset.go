package resources

import "gitlab.com/tokend/regources"

type Asset struct {
	Code                 string                 `json:"code"`
	Owner                string                 `json:"owner"`
	AvailableForIssuance Amount                 `json:"available_for_issuance"`
	Details              regources.AssetDetails `json:"details"`
}
