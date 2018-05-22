package figure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func TestFigurator_Please(t *testing.T) {
	type Config struct {
		SomeField     int `fig:"foo"`
		AnotherField  int `fig:"-"`
		RequiredField int `fig:"bar,required"`
	}

	cases := []struct {
		name       string
		configData map[string]interface{}
		expected   Config
		err        error
	}{
		{
			name: "check field with tag ignore is not set",
			configData: map[string]interface{}{
				"foo":         123,
				"another_var": 321,
				"bar":         666,
			},
			expected: Config{SomeField: 123, AnotherField: 0, RequiredField: 666},
			err:      nil,
		},
		{
			name: "check err if not enough data for tag required",
			configData: map[string]interface{}{
				"foo":         123,
				"another_var": 321,
			},
			expected: Config{},
			err:      ErrRequiredValue,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := Config{}
			err := Out(&config).From(c.configData).Please()
			if c.err == nil {
				assert.EqualValues(t, c.expected, config)
			}

			assert.Equal(t, c.err, errors.Cause(err))
		})
	}
}

func TestNoHook(t *testing.T) {

	type customType uint32
	type TestStruct struct {
		SomeField customType `fig:"some_field"`
	}

	testData := TestStruct{
		SomeField: 123,
	}

	err := Out(&testData).From(map[string]interface{}{"some_field": 123}).Please()
	assert.Error(t, err)
	assert.Equal(t, ErrNoHook, errors.Cause(err))

}
