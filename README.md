# errors [![Build Status](https://travis-ci.org/OrenRosen/errors.svg?branch=master)](https://travis-ci.org/OrenRosen/errors)


pkg errors is a very minimal wrapper to [`pkg/errors`](https://github.com/pkg/errors) for supporting structured errors for logging
(`github.com/OrenRosen/errors` is temporary location)

---
## why structured
From [go-kit/log](https://github.com/go-kit/kit/tree/master/log)
> Structured logging is, basically, conceding to the reality that logs are data, and warrant some level of schematic rigor. Using a stricter, key/value-oriented message format for our logs, containing contextual and semantic information, makes it much easier to get insight into the operational activity of the systems we build  
In short, A structured log message is not only a simple plain text message for a human to read, but a collection of key-values attributes which can be processed and analysed.

Treating errors as structured as well may give us better insight about our errors.
This package is an experiment for using structured errors for logging as well. Although we may find this is not as useful as in logging general information.

## Wrapping an error:
Wrapping an error is done using `Wrap` function:
```golang
func Wrap(err error, msg string, keyvals ...interface{}) error
```
For example:
```golang
if err != nil {
    return errors.Wrap(err, "somePkg.someFunc",
        "someKey", "some value",
    )
}
```
## Unwrapping an error:
Like in `pkg/errors`, getting the wrapped error can be done by the `causer` interface, which each error created with this package implements
```golang
type causer interface {
        Cause() error
}
``` 
Meaning, calling `pkgerrors.Cause` is supported on errors of this package. For convinience this package implemts `Cause` as well
```golang
func Cause(err error) error
```

---
## Using custom Types:
Usually the error type which is reported (by 3rd party) is the type of the original error. This may be by supporting `pkg/errors` and retrieving the original error by the `Cause()` function   

Easiest way to do this is for example:

```golang
package example

type someError string
func (e someError) Error() string { return string(e) }

const veryBadError = someError("failed desperately")

// inside a function:
return errors.Wrap(veryBadError, "someFunc")
```

This way when reporting, you'll get the type of your error.

### Wrap and Merge error:
Another way to achieve this is by using `WrapAndMerge` function.\
```golang
func WrapAndMerge(cause, err error, keyvals ...interface{}) error
```
Let's say you want to report on your custom error, but also you don't want to lose the data for the other error, as it may support `pkg/errors` as well.\
Using `WrapAndMerge` you can set your custom error as the cause, and merge the data of the other error (message, structured data and stacktrace)


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
