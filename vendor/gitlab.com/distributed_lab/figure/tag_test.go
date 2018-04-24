package figure

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseFieldTag(t *testing.T) {
	cases := []struct {
		name        string
		field       reflect.StructField
		expectedTag *Tag
	}{
		{
			name:        `field name set as tag key in snake case`,
			field:       reflect.StructField{Name: `FooBar`, Tag: ``},
			expectedTag: &Tag{Key: `foo_bar`, IsRequired: false},
		},
		{
			name:        `check value was recognized`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"foo"`},
			expectedTag: &Tag{Key: `foo`, IsRequired: false},
		},
		{
			name:        `check ignore tag`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"-"`},
			expectedTag: nil,
		},
		{
			name:        `recognition the tag and attribute`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"foo,required"`},
			expectedTag: &Tag{Key: `foo`, IsRequired: true},
		},
		{
			name:        `tag with multiple attributes`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"foo,required,trash"`},
			expectedTag: &Tag{Key: `foo`, IsRequired: true},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualTag, err := parseFieldTag(c.field, keyTag)

			assert.Equal(t, actualTag, c.expectedTag)
			assert.NoError(t, err)
		})
	}
}

func TestParseFieldTagErr(t *testing.T) {
	cases := []struct {
		name        string
		field       reflect.StructField
		expectedErr error
	}{
		{
			name:        `Conflicting attributes`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"-,required"`},
			expectedErr: ErrConflictingAttributes,
		},
		{
			name:        `Unknown attribute`,
			field:       reflect.StructField{Name: ``, Tag: `fig:"foo,yoba"`},
			expectedErr: ErrUnknownAttribute,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := parseFieldTag(c.field, keyTag)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}

func TestFigurator_Please(t *testing.T) {
	var config struct {
		SomeVar     int `fig:"foo"`
		AnotherVar  int `fig:"-"`
		RequiredVar int `fig:"bar,required"`
	}

	figurator := Out(&config).From(map[string]interface{}{
		"foo":         123,
		"another_var": 321,
		"bar":         666,
	})
	var err error
	assert.NotPanics(t, func() {
		err = figurator.Please()
	})

	assert.NoError(t, err)
	assert.Equal(t, 123, config.SomeVar)
	assert.Equal(t, 0, config.AnotherVar)
	assert.Equal(t, 666, config.RequiredVar)
}

func TestFailedFigurator_Please(t *testing.T) {
	var config struct {
		SomeVar     int `fig:"foo"`
		AnotherVar  int `fig:"-"`
		RequiredVar int `fig:"bar,required"`
	}

	err := Out(&config).From(map[string]interface{}{
		"foo":         123,
		"another_var": 321,
	}).Please()

	if assert.Error(t, err) {
		assert.Equal(t, ErrRequiredValue, errors.Cause(err))
	}
}
