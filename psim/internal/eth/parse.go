package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Marshal(tx types.Transaction) (string, error) {
	var buf bytes.Buffer

	if err := tx.EncodeRLP(&buf); err != nil {
		return "", errors.Wrap(err, "failed to encode tx")
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func Unmarshal(txHex string) (*types.Transaction, error) {
	rlpBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode tx hex")
	}

	var tx types.Transaction
	err = tx.DecodeRLP(rlp.NewStream(bytes.NewReader(rlpBytes), 0))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode tx rlp")
	}

	return &tx, nil
}
