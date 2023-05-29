# meh

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![Go](https://github.com/lefinal/meh/workflows/Go/badge.svg?branch=main)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lefinal/meh)
[![GoReportCard example](https://goreportcard.com/badge/github.com/lefinal/meh)](https://goreportcard.com/report/github.com/lefinal/meh)
[![codecov](https://codecov.io/gh/lefinal/meh/branch/main/graph/badge.svg?token=ema8Z2HEk5)](https://codecov.io/gh/lefinal/meh)
[![GitHub issues](https://img.shields.io/github/issues/lefinal/meh)](https://github.com/lefinal/meh/issues)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/lefinal/meh)

Mastery of Error Handling

Convenient error handling for Go using custom error containers and bubbling features.

This document provides an overview over provided features but is not complete.
Documentation available [here](https://pkg.go.dev/github.com/lefinal/meh).

# Installation

In order to use this package, run:

```shell
go get github.com/lefinal/meh
```

# What errors are made of

Errors consist of the following properties:

## Code

Error codes that also describe the severity.
When an error is created, the code is set.
The following default codes ship with `meh`:

- _internal_: Basic internal errors like a failed database query.
- _bad-input_: Bad user input/request.
- _not-found_: The requested resource could not be found.
- _unauthorized_: Authentication is required for accessing the resource or performing the action.
- _forbidden_: Invalid permissions for accessing the resource or performing the action.
- _neutral_: Used for wrapping errors without changing the code.
- _(unexpected)_: No code specified.

Codes can be used for error handling based on the occurred problem.
They may also be used for choosing HTTP status codes (see the `mehhttp`-package).

Of course, you can define custom codes.
However, you should use the `__`-prefix in order to avoid collisions with codes being added natively in the future.
Example: `__myapp_my_code`

## Wrapped error

Errors are meant to be wrapped when being returned to the caller.
This allows details to bubble up and be logged as well as the error message to be displayed in a stacktrace-like manner.
If the error is a root error (the original cause), this property will not be set.

## Message

The actual error message.
This is a string and should describe the action that was performed.

**Bad example**: error while loading file

**Good example**: load file

The error messages will be concatenated with colons.
Therefore, the final format with the original error wrapped two times looks like this:

_get users: load user file: file not found_

## Details

Details are one of the main reasons why `meh` was developed.
They are of type `map[string]interface{}` and allow passing arbitrary details as key-value pairs.
Most times you will want to pass function call arguments here in order to allow easier inspection of issues.
Because of wrapping, details are persisted and returned back all to the top caller which handles the error.
If no details are provided, this can be kept unset (`nil`).

# Creating errors

Errors can be created manually using the `Error`-struct:

```go
return &meh.Error{
	Code: meh.ErrInternal,
	WrappedErr: err,
	Message: "read file",
	Details: meh.Details{
		"file_name": "my-file.txt"
    }
}
```

However, you want to use generators most times because of the syntactic sugar they provide.

General ones:
```go
func NewErr(code Code, message string, details Details) error
func NewErrFromErr(err error, code Code, message string, details Details) error
```
These allow creating a new error with the given code, message and details.
The ones with `FromErr`-suffix create a new error with the given one used as wrapped error.

Often, you use the native error codes.
That's why there are generators, including codes:

```go
func NewInternalErr(message string, details Details) error
func NewInternalErrFromErr(err error, message string, details Details) error
func NewBadInputErr(message string, details Details) error
func NewBadInputErrFromErr(err error, message string, details Details) error
func NewNotFoundErr(message string, details Details) error
func NewNotFoundErrFromErr(err error, message string, details Details) error
func NewUnauthorizedErr(message string, details Details) error
func NewUnauthorizedErrFromErr(err error, message string, details Details) error
func NewForbiddenErr(message string, details Details) error
func NewForbiddenErrFromErr(err error, message string, details Details) error
```

# Wrapping errors

Most of the time, you do not want to create new error but wrap it for passing it over to the caller.
This preserves the error code from the underlying error.

Let's have a look at an example which describes a situation where errors are wrapped:

```go
struct Fruit {
	// ...
}

func IncrementApplesForUser(userID uuid.UUID, includePineapples bool) error {
	apples, err := applesByUser(userID, includePineapples)
	if err!= nil {
		return meh.Wrap(err, "apples by user", meh.Details{
			"user_id": userID,
			"incude_pineapples": includePineapples,
		})
	}
	// ...
}

func applesByUser(userID uuid.UUID, includePineapples bool) (int, error) {
	fruits, err := fruitsByUser(userID)
	if err != nil {
		return 0, meh.Wrap(err, "fruits by user", meh.Details{"user_id": userID})
	}
	// ...
}

func fruitsByUser(userID uuid.UUID) ([]Fruit, error) {
	fruitsRaw, err := readFruitsFile()
	if err != nil {
		return nil, meh.Wrap(err, "read fruits file", nil)
	}
	// ...
}

func readFruitsFile() ([]byte, error) {
	const fruitsFilename = "fruits.txt"
	b, err := os.ReadFile(fruitsFilename)
	if err != nil {
		return nil, meh.NewInternalErrFromErr(err, "read fruits file", meh.Details{
			"fruits_filename": fruitsFilename,
		})
	}
	return b, nil
}
```

By wrapping the error, the `meh.ErrInternal`-code from `readFruitsFile` is preserved.
Details are passed as well and the final error message would look like this:

_apples by user: fruits by user: read fruits file: file not found_

Of course, the error code could also be changed.
For example, if the returned error from `os.ReadFile` is checked to be a `os.ErrNotExist` and `meh.ErrNotFound` is returned, `fruitsByUser` could then return `meh.NewInternalErrFromErr(err, ...)` in order to change the code to `ErrInternal` as `readFruitsFile` is expected to not fail.

# Checking the error code

As already mentioned, each layer of "wrapping" is represented another `meh.Error` with its own code.
If you want to check the actual error code, use `meh.ErrorCode(err error)`.
This will return the error code of the first error without `meh.ErrNeutral`-code, which is set when wrapping errors.

# Logging

[Documentation](https://pkg.go.dev/github.com/lefinal/meh/mehlog)

Support for logging with [zap](https://github.com/uber-go/zap) comes out of the box and is provided with the package `mehlog`.

Set the log-level translation with `mehlog.SetDefaultLevelTranslator` and log with `mehlog.Log`.
This logs the error to the level which is determined by the error code (same as `meh.ErrorCode`).

# HTTP support

[Documentation](https://pkg.go.dev/github.com/lefinal/meh/mehhttp)

In combination with `mehlog`, support for responding with the correct status code and logging error details with request details is provided.
Currently, features are rather limited and serve more as an example.

Set the status code mapping with `mehhttp.SetHTTPStatusCodeMapping`.
You can then log and respond using `mehhttp.LogAndRespondError`.
This logs the error along with request details and responds with the determined HTTP status code and an empty message.

The following additional error codes are provided:

- _mehhttp-communication_: Used for all problems regarding client communication because communication is unstable by nature and not always an internal error.
- _mehhttp-service-not-reachable_: Used for problems with requesting third-party services.

# PostgreSQL support

[Documentation](https://pkg.go.dev/github.com/lefinal/meh/mehpg)

Support for errors of type `pgconn.PgError` is provided using these generators:

```go
func NewQueryDBErr(err error, message string, query string) error
func NewScanRowsErr(err error, message string, query string) error 
```

`NewQueryDBErr` returns a `meh.ErrBadInput`-error if the error code has prefix 22 (data exception) or 23 (integrity constraint violation).
