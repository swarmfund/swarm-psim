package mtproto

import (
	"encoding/hex"
)

type StringSecretsStorage struct {
	data []byte
}

func NewStringSecretsStorage(s string) (*StringSecretsStorage, error) {
	bb, err := hex.DecodeString(s)
	if err != nil {
		return &StringSecretsStorage{}, err
	}

	return &StringSecretsStorage{
		data: bb,
	}, nil
}

func (s StringSecretsStorage) Read(dst []byte) (int, error) {
	copy(dst, s.data)

	if len(dst) < len(s.data) {
		return len(dst), nil
	} else {
		return len(s.data), nil
	}
}

func (s *StringSecretsStorage) Write(p []byte) (int, error) {
	s.data = make([]byte, len(p))
	copy(s.data, p)

	if len(p) < len(s.data) {
		return len(p), nil
	} else {
		return len(s.data), nil
	}
}
