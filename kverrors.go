// Package kverrors is a small wrapper to https://github.com/pkg/errors for supporting key values. e.g structured errors.
//
// Adding keyvalues to errors is done by the method
//	Wrap(err error, msg string, keyvals ...interface{}) error
// Where keyvals are a key value pairs. The key must be string.
//	func DoSomethingWithUser(userID string) error {
//		user, err := userRepo.FindUser(123)
//		if err != nil {
//			return kverrors.Wrap(err, "DoSomethingWithUser",
//				"userID", 123,
//				"otherKey", "some value",
//			)
//		}
//
//		...
//	}
//
// Unwrapping
//
// Every error created with kverrors implements the unwrapper interface:
//
//	type unwrapper interface {
//		Unwrap() error
//	}
//
// For getting the original error, use the package function kverros.Unwrap
// which will recursively iterate through the error chain and will retrieve
// the original error, which is the first error which doesn't implement unwraper
// (or causer from pkg/erros)
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

// Errorf returns a new error with formatted message according to a format specifier.
// keyvals are structured key-value pairs, and usually used for infrastructure frameworks
// as logs and error reporting. (should not be analysed in the code, as the message isn't).
// the returned error holds a stacktrace created with pkg/errors
func Errorf(format string, keyvals ...interface{}) error {
	if strings.Contains(format, "%w") {
		return fmt.Errorf(format, keyvals...)
	}

	msg := fmt.Sprintf(format, keyvals...)
	return New(msg)
}

// KeyVals returns all the aggregated keyvals across the error chain
// the error chain considered to be stopped when the error doesn't
// unwraps to an inner error.
func KeyVals(err error) []interface{} {
	var keyvals []interface{}
	for err != nil {
		if kver, ok := err.(keyvaluer); ok {
			keyvals = append(keyvals, kver.KeyVals()...)
		}

		if unw, ok := err.(unwrapper); ok {
			err = unw.Unwrap()
		} else {
			err = nil
		}
	}

	if keyvals == nil {
		keyvals = []interface{}{}
	}

	return keyvals
}

// KeyValsMap returns the key value pairs across the error chain as a map
// the error chain considered to be stopped when the error doesn't
// unwraps to an inner error.
func KeyValsMap(err error) map[string]interface{} {
	keyvals := KeyVals(err)
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
	if keyvals == nil {
		return map[string]interface{}{}
	}

	if len(keyvals)%2 != 0 {
		temp := make([]interface{}, len(keyvals)+1)
		temp[0] = "missingKey"
		for i, kv := range keyvals {
			temp[i+1] = kv
		}

		keyvals = temp
	}

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
		keyStr := keyToString(key)
		extra[keyStr] = value
	}

	return extra
}

func keyToString(key interface{}) string {
	keyStr, ok := key.(string)
	if !ok {
		keyStr = "keyIsNotAStringError"
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
