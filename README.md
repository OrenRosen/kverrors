# kverrors [![Build Status](https://travis-ci.org/OrenRosen/kverrors.svg?branch=master)](https://travis-ci.org/OrenRosen/kverrors) [![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/OrenRosen/kverrors)


pkg `kverrors` is a small wrapper to [`pkg/errors`](https://github.com/pkg/errors) for supporting structured errors for logging.

---
## why structured
From [go-kit/log](https://github.com/go-kit/kit/tree/master/log)
> Structured logging is, basically, conceding to the reality that logs are data, and warrant some level of schematic rigor. Using a stricter, key/value-oriented message format for our logs, containing contextual and semantic information, makes it much easier to get insight into the operational activity of the systems we build  
In short, A structured log message is not only a simple plain text message for a human to read, but a collection of key-values attributes which can be processed and analysed.

Errors most common method is logging, and treating errors as structured as well may give us better insight about our errors.
This package is an experiment for using structured errors for logging.

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

The function `KeyVals` returns the aggregated keyvals across the error chain. 
the error chain considered to be stopped when the error doesn't unwraps to an inner error.

```go
KeyVals(err error) []interface{}
```

The function `KeyValsMap` returns the keyvals as map of key-value pairs across the error chain.
This function have some constraints:
- If the keyvals length across the error chain isn't even, it adds a `"missing"` key.
- If a key isn't a string, it ignores it and put a placeholder.

```go
KeyValsMap(err error) map[string]interface{}
````


## Stack Trace

`kverrors` depends on `pkg/errors` for generating the stack trace.
All errors returned with pkg `kverrors` implements the `stacker` interfface{}
```go
type stacker interface {
    StackTrace() []pkgerrors.StackTrace
}
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
