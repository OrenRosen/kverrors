package kverrors_test

import (
	"fmt"

	"github.com/OrenRosen/kverrors"
)

func ExampleWrap() {
	err := fmt.Errorf("oops")
	err = kverrors.Wrap(err, "wrappy",
		"id", "some_id",
		"count", 2,
	)

	err = kverrors.Wrap(err, "wrappy2",
		"field2", "value2",
	)

	fmt.Println(kverrors.KeyVals(err))

	// Output: [field2 value2 id some_id count 2]
}

func ExampleWrap_2() {
	err := fmt.Errorf("oops")
	err = kverrors.Wrap(err, "wrappy",
		"key", "value",
	)

	err = fmt.Errorf("oops2: %w", err)
	err = kverrors.Wrap(err, "wrappy2",
		"key2", "value2",
	)

	err = kverrors.Wrapf(err, "oops3")

	fmt.Println(err.Error())
	fmt.Println(kverrors.KeyVals(err))

	// Output: oops3: wrappy2: oops2: wrappy: oops
	// [key2 value2 key value]
}

func ExampleWrapf() {
	err := fmt.Errorf("oops")
	err = kverrors.Wrapf(err, "%s not found", "something")

	fmt.Println(err.Error())

	// Output: something not found: oops
}
