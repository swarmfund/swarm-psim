package figure_test

import (
	"reflect"
	"testing"
	"time"

	"math/big"

	"gitlab.com/distributed_lab/figure"
)

func TestSimpleUsage(t *testing.T) {
	type Config struct {
		SomeInt                int            `fig:"some_int"`
		SomeString             string         `fig:"some_string"`
		Missing                int            `fig:"missing"`
		Default                int            `fig:"default"`
		DurationStr            time.Duration  `fig:"duration_str"`
		DurationInt            time.Duration  `fig:"duration_int"`
		DurationPointer        *time.Duration `fig:"duration_pointer"`
		DurationPointerMissing *time.Duration `fig:"duration_pointer_missing"`
		BigInt                 *big.Int       `fig:"big_int"`
		BigIntStr              *big.Int       `fig:"big_int_str"`
		Int64                  int64          `fig:"int_64"`
		Uint                   uint           `fig:"uint"`
		Uint32                 uint32         `fig:"uint_32"`
		Uint64                 uint64         `fig:"uint_64"`
		Float64                float64        `fig:"float_64"`
	}

	c := Config{
		Default: 42,
	}
	err := figure.Out(&c).From(map[string]interface{}{
		"some_int":         1,
		"some_string":      "satoshi",
		"duration_str":     "1s",
		"duration_int":     1,
		"duration_pointer": "1h",
		"big_int":          42,
		"big_int_str":      "42",
		"int_64":           17,
		"uint":             18,
		"uint_32":          19,
		"uint_64":          20,
		"float_64":         21.9,
	}).Please()
	if err != nil {
		t.Fatalf("expected nil error got %s", err)
	}

	duration := 1 * time.Hour
	expectedConfig := Config{
		SomeInt:         1,
		SomeString:      "satoshi",
		Default:         42,
		DurationStr:     1 * time.Second,
		DurationInt:     1,
		DurationPointer: &duration,
		BigInt:          big.NewInt(42),
		BigIntStr:       big.NewInt(42),
		Int64:           17,
		Uint:            18,
		Uint32:          19,
		Uint64:          20,
		Float64:         21.9,
	}
	if !reflect.DeepEqual(c, expectedConfig) {
		t.Errorf("expected %#v got %#v", expectedConfig, c)
	}

}
