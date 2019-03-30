package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type wraper struct {
	cause      error
	msg        string
	keyvals    []interface{}
	stackTrace errors.StackTrace
}

// inner implements error, stacker, keyvaluer and causer
func (w *wraper) Error() string                 { return w.msg + ": " + w.cause.Error() }
func (w *wraper) Cause() error                  { return w.cause }
func (w *wraper) StackTrace() errors.StackTrace { return w.stackTrace }
func (w *wraper) KeyVals() []interface{}        { return w.keyvals }

// Wrap Creates a new error with added msg and keyvals
// It saves the passed error as the cause, so it could trace back to it when logging/reporting
func Wrap(err error, msg string, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wraper{
		cause:      err,
		msg:        msg,
		keyvals:    keyvals,
		stackTrace: getOrNewStackTrace(err),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	return Wrap(err, msg)
}

// WrapAndMerge can wrap your custom error, and merge data from another error.
// Data means stack trace and keyvals and message
func WrapAndMerge(cause, err error, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wraper{
		cause:      cause,
		msg:        err.Error(),
		keyvals:    KeyVals(err, keyvals...),
		stackTrace: getOrNewStackTrace(err),
	}
}

// private

func getOrNewStackTrace(err error) errors.StackTrace {
	if tracer, ok := err.(stacker); ok {
		return tracer.StackTrace()
	}

	return newStackTrace(2)
}
