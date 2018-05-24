package internal

import (
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

type BTCFamilyMaster struct {
	master *hdkeychain.ExtendedKey
}

func NewBTCFamilyMaster(network NetworkType) (*BTCFamilyMaster, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate seed")
	}
	master, err := hdkeychain.NewMaster(seed, NetworkParams(network))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate master")
	}
	return &BTCFamilyMaster{master}, nil
}

func (s *BTCFamilyMaster) ExtendedPrivate() (string, error) {
	return s.master.String(), nil
}

func (s *BTCFamilyMaster) ExtendedPublic() (string, error) {
	public, err := s.master.Neuter()
	if err != nil {
		return "", err
	}
	return public.String(), nil
}
