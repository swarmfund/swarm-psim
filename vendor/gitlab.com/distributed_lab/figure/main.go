package figure

import (
	"reflect"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	keyTag   = "fig"
	required = "required"
	ignore   = "-"
)

var ErrRequiredValue = errors.New("You must set the value in field")

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
	vle := reflect.Indirect(reflect.ValueOf(f.target))
	tpe := vle.Type()
	for fi := 0; fi < tpe.NumField(); fi++ {
		fieldType := tpe.Field(fi)
		fieldValue := vle.Field(fi)

		_, err := f.SetField(fieldValue, fieldType, keyTag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Figurator) SetField(fieldValue reflect.Value, field reflect.StructField, keyTag string) (bool, error) {

	tag, err := parseFieldTag(field, keyTag)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse tag")
	}

	if tag == nil {
		return false, nil
	}

	raw, hasRaw := f.values[tag.Key]

	isSet := false
	if hook, ok := f.hooks[field.Type.String()]; ok && hasRaw {
		value, err := hook(raw)
		if err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", field))
		}
		fieldValue.Set(value)
		isSet = true
	}

	if !isSet {
		if tag.IsRequired {
			return false, errors.Wrap(ErrRequiredValue, tag.Key)
		}
		return false, nil
	}

	return true, nil
}
