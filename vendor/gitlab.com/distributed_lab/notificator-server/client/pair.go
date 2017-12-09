package notificator

import (
	"crypto/sha256"
	"encoding/base64"
)

type Pair struct {
	Public string
	Secret string
}

func (p *Pair) Signature(msg []byte) string {
	msg = append(msg, []byte(p.Secret)...)
	hash := sha256.Sum256(msg)
	return base64.StdEncoding.EncodeToString(hash[:])
}
