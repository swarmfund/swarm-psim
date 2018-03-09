package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"strings"
)

const (
	pairLength = 64
)

type Pair struct {
	ID     int64  `db:"id"`
	Public string `db:"public"`
	Secret string `db:"secret"`
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func safeString(raw []byte) string {
	return strings.Replace(base32.HexEncoding.EncodeToString(raw), "=", "", -1)
}

func GeneratePair() (*Pair, error) {
	b, err := randomBytes(pairLength)
	if err != nil {
		return nil, err
	}
	return &Pair{
		Public: safeString(b[:pairLength/2]),
		Secret: safeString(b[pairLength/2:]),
	}, nil
}

//todo refactor
func Verify(pair *Pair, body []byte, signature string) bool {
	msg := append(body, []byte(pair.Secret)...)
	msgHash := sha256.Sum256(msg)
	return base64.StdEncoding.EncodeToString(msgHash[:]) == signature
}
