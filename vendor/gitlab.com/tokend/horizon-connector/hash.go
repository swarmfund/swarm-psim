package horizon

import "encoding/hex"

type Hash struct {
	raw [32]byte
}

func (h *Hash) Slice() []byte {
	return h.raw[:]
}

func (h *Hash) Hex() string {
	return hex.EncodeToString(h.Slice())
}
