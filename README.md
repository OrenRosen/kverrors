# kverrors [![Build Status](https://travis-ci.org/OrenRosen/kverrors.svg?branch=master)](https://travis-ci.org/OrenRosen/kverrors) [![GoDoc](https://godoc.org/github.com/OrenRosen/kverrors?status.svg)](http://godoc.org/github.com/OrenRosen/kverrors)


pkg `kverrors` is a small wrapper to [`pkg/errors`](https://github.com/pkg/errors) for supporting structured errors for logging.

---
## why structured
From [go-kit/log](https://github.com/go-kit/kit/tree/master/log)
> Structured logging is, basically, conceding to the reality that logs are data, and warrant some level of schematic rigor. Using a stricter, key/value-oriented message format for our logs, containing contextual and semantic information, makes it much easier to get insight into the operational activity of the systems we build  
In short, A structured log message is not only a simple plain text message for a human to read, but a collection of key-values attributes which can be processed and analysed.

Treating errors as structured as well may give us better insight about our errors.
This package is an experiment for using structured errors for logging as well.

## Wrapping an error:
Wrapping an error is done using `Wrap` function:
```golang
func Wrap(err error, msg string, keyvals ...interface{}) error
```
For example:
```golang
if err != nil {
    return kverrors.Wrap(err, "somePkg.someFunc",
        "someKey", "some value",
    )
}
```

## KeyVals  

The funciton
```go
kverrors.KeyVals(err error) map[string]interface{}
```
returns the key value pairs across the error chain. 
the error chain considered to be stopped when the error doesn't unwraps to an inner error.

## Stack Trace

`kverrors` depends on `pkg/errors` for generating the stack trace.
For getting the stack trace, the error returned from `Wrap` implements
```go
StackTrace() []pkgerrors.StackTrace
``` 

## Unwrapping an error

Errors returned in `kverrors` can unwrap to the inner error by implementing
```go
Unwrap() error
```

### Unwrap to the original error

Package `kverrors` exports function
```go
UnwrapAll()
```
USing this function you can get the original error in the error chain. This function
iteratively goes over the error chain using the method `Unwrap`, and returns the first
error which doesn't implement this interface.


## Support of Go 1.13 and pkg/errors

Implementing `Unwrap` makes it possible to use Go's `errors` package `Is` and `As` method.

Also, the error chain may contain any error which implement `Unwrap`,
meaning you can use `kverrors.Wrap` together with `pkgerrors.Wrap` and `fmt.Errorf("%w", err)`
