package figure

import (
	"reflect"

	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrUnknownAttribute      = errors.New("Unknown syntax of tag")
	ErrConflictingAttributes = errors.New("Conflict attributes")
)

type Tag struct {
	Key        string
	IsRequired bool
}

func parseFieldTag(field reflect.StructField, tagKey string) (*Tag, error) {
	tag := &Tag{}

	fieldTag := field.Tag.Get(tagKey)
	if fieldTag == "" {
		tag.Key = toSnakeCase(field.Name)
		return tag, nil
	}

	splitedTag := strings.Split(fieldTag, `,`)
	tag.Key = splitedTag[0]

	if len(splitedTag) == 1 {
		if tag.Key == ignore {
			return nil, nil
		}
	}

	if len(splitedTag) > 1 {
		if tag.Key == ignore {
			return nil, ErrConflictingAttributes
		}

		if splitedTag[1] == required {
			tag.IsRequired = true
			return tag, nil
		}

		return nil, ErrUnknownAttribute
	}

	return tag, nil
}
