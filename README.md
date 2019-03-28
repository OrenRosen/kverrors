# errors [![Build Status](https://travis-ci.org/OrenRosen/errors.svg?branch=master)](https://travis-ci.org/OrenRosen/errors)
pkg errors is a wrapping to pkg/errors which adds key value params.

## Wrapping an error:
Wrapping an error is done using `Wrap` function:
```golang
func Wrap(err error, msg string, keyvals ...interface{}) error
```
For example:
```golang
if err != nil {
    return errors.Wrap(err, "somePkg.someFunc",
        "userSlug", "someUser",
    )
}
```
## Using custom Types:
For example:

```golang
package example

type someError string
func (e someError) Error() string { return string(e) }

const veryBadError = someError("failed desperately")

// inside a function:
return errors.Wrap(veryBadError, "someFunc")
```

This way when reporting, you'll get the type of your error.
> For getting the type and the stack trace in sentry, you need to wrap all the errors in the error chain

## Wrap and Merge error:
Let's say you want to report on your custom error, but also not to lose the data for other error.
This can be made with `WrapAndMerge` (tentative name)

```golang
package example

type someError string
func (e someError) Error() string { return string(e) }

const veryBadError = someError("failed desperately")

// inside a function:
if err != nil {
    return errors.WrapAndMerge(veryBadError, err)
}
```

## Unwrapping an error:
Unwrapping an error to its original error can be done using `Cause` function:
```golang
func Cause(err error) error
```
Get ke## yvals:
```golang
func KeyVals(err error, keyvals ...interface{}) error
```
Using this you can get the `keyvals` for an errors chain plus extra keyvals

TODOs
* enable formatting the error
* add `Wrapf` and `Errorf` for creating with `fmt.Sprintf`



