package figure

import (
	"reflect"

	"github.com/pkg/errors"
)

const (
	tag = "fig"
)

// Hook signature for custom hooks.
// Takes raw value expected to return target value
type Hook func(value interface{}) (reflect.Value, error)

// Hooks is mapping raw type -> `Hook` instance
type Hooks map[string]Hook

// With accepts hooks to be used for figuring out target from raw values.
// `BaseHooks` will be used implicitly if no hooks are provided
func (f *Figurator) With(hooks ...Hooks) *Figurator {
	merged := Hooks{}
	for _, partial := range hooks {
		for key, hook := range partial {
			merged[key] = hook
		}
	}
	f.hooks = merged
	return f
}

// Figurator holds state for chained call
type Figurator struct {
	values map[string]interface{}
	hooks  Hooks
	target interface{}
}

// Out is main entry point for package, used to start figure out chain
func Out(target interface{}) *Figurator {
	return &Figurator{
		target: target,
	}
}

// From takes raw config values to be used in figure out process
func (f *Figurator) From(values map[string]interface{}) *Figurator {
	f.values = values
	return f
}

// Please exit point for figure out chain.
// Will modify target partially in case of error
func (f *Figurator) Please() error {
	// if hooks were not explicitly set use default
	if len(f.hooks) == 0 {
		f.With(BaseHooks)
	}
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
				return errors.Wrapf(err, "failed to figure out %s", fieldType.Name)
			}
			fieldValue.Set(value)
		}
	}

	return nil
}
