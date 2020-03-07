package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type wraper struct {
	wrappedErr error
	msg        string
	keyvals    []interface{}
	stackTrace errors.StackTrace
}

// wraper implements error, stacker, keyvaluer and unwrapper
func (w *wraper) Error() string                 { return w.msg + ": " + w.wrappedErr.Error() }
func (w *wraper) Unwrap() error                 { return w.wrappedErr }
func (w *wraper) StackTrace() errors.StackTrace { return w.stackTrace }
func (w *wraper) KeyVals() []interface{}        { return w.keyvals }

// Wrap returns a new error with added msg and keyvals.
// the returned error wraps supplied error.
// the stacktrace will be either of the deepest error which implement tracer
// or a new one from this location (by pkg/errors).
func Wrap(err error, msg string, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wraper{
		wrappedErr: err,
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

// UnwrapAll returns the most inner error.
// in cases all errors can unwrap, the returned error is the original error.
func UnwrapAll(err error) error {
	for err != nil {
		unw, ok := err.(unwrapper)
		if !ok {
			break
		}
		err = unw.Unwrap()
	}
	return err
}

// private

func getOrNewStackTrace(err error) errors.StackTrace {
	if tracer, ok := err.(stacker); ok {
		return tracer.StackTrace()
	}

	return newStackTrace(2)
}
