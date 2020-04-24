package kverrors_test

import (
	"fmt"
	"testing"

	"github.com/OrenRosen/kverrors"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	examples := []struct {
		msg         string
		kv          []interface{}
		expectedKVs []interface{}
		expectedMap map[string]interface{}
		desc        string
	}{
		{
			desc:        "one pair of keyvals",
			msg:         "exemaple 1",
			kv:          []interface{}{"firstKey", "first value"},
			expectedKVs: []interface{}{"firstKey", "first value"},
			expectedMap: map[string]interface{}{"firstKey": "first value"},
		},
		{
			desc:        "many keyvals",
			msg:         "exemaple 2",
			kv:          []interface{}{"firstKey", "first value", "secondKey", "second value"},
			expectedKVs: []interface{}{"firstKey", "first value", "secondKey", "second value"},
			expectedMap: map[string]interface{}{"firstKey": "first value", "secondKey": "second value"},
		},
		{
			desc:        "empty keyvals",
			msg:         "exemaple 3",
			kv:          []interface{}{},
			expectedKVs: []interface{}{},
			expectedMap: map[string]interface{}{},
		},
		{
			desc:        "key not a string",
			msg:         "exemaple 4",
			kv:          []interface{}{2, "val2"},
			expectedKVs: []interface{}{2, "val2"},
			expectedMap: map[string]interface{}{"keyIsNotAStringError": "val2"},
		},
		{
			desc:        "no key valuse",
			msg:         "exemaple 5",
			kv:          nil,
			expectedKVs: []interface{}{},
			expectedMap: map[string]interface{}{},
		},
		{
			desc:        "empty value",
			msg:         "exemaple 6",
			kv:          []interface{}{"key1", ""},
			expectedKVs: []interface{}{"key1", ""},
			expectedMap: map[string]interface{}{"key1": "empty"},
		},
	}

	for _, e := range examples {
		err := kverrors.New(e.msg, e.kv...)
		require.Equal(t, e.msg, err.Error())
		require.Equal(t, e.expectedKVs, kverrors.KeyVals(err), e.desc)
		require.Equal(t, e.expectedMap, kverrors.KeyValsMap(err), e.desc)
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
		{
			msg:      "some msg with error format: %w",
			format:   []interface{}{fmt.Errorf("base error")},
			expected: "some msg with error format: base error",
		},
	}

	for _, example := range examples {
		err := kverrors.Errorf(example.msg, example.format...)
		if err == nil {
			t.Fatalf("err is nil after ErrorF")
		}

		require.Equal(t, example.expected, err.Error())
	}
}

type keyvaluer interface {
	KeyVals() []interface{}
}
