package log

import (
	"github.com/Sirupsen/logrus"
)

var (
	DefaultLogger *logrus.Logger
	DefaultEntry  *logrus.Entry
)

const (
	PanicLevel = logrus.PanicLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)


func WithField(key string, value interface{}) *logrus.Entry {
	result := DefaultEntry.WithField(key, value)
	return result
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return DefaultEntry.WithFields(fields)
}

// ===== Delegations =====

// Debugf logs a message at the debug severity.
func Debugf(format string, args ...interface{}) {
	DefaultEntry.Debugf(format, args...)
}

// Debug logs a message at the debug severity.
func Debug(args ...interface{}) {
	DefaultEntry.Debug(args...)
}

// Infof logs a message at the Info severity.
func Infof(format string, args ...interface{}) {
	DefaultEntry.Infof(format, args...)
}

// Info logs a message at the Info severity.
func Info(args ...interface{}) {
	DefaultEntry.Info(args...)
}

// Warnf logs a message at the Warn severity.
func Warnf(format string, args ...interface{}) {
	DefaultEntry.Warnf(format, args...)
}

// Warn logs a message at the Warn severity.
func Warn(args ...interface{}) {
	DefaultEntry.Warn(args...)
}

// Errorf logs a message at the Error severity.
func Errorf(format string, args ...interface{}) {
	DefaultEntry.Errorf(format, args...)
}

// Error logs a message at the Error severity.
func Error(args ...interface{}) {
	DefaultEntry.Error(args...)
}

// Panicf logs a message at the Panic severity.
func Panicf(format string, args ...interface{}) {
	DefaultEntry.Panicf(format, args...)
}

// Panic logs a message at the Panic severity.
func Panic(args ...interface{}) {
	DefaultEntry.Panic(args...)
}
