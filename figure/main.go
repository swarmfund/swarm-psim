package figure

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

const (
	tag = "fig"
)

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
type Figurator struct {
	values map[string]interface{}
	hooks  Hooks
	target interface{}
}

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
func (f *Figurator) With(hooks ...Hooks) *Figurator {
	merged := Merge(hooks...)
	f.hooks = merged
	return f
}

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
func Out(target interface{}) *Figurator {
	return &Figurator{
		target: target,
	}
}

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
func (f *Figurator) From(values map[string]interface{}) *Figurator {
	f.values = values
	return f
}

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
func (f *Figurator) Please() error {
	tpe := reflect.Indirect(reflect.ValueOf(f.target)).Type()
	vle := reflect.Indirect(reflect.ValueOf(f.target))
	for fi := 0; fi < tpe.NumField(); fi++ {
		fieldType := tpe.Field(fi)
		fieldValue := vle.Field(fi)
		figTag := fieldType.Tag.Get(tag)
		if figTag == "" {
			figTag = toSnakeCase(fieldType.Name)
		}
		raw, hasRaw := f.values[figTag]
		if !hasRaw {
			continue
		}
		if hook, ok := f.hooks[fieldType.Type.String()]; ok {
			value, err := hook(raw)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to figure out %s", fieldType.Name))
			}
			fieldValue.Set(value)
		}
	}

	return nil
}
