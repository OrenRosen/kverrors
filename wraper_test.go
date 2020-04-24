package kverrors_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/OrenRosen/kverrors"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	examples := []struct {
		err           error
		secondMsg     string
		secondKV      []interface{}
		expectedError string
		expectedKV    map[string]interface{}
		desc          string
	}{
		{
			err:           kverrors.New("some msg", "firstKey", "first value"),
			secondMsg:     "other msg",
			secondKV:      []interface{}{"second key", "second value"},
			expectedError: "other msg: some msg",
			expectedKV:    map[string]interface{}{"firstKey": "first value", "second key": "second value"},
		},
		{
			err:           io.EOF,
			secondMsg:     "other msg",
			secondKV:      []interface{}{"second key", "second value"},
			expectedError: "other msg: EOF",
			expectedKV:    map[string]interface{}{"second key": "second value"},
		},
		{
			err:           kverrors.Wrap(io.EOF, "other msg", "other key", "other value"),
			secondMsg:     "another msg",
			secondKV:      []interface{}{"second key", "second value"},
			expectedError: "another msg: other msg: EOF",
			expectedKV:    map[string]interface{}{"other key": "other value", "second key": "second value"},
		},
		{
			desc:          "odd number of keyvals",
			err:           kverrors.Wrap(io.EOF, "other msg", "key1", "val1"),
			secondMsg:     "another msg",
			secondKV:      []interface{}{"value with no key", "key2", "val2"},
			expectedError: "another msg: other msg: EOF",
			expectedKV:    map[string]interface{}{"key1": "val1", "key2": "val2", "missingKey": "value with no key"},
		},
		{
			desc:          "nil keyvals",
			err:           kverrors.Wrap(io.EOF, "other msg", "key1", "val1", "key2", "val2"),
			secondMsg:     "another msg",
			secondKV:      nil,
			expectedError: "another msg: other msg: EOF",
			expectedKV:    map[string]interface{}{"key1": "val1", "key2": "val2"},
		},
	}

	for _, e := range examples {
		err := kverrors.Wrap(e.err, e.secondMsg, e.secondKV...)
		require.NotNil(t, err)
		require.Equal(t, e.expectedError, err.Error(), e.desc)
		require.Equal(t, e.expectedKV, kverrors.KeyValsMap(err), e.desc)
	}
}

func TestWrapFmt(t *testing.T) {
	err := kverrors.New("oops", "k1", "v1")
	err = fmt.Errorf("fmt1: %w", err)
	err = fmt.Errorf("fmt2: %w", err)
	require.Equal(t, "fmt2: fmt1: oops", err.Error())

	err = kverrors.Wrap(err, "oops2", "k2", "v2")
	require.NotNil(t, err)
	require.Equal(t, "oops2: fmt2: fmt1: oops", err.Error())

	err = fmt.Errorf("fmt3: %w", err)
	require.Equal(t, "fmt3: oops2: fmt2: fmt1: oops", err.Error())

	expectedKV := []interface{}{"k2", "v2", "k1", "v1"}
	expectedKVMap := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	require.Equal(t, expectedKVMap, kverrors.KeyValsMap(err))
	require.Equal(t, expectedKV, kverrors.KeyVals(err))
}

type MyErr string

func (e MyErr) Error() string { return string(e) }

func TestWrapCustomType(t *testing.T) {
	myErr := MyErr("something failed")
	msg := "wraper"
	kv := []interface{}{"firstKey", "first value"}
	err := kverrors.Wrap(myErr, msg, kv...)

	msg2 := "wraper2"
	kv2 := []interface{}{"second key", "second value"}
	err2 := kverrors.Wrap(err, msg2, kv2...)

	msg3 := "wraper3"
	kv3 := []interface{}{"third key", 343}
	err3 := kverrors.Wrap(err2, msg3, kv3...)

	c := kverrors.UnwrapAll(err)
	require.Equal(t, myErr, c)
	expectedKV := map[string]interface{}{
		"firstKey":   "first value",
		"second key": "second value",
		"third key":  343,
	}
	require.Equal(t, expectedKV, kverrors.KeyValsMap(err3))

	expectedMsg := fmt.Sprintf("%s: %s: %s: %s", msg3, msg2, msg, "something failed")
	require.Equal(t, expectedMsg, err3.Error())

	// test stack
	expectedStack := err.(stacker).StackTrace()
	require.Equal(t, expectedStack, err3.(stacker).StackTrace())
}

func TestUnwrapAll(t *testing.T) {
	myErr := MyErr("something failed")
	require.Equal(t, myErr, kverrors.UnwrapAll(myErr))

	err := kverrors.Wrap(myErr, "oops")
	require.Equal(t, myErr, kverrors.UnwrapAll(err))

	err = fmt.Errorf("wrap: %w", err)
	require.Equal(t, myErr, kverrors.UnwrapAll(err))

	err = pkgerrors.Wrap(err, "pkg oops")
	require.Equal(t, myErr, kverrors.UnwrapAll(err))
}

type stacker interface {
	StackTrace() pkgerrors.StackTrace
}
