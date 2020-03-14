package kverrors

import (
	"fmt"
	"io"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	examples := []struct {
		err           error
		msg           string
		kv            []interface{}
		expectedError string
		expectedKV    map[string]interface{}
	}{
		{
			err:           New("some msg", "firstKey", "first value"),
			msg:           "other msg",
			kv:            []interface{}{"second key", "second value"},
			expectedError: "other msg: some msg",
			expectedKV:    map[string]interface{}{"firstKey": "first value", "second key": "second value"},
		},
		{
			err:           io.EOF,
			msg:           "other msg",
			kv:            []interface{}{"second key", "second value"},
			expectedError: "other msg: EOF",
			expectedKV:    map[string]interface{}{"second key": "second value"},
		},
		{
			err:           Wrap(io.EOF, "other msg", "other key", "other value"),
			msg:           "another msg",
			kv:            []interface{}{"second key", "second value"},
			expectedError: "another msg: other msg: EOF",
			expectedKV:    map[string]interface{}{"other key": "other value", "second key": "second value"},
		},
	}

	for _, example := range examples {
		err := Wrap(example.err, example.msg, example.kv...)
		require.NotNil(t, err)
		require.Equal(t, example.expectedError, err.Error())
		require.Equal(t, example.expectedKV, KeyVals(err))
	}
}

func TestWrapFmt(t *testing.T) {
	err := New("oops", "k1", "v1")
	err = fmt.Errorf("fmt1: %w", err)
	err = fmt.Errorf("fmt2: %w", err)
	require.Equal(t, "fmt2: fmt1: oops", err.Error())

	err = Wrap(err, "oops2", "k2", "v2")
	require.NotNil(t, err)
	require.Equal(t, "oops2: fmt2: fmt1: oops", err.Error())

	err = fmt.Errorf("fmt3: %w", err)
	require.Equal(t, "fmt3: oops2: fmt2: fmt1: oops", err.Error())

	expectedKV := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	require.Equal(t, expectedKV, KeyVals(err))
}

type MyErr string

func (e MyErr) Error() string { return string(e) }

func TestWrapCustomType(t *testing.T) {
	myErr := MyErr("something failed")
	msg := "wraper"
	kv := []interface{}{"firstKey", "first value"}
	err := Wrap(myErr, msg, kv...)

	msg2 := "wraper2"
	kv2 := []interface{}{"second key", "second value"}
	err2 := Wrap(err, msg2, kv2...)

	msg3 := "wraper3"
	kv3 := []interface{}{"third key", 343}
	err3 := Wrap(err2, msg3, kv3...)

	c := UnwrapAll(err)
	require.Equal(t, myErr, c)
	expectedKV := map[string]interface{}{
		"firstKey":   "first value",
		"second key": "second value",
		"third key":  343,
	}
	require.Equal(t, expectedKV, KeyVals(err3))

	expectedMsg := fmt.Sprintf("%s: %s: %s: %s", msg3, msg2, msg, "something failed")
	require.Equal(t, expectedMsg, err3.Error())

	// test stack
	expectedStack := err.(stacker).StackTrace()
	require.Equal(t, expectedStack, err3.(stacker).StackTrace())
}

func TestUnwrapAll(t *testing.T) {
	myErr := MyErr("something failed")
	require.Equal(t, myErr, UnwrapAll(myErr))

	err := Wrap(myErr, "oops")
	require.Equal(t, myErr, UnwrapAll(err))

	err = fmt.Errorf("wrap: %w", err)
	require.Equal(t, myErr, UnwrapAll(err))

	err = pkgerrors.Wrap(err, "pkg oops")
	require.Equal(t, myErr, UnwrapAll(err))
}
