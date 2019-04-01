package errors

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestWrap(t *testing.T) {
	examples := []struct {
		err           error
		msg           string
		kv            []interface{}
		expectedError string
		expectedKV    []interface{}
	}{
		{New("some msg", "firstKey", "first value"), "other msg", []interface{}{"second key", "second value"},
			"other msg: some msg", []interface{}{"firstKey", "first value", "second key", "second value"}},
		{io.EOF, "other msg", []interface{}{"second key", "second value"},
			"other msg: EOF", []interface{}{"second key", "second value"}},
		{Wrap(io.EOF, "other msg", "other key", "other value"), "another msg", []interface{}{"second key", "second value"},
			"another msg: other msg: EOF", []interface{}{"other key", "other value", "second key", "second value"}},
	}

	for _, example := range examples {
		err := Wrap(example.err, example.msg, example.kv...)
		require.Equal(t, example.expectedError, err.Error())
		require.Equal(t, example.expectedKV, KeyVals(err))
	}
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

	c := Cause(err)
	require.Equal(t, myErr, c)
	expectedKV := append(kv, append(kv2, kv3...)...)
	require.Equal(t, expectedKV, KeyVals(err3))

	expectedMsg := fmt.Sprintf("%s: %s: %s: %s", msg3, msg2, msg, "something failed")
	require.Equal(t, expectedMsg, err3.Error())

	// test stack
	expectedStack := err.(stacker).StackTrace()
	require.Equal(t, expectedStack, err3.(stacker).StackTrace())

}
