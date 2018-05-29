package internal

import (
	"testing"
)

func TestTXMarshal(t *testing.T) {
	h := "f88908850ba43b740083061a809410460a6eff829c4e0c6f25cf4280fb5579b82f7080a42c48e7dbecf75b18e4a2182a214635d1051764125cf30e312e207e56380d55fdbc87d3031ca0ef3ca21876f5dc09ad0c4dad94b72f46dfb98d01df47333ad9481191ef8b1cdda07b451785f3057953287df7abca3079942bee196431615162d1fb70403e7ed670"
	m := TxMarshaller{}
	tx, err := m.Unmarshal(h)
	if err != nil {
		panic(err)
	}
	t.Log(tx.Nonce())
}
