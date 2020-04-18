package kverrors_test

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/OrenRosen/kverrors"
	"github.com/stretchr/testify/require"
)

func TestIs(t *testing.T) {
	origErr := fmt.Errorf("original error")
	wrapErr := kverrors.Wrap(origErr, "mm wrapping")
	fmtErr := fmt.Errorf("fmt wrapping: %w", wrapErr)

	require.True(t, goerrors.Is(fmtErr, origErr), "wrapping errors failed to be")
	require.Equal(t, "fmt wrapping: mm wrapping: original error", fmtErr.Error())
}

type myError string

func (e myError) Error() string { return string(e) }

func TestAs(t *testing.T) {
	myErr := myError("my error")
	wrapErr := kverrors.Wrap(myErr, "mm wrapping")
	fmtErr := fmt.Errorf("fmt wrapping: %w", wrapErr)

	var myErr2 myError
	require.True(t, goerrors.As(fmtErr, &myErr2))
	require.Equal(t, myErr2.Error(), "my error")
}

func TestUnwrapAllWithFmt(t *testing.T) {
	myErr := myError("my error")
	err := kverrors.Wrap(myErr, "wrapping")
	errW := fmt.Errorf("wrapping: %w", err)
	errW = kverrors.Wrap(err, "mmwrapping")

	origErr := kverrors.UnwrapAll(errW)
	c, ok := origErr.(myError)

	require.True(t, ok, "failed unwrapping to myError")
	require.Equal(t, c.Error(), "my error")
}
