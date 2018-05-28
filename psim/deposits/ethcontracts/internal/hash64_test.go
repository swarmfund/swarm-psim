package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash64(t *testing.T) {
	t.Run("is deterministic", func(t *testing.T) {
		msg := []byte{42, 37, 83}
		assert.Equal(t, Hash64(msg), Hash64(msg))
	})

	t.Run("is uniform", func(t *testing.T) {
		a := []byte{1}
		b := []byte{2}
		assert.NotEqual(t, Hash64(a), Hash64(b))
	})
}
