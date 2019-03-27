package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	msg := "some msg"
	kv := []interface{}{"firstKey", "first value"}
	err := New(msg, kv...)
	require.Equal(t, msg, err.Error())
	kver, ok := err.(keyvaluer)
	require.True(t, ok)
	require.Equal(t, kv, kver.KeyVals())
	require.Equal(t, kv, KeyVals(err))
}

func TestErrorf(t *testing.T) {
	msg := "some msg %s"
	arg := "some arg"
	expected := "some msg some arg"
	err := Errorf(msg, arg)
	require.Equal(t, expected, err.Error())
}

func TestWrap(t *testing.T) {
	msg := "some msg"
	kv := []interface{}{"firstKey", "first value"}
	mmerr := New(msg, kv...)

	msg2 := "other msg"
	kv2 := []interface{}{"second key", "second value"}
	err2 := Wrap(mmerr, msg2, kv2...)

	expectedError := fmt.Sprintf("%s: %s", msg2, msg)
	require.Equal(t, expectedError, err2.Error())
	require.Equal(t, append(kv, kv2...), KeyVals(err2))

	// test stack
	tracer, ok := mmerr.(stacker)
	require.True(t, ok)
	expectedStack := tracer.StackTrace()
	require.Equal(t, expectedStack, err2.(stacker).StackTrace())
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
