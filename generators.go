// Package meh adds typed errors, codes and reusable details to Go's native
// errors.
package meh

// NewErr creates a new Error with the given Code, message and details.
func NewErr(code Code, message string, details Details) error {
	return NewErrFromErr(nil, code, message, details)
}

// NewErrFromErr creates a new Error with the given wrapped one, Code, message
// and details.
func NewErrFromErr(err error, code Code, message string, details Details) error {
	return &Error{
		Code:       code,
		WrappedErr: ClearPassThrough(err),
		Message:    message,
		Details:    details,
	}
}

// NewPassThroughErr is similar to NewErrFromErr with the only difference that
// Error.WrappedErrPassThrough is set to true. This makes Finalize return the
// given error, even if wrapped inside *Error.
func NewPassThroughErr(err error, code Code, message string, details Details) error {
	return &Error{
		Code:                  code,
		WrappedErr:            err,
		WrappedErrPassThrough: true,
		Message:               message,
		Details:               details,
	}
}

// NewInternalErrFromErr creates a new ErrInternal with the given error to be
// wrapped, message and details.
func NewInternalErrFromErr(err error, message string, details Details) error {
	return NewErrFromErr(err, ErrInternal, message, details)
}

// NewInternalErr creates a new ErrInternal with the given message and details.
func NewInternalErr(message string, details Details) error {
	return NewErr(ErrInternal, message, details)
}

// NewBadInputErr creates a new ErrBadInput with the given message and details.
func NewBadInputErr(message string, details Details) error {
	return NewBadInputErrFromErr(nil, message, details)
}

// NewBadInputErrFromErr creates a new ErrBadInput with the given error to be
// wrapped, message and details.
func NewBadInputErrFromErr(err error, message string, details Details) error {
	return NewErrFromErr(err, ErrBadInput, message, details)
}

// NewNotFoundErr creates a new ErrNotFound with the given message and
// details.
func NewNotFoundErr(message string, details Details) error {
	return NewErr(ErrNotFound, message, details)
}

// NewNotFoundErrFromErr creates a new ErrNotFound with the given error to be
// wrapped, message and details.
func NewNotFoundErrFromErr(err error, message string, details Details) error {
	return NewErrFromErr(err, ErrNotFound, message, details)
}

// NewUnauthorizedErr creates a new ErrUnauthorized with the given message and
// details.
func NewUnauthorizedErr(message string, details Details) error {
	return NewErr(ErrUnauthorized, message, details)
}

// NewUnauthorizedErrFromErr creates a new ErrUnauthorized with the given error
// to be wrapped, message and details.
func NewUnauthorizedErrFromErr(err error, message string, details Details) error {
	return NewErrFromErr(err, ErrUnauthorized, message, details)
}

// NewForbiddenErr creates a new ErrForbidden with the given message and
// details.
func NewForbiddenErr(message string, details Details) error {
	return NewErr(ErrForbidden, message, details)
}

// NewForbiddenErrFromErr creates a new ErrForbidden with the given error to be
// wrapped, message and details.
func NewForbiddenErrFromErr(err error, message string, details Details) error {
	return NewErrFromErr(err, ErrForbidden, message, details)
}
