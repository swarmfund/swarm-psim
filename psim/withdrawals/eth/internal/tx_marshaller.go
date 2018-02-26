package internal

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

type TxMarshaller struct{}

func (h *TxMarshaller) Marshal(tx *types.Transaction) (string, error) {
	var buf bytes.Buffer
	if err := tx.EncodeRLP(&buf); err != nil {
		return "", errors.Wrap(err, "failed to encode tx")
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (h *TxMarshaller) Unmarshal(txhex string) (*types.Transaction, error) {
	rlpbytes, err := hex.DecodeString(txhex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode tx hex")
	}

	var tx types.Transaction
	err = tx.DecodeRLP(rlp.NewStream(bytes.NewReader(rlpbytes), 0))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode tx rlp")
	}

	return &tx, nil
}
