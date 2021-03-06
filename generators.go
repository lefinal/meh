package meh

// NewInternalErrFromErr creates a new ErrInternal with the given error to be
// wrapped, message and details.
func NewInternalErrFromErr(err error, message string, details Details) error {
	return &Error{
		Code:       ErrInternal,
		WrappedErr: err,
		Message:    message,
		Details:    details,
	}
}

// NewInternalErr creates a new ErrInternal with the given message and details.
func NewInternalErr(message string, details Details) error {
	return &Error{
		Code:    ErrInternal,
		Message: message,
		Details: details,
	}
}

// NewBadInputErr creates a new ErrBadInput with the given message and details.
func NewBadInputErr(message string, details Details) error {
	return NewBadInputErrFromErr(nil, message, details)
}

// NewBadInputErrFromErr creates a new ErrBadInput with the given error to be
// wrapped, message and details.
func NewBadInputErrFromErr(err error, message string, details Details) error {
	return &Error{
		Code:       ErrBadInput,
		WrappedErr: err,
		Message:    message,
		Details:    details,
	}
}

// NewNotFoundErr creates a new ErrNotFound with the given message and
// details.
func NewNotFoundErr(message string, details Details) error {
	return &Error{
		Code:    ErrNotFound,
		Message: message,
		Details: details,
	}
}

// NewNotFoundErrFromErr creates a new ErrNotFound with the given error to be
// wrapped, message and details.
func NewNotFoundErrFromErr(err error, message string, details Details) error {
	return &Error{
		Code:       ErrNotFound,
		WrappedErr: err,
		Message:    message,
		Details:    details,
	}
}

// NewUnauthorizedErr creates a new ErrUnauthorized with the given message and
// details.
func NewUnauthorizedErr(message string, details Details) error {
	return &Error{
		Code:    ErrUnauthorized,
		Message: message,
		Details: details,
	}
}

// NewUnauthorizedErrFromErr creates a new ErrUnauthorized with the given error
// to be wrapped, message and details.
func NewUnauthorizedErrFromErr(err error, message string, details Details) error {
	return &Error{
		Code:       ErrUnauthorized,
		WrappedErr: err,
		Message:    message,
		Details:    details,
	}
}

// NewForbiddenErr creates a new ErrForbidden with the given message and
// details.
func NewForbiddenErr(message string, details Details) error {
	return &Error{
		Code:    ErrForbidden,
		Message: message,
		Details: details,
	}
}

// NewForbiddenErrFromErr creates a new ErrForbidden with the given error to be
// wrapped, message and details.
func NewForbiddenErrFromErr(err error, message string, details Details) error {
	return &Error{
		Code:       ErrForbidden,
		WrappedErr: err,
		Message:    message,
		Details:    details,
	}
}
