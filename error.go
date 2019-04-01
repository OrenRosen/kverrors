package errors

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

type inner struct {
	msg        string
	keyvals    []interface{}
	stackTrace pkgerrors.StackTrace
}

// inner implements error, stacker and keyvaluer
func (e inner) Error() string                    { return e.msg }
func (e inner) StackTrace() pkgerrors.StackTrace { return e.stackTrace }
func (e inner) KeyVals() []interface{}           { return e.keyvals }

func New(msg string, keyvals ...interface{}) error {
	return &inner{
		msg:        msg,
		keyvals:    keyvals,
		stackTrace: newStackTrace(1),
	}
}

func Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return New(msg)
}

func Cause(err error) error {
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			break
		}
		err = c.Cause()
	}
	return err
}

func KeyVals(err error, keyvals ...interface{}) []interface{} {
	for err != nil {
		if kver, ok := err.(keyvaluer); ok {
			keyvals = append(kver.KeyVals(), keyvals...)
		}

		if causer, ok := err.(causer); ok {
			err = causer.Cause()
			continue
		}

		err = nil
	}

	return keyvals
}

// private

type causer interface {
	Cause() error
}

type stacker interface {
	StackTrace() pkgerrors.StackTrace
}

type keyvaluer interface {
	KeyVals() []interface{}
}

func newStackTrace(skip int) pkgerrors.StackTrace {
	pkgErr := pkgerrors.New("")
	stack := stackTraceFrom(pkgErr)
	if len(stack) > skip {
		stack = stack[skip:]
	}

	return stack
}

func stackTraceFrom(err error) pkgerrors.StackTrace {
	tracer, ok := Cause(err).(interface {
		StackTrace() pkgerrors.StackTrace
	})

	if !ok {
		return pkgerrors.StackTrace{}
	}

	s := tracer.StackTrace()
	return s
}
