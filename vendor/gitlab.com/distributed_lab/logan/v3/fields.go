package logan

import "gitlab.com/distributed_lab/logan/v3/fields"

// F type is for fields, connected to `withFields` error.
type F map[string]interface{}

// WithField creates new `F` fields map and add provided key-value pair into it
// using Add method.
func Field(key string, value interface{}) F {
	result := make(F)
	result.Add(key, value)
	return result
}

// Add tries to extract fields from `value`, if `value` implements fields.Provider interface:
//
//		type Provider interface {
//			GetLoganFields() map[string]interface{}
//		}
//
// And adds these fields using AddFields.
// If `value` does not implement Provider - a single key-value pair is added.
func (f F) Add(key string, value interface{}) F {
	return f.AddFields(fields.Obtain(key, value))
}

// AddFields returns `F` map, which contains key-values from both maps.
// If both maps has some key - the value from the `newF` will be used.
func (f F) AddFields(newF F) F {
	return F(fields.Merge(f, newF))
}
