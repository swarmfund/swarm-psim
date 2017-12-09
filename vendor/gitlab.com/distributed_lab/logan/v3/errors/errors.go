package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3/fields"
)

// FromPanic extracts the err from the result of a recover() call.
func FromPanic(rec interface{}) error {
	err, ok := rec.(error)
	if !ok {
		err = fmt.Errorf("%s", rec)
	}

	return err
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(msg string) error {
	return errors.New(msg)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
//
// Fields can optionally be added. If provided multiple - fields will be merged.
//
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string, errorFields... map[string]interface{}) error {
	wrapped := errors.Wrap(err, msg)
	if wrapped == nil {
		return nil
	}

	var mergedFields map[string]interface{}
	for _, f := range errorFields {
		mergedFields = fields.Merge(mergedFields, f)
	}

	return &withFields{
		wrapped,
		mergedFields,
	}
}

// From returns an error annotating err with a stack trace
// at the point From is called, and the provided fields.
//
// If err is nil, From returns nil.
func From(err error, fields map[string]interface{}) error {
	withStack := errors.WithStack(err)

	if withStack == nil {
		return nil
	}

	return &withFields{
		withStack,
		fields,
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	return errors.Cause(err)
}

// GetFields returns the underlying fields of the error and its nested cause-errors, if possible.
// An error value has fields if it (or any of its nested cause) implements the following interface:
//
//     type fieldsProvider interface {
//            GetFields() F
//     }
//
// If the error and all of its nested causes do not implement GetFields, empty fields map will
// be returned.
func GetFields(err error) map[string]interface{} {
	type fieldsProvider interface {
		GetFields() eFields
	}

	type causer interface {
		Cause() error
	}

	result := eFields{}
	for err != nil {
		fError, ok := err.(fieldsProvider)
		if ok {
			result = fields.Merge(result, fError.GetFields())
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return result
}

type withFields struct {
	error
	eFields
}

func (w *withFields) Error() string {
	return w.error.Error()
}

func (w *withFields) GetFields() eFields {
	return w.eFields
}

func (w *withFields) Cause() error {
	return w.error
}
