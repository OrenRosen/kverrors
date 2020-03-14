package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	examples := []struct {
		msg         string
		kv          []interface{}
		expectedMap map[string]interface{}
	}{
		{
			msg:         "some msg",
			kv:          []interface{}{"firstKey", "first value"},
			expectedMap: map[string]interface{}{"firstKey": "first value"},
		},
		{
			msg:         "some msg",
			kv:          []interface{}{"firstKey", "first value", "secondKey", "second value"},
			expectedMap: map[string]interface{}{"firstKey": "first value", "secondKey": "second value"},
		},
		{
			msg:         "some msg",
			kv:          []interface{}{""},
			expectedMap: map[string]interface{}{},
		},
	}

	for _, example := range examples {
		err := New(example.msg, example.kv...)
		require.Equal(t, example.msg, err.Error())
		kver, ok := err.(keyvaluer)
		require.True(t, ok)
		require.Equal(t, example.kv, kver.KeyVals())
		require.Equal(t, example.expectedMap, KeyVals(err))
	}
}

func TestErrorf(t *testing.T) {
	examples := []struct {
		msg      string
		format   []interface{}
		expected string
	}{
		{
			msg:      "some msg without format",
			format:   []interface{}{},
			expected: "some msg without format",
		},
		{
			msg:      "some msg with %d format",
			format:   []interface{}{1},
			expected: "some msg with 1 format",
		},
	}

	for _, example := range examples {
		err := Errorf(example.msg, example.format...)
		require.Equal(t, example.expected, err.Error())
	}
}
