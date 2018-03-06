package internal

import (
	"testing"
)

func TestTXMarshal(t *testing.T) {
	h := "f8aa1085012a05f20083030d40949e88613418cf03dca54d6a2cf6ad934a78c7a17a80b844a9059cbb0000000000000000000000007189adeafc3ef75c7f35ad6b4e907969b7fc369d0000000000000000000000000000000000000000000000055a6e79ccd1d300001ca0e26d5fde788df5613ff0b7d2168c09fef133df7d42c007dca6d00030c2cf8aa7a04337fa6004a6b4394dc67f0f8d1c23bb2f599e17ac47a63d550eb86ddff421ad"
	m := TxMarshaller{}
	tx, err := m.Unmarshal(h)
	if err != nil {
		panic(err)
	}
	t.Log(tx.String())
}
