package utils

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair"
)

func GenerateToken() (string, error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}
	return kp.Address(), nil
}
