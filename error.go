package kverrors

import (
	"fmt"
	"strings"

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

// New returns a new error with the provided msg and keyvals.
// keyvals are structured key-value pairs, and usually used for infrastructure frameworks
// as logs and error reporting. (should not be analysed in the code, as the message isn't).
// the returned error holds a stacktrace created with pkg/errors
func New(msg string, keyvals ...interface{}) error {
	return &inner{
		msg:        msg,
		keyvals:    keyvals,
		stackTrace: newStackTrace(1),
	}
}

// New returns a new error with formatted msg.
// the returned error holds a stacktrace created with pkg/errors
func Errorf(format string, args ...interface{}) error {
	if strings.Contains(format, "%w") {
		err := fmt.Errorf(format, args...)
		return Wrap(err, "Errorf")
	}

	msg := fmt.Sprintf(format, args...)
	return New(msg)
}

// KeyVals returns the key value pairs across the error chain
// the error chain considered to be stopped when the error doesn't
// unwraps to an inner error.
func KeyVals(err error, keyvals ...interface{}) map[string]interface{} {
	for err != nil {
		if kver, ok := err.(keyvaluer); ok {
			keyvals = append(kver.KeyVals(), keyvals...)
		}

		if unw, ok := err.(unwrapper); ok {
			err = unw.Unwrap()
			continue
		}

		err = nil
	}

	return paramsFromKeyvals(keyvals)
}

// private

type unwrapper interface {
	Unwrap() error
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
	tracer, ok := UnwrapAll(err).(interface {
		StackTrace() pkgerrors.StackTrace
	})

	if !ok {
		return pkgerrors.StackTrace{}
	}

	s := tracer.StackTrace()
	return s
}

func paramsFromKeyvals(keyvals []interface{}) map[string]interface{} {
	extra := make(map[string]interface{}, len(keyvals)/2)

	for i, key := range keyvals {
		if i%2 != 0 {
			continue
		}

		valIndex := i + 1
		if valIndex >= len(keyvals) {
			break
		}

		value := val(keyvals[valIndex])
		keyStr := keyStr(key)
		extra[keyStr] = value
	}

	return extra
}

func keyStr(key interface{}) string {
	keyStr, ok := key.(string)
	if !ok {
		keyStr = fmt.Sprintf("%v", key)
	}

	return keyStr
}

type valueFunc func() interface{}

func val(value interface{}) interface{} {
	if v, ok := value.(valueFunc); ok {
		value = v()
	}

	if value == "" {
		value = "empty"
	}

	return value
}
