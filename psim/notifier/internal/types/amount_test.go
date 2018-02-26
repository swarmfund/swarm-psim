package types

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestAmount_String(t *testing.T) {
	var (
		val  = Amount(421234)
		want = "0.421234"
	)
	got := fmt.Sprintf("%s", val)
	assert.Equal(t, got, want)
}

func TestAmount_UnmarshalJSON(t *testing.T) {
	var (
		valid = []byte(`{"v":"0.421234"}`)
		want  = Amount(421234)
	)
	dest := struct {
		V Amount `json:"v"`
	}{}
	err := json.Unmarshal(valid, &dest)
	assert.Equal(t, err, nil)
	assert.Equal(t, dest.V, want)
}

func TestAmount_MarshalJSON(t *testing.T) {
	var (
		valid = Amount(421234)
		want  = `{"v":"0.421234"}`
	)
	val := struct {
		V Amount `json:"v"`
	}{valid}
	result, err := json.Marshal(val)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(result), want)
}
