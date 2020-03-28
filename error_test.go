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
			msg:         "exemaple 1",
			kv:          []interface{}{"firstKey", "first value"},
			expectedMap: map[string]interface{}{"firstKey": "first value"},
		},
		{
			msg:         "exemaple 1",
			kv:          []interface{}{"firstKey", "first value", "secondKey", "second value"},
			expectedMap: map[string]interface{}{"firstKey": "first value", "secondKey": "second value"},
		},
		{
			msg:         "exemaple 3",
			kv:          []interface{}{""},
			expectedMap: map[string]interface{}{},
		},
	}

	for _, example := range examples {
		err := New(example.msg, example.kv...)
		require.Equal(t, example.msg, err.Error())
		kver, ok := err.(keyvaluer)
		require.True(t, ok)
		require.Equal(t, example.expectedMap, kver.KeyVals(), example.msg)
		require.Equal(t, example.expectedMap, KeyVals(err), example.msg)
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
