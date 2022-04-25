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
func NewBadInputErr(publicMessage string, details Details) error {
	return NewBadInputErrFromErr(nil, publicMessage, details)
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
func NewNotFoundErrFromErr(message string, details Details) error {
	return &Error{
		Code:    ErrNotFound,
		Message: message,
		Details: details,
	}
}
