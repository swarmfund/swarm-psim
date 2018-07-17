package internal

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

var (
	ErrAssetNotFound             = errors.New("asset not found")
	ErrNoAssetExternalSystemType = errors.New("asset external type not set")
)

//go:generate mockery -case underscore -output ./internal/mocks -name AssetsQ
type AssetsQ interface {
	ByCode(string) (*horizon.Asset, error)
}

// GetExternalSystemType will try to get external system type from asset details
func GetExternalSystemType(q AssetsQ, code string) (int32, error) {
	// external system type is not set, let's check asset details for that
	asset, err := q.ByCode(code)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get asset details")
	}
	if asset == nil {
		return 0, errors.From(ErrAssetNotFound, logan.F{"deposit_asset": code})
	}
	if asset.Details.ExternalSystemType == 0 {
		return 0, errors.From(ErrNoAssetExternalSystemType, logan.F{"deposit_asset": code})
	}
	return asset.Details.ExternalSystemType, nil
}

// MustGetExternalSystemType will try to get external system type from asset details
// and panic if it's not sure about result
func MustGetExternalSystemType(q AssetsQ, code string) int32 {
	system, err := GetExternalSystemType(q, code)
	if err != nil {
		panic(err)
	}
	return system
}
