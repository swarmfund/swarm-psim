package logan

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3/fields"
)

var (
	ErrorKey = logrus.ErrorKey
	StackKey = "stack"
)

type Entry struct {
	*logrus.Entry
}

// WithRecover creates error from the `recoverData` if it isn't actually an error already
// and returns Entry with this error and its stack.
func (e *Entry) WithRecover(recoverData interface{}) *Entry {
	err := errors.FromPanic(recoverData)
	return e.WithStack(err).WithError(err)
}

func (e *Entry) WithError(err error) *Entry {
	errorFields := errors.GetFields(err)

	return &Entry{
		Entry: e.Entry.WithFields(logrus.Fields(errorFields)).WithError(err),
	}
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(fields.Obtain(key, value))
}

func (e *Entry) WithFields(fields F) *Entry {
	return &Entry{e.Entry.WithFields(logrus.Fields(fields))}
}

func (e *Entry) WithStack(err error) *Entry {
	return e.WithField(StackKey, errors.GetStack(err))
}

// Debugf logs a message at the debug severity.
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.Debugf(format, args...)
}

// Debug logs a message at the debug severity.
func (e *Entry) Debug(args ...interface{}) {
	e.Entry.Debug(args...)
}

// Infof logs a message at the Info severity.
func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.Infof(format, args...)
}

// Info logs a message at the Info severity.
func (e *Entry) Info(args ...interface{}) {
	e.Entry.Info(args...)
}

// Warnf logs a message at the Warn severity.
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.Warnf(format, args...)
}

// Warn logs a message at the Warn severity.
func (e *Entry) Warn(args ...interface{}) {
	e.Entry.Warn(args...)
}

// Errorf logs a message at the Error severity.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.Errorf(format, args...)
}

// Error logs a message at the Error severity.
func (e *Entry) Error(args ...interface{}) {
	e.Entry.Error(args...)
}

// Panicf logs a message at the Panic severity.
func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Entry.Panicf(format, args...)
}

// Panic logs a message at the Panic severity.
func (e *Entry) Panic(args ...interface{}) {
	e.Entry.Panic(args...)
}
