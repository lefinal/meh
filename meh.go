package meh

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime/debug"
	"strings"
)

// Code is the type of error in Error.
type Code string

const (
	// ErrUnexpected is the default error that is used when no other Code is
	// specified.
	ErrUnexpected Code = ""
	// ErrInternal is used for basic internal errors.
	ErrInternal Code = "internal"
	// ErrBadInput is used when submitted data was invalid. This should be used,
	// when handling external input, e.g. client data or database requests failing
	// because of constraint violation, etc.
	ErrBadInput Code = "bad-input"
	// ErrNotFound is used when requested resources where not found.
	ErrNotFound Code = "not-found"
	// ErrNeutral is used mainly for wrapping in order to not change the Code.
	ErrNeutral Code = "neutral"
	// ErrUnauthorized is used for when the caller is not known but a resource
	// required authorized access.
	ErrUnauthorized Code = "unauthorized"
	// ErrForbidden is used for unauthorized access to resources.
	ErrForbidden Code = "forbidden"
)

// Details are optionally provided error details in Error.Details that are used
// for easier debugging and error locating.
type Details map[string]interface{}

// StackTrace holds an errors.StackTrace as well as a formatted stack trace for
// usage in logging. You can add one using ApplyStackTrace.
type StackTrace struct {
	// StackTrace is the actual errors.StackTrace.
	StackTrace errors.StackTrace
	// StackTraceStr is a formatted stack trace from debug.Stack.
	StackTraceStr string
}

// genStackTrace generates a StackTrace from current position.
func genStackTrace(err error) StackTrace {
	return StackTrace{
		StackTrace:    errors.WithStack(err).(stackTracer).StackTrace(),
		StackTraceStr: string(debug.Stack()),
	}
}

// Error is the container for any relevant error information that needs to be
// kept when bubbling. For wrapping errors use Wrap. You can create an Error
// manually or by using generators like NewInternalErrFromErr.
type Error struct {
	// Code is the type of Error.
	//
	// Warning: The Code only applies to the current level. For checking the actual
	// level (ignoring ErrNeutral like being added by ApplyCode or Wrap), you need to
	// use ErrorCode!
	Code Code
	// WrappedErr is an optionally wrapped error for example added when wrapping
	// low-level errors or using Wrap.
	WrappedErr error
	// WrappedErrPassThrough describes whether the WrappedErr should be returned on
	// Finalize. This is useful if the wrapped error is relevant for other libraries,
	// etc.
	WrappedErrPassThrough bool
	// Message is an internal error message that is used when generating the error
	// message if not an empty string.
	Message string
	// Details is any optionally added information.
	Details Details
	// Trace is the stack trace to use (set it via ApplyStackTrace).
	Trace StackTrace
}

type jsonError struct {
	Code                  Code            `json:"code"`
	WrappedErr            json.RawMessage `json:"wrappedErr"`
	WrappedErrPassThrough bool            `json:"wrappedErrPassThrough"`
	Message               string          `json:"message"`
	Details               Details         `json:"details"`
	Trace                 StackTrace      `json:"trace"`
}

// MarshalJSON marshals the Error into a JSON representation.
func (e *Error) MarshalJSON() ([]byte, error) {
	var err error
	// Marshal wrapped error.
	var wrappedErrJSON json.RawMessage
	if e.WrappedErr != nil {
		var wrappedErrToMarshal any = e.WrappedErr
		if _, ok := e.WrappedErr.(*Error); !ok {
			// No meh error. Marshal and parse again into details map.
			wrappedErrToMarshalInMehRepresentation := &Error{
				Message: e.WrappedErr.Error(),
			}
			if wrappedForeignErrJSON, err := json.Marshal(e.WrappedErr); err != nil {
				// We cannot marshal the error. Fallback to native representation.
				e.Details["_native"] = fmt.Sprintf("%+v", e.WrappedErr)
			} else {
				err = json.Unmarshal(wrappedForeignErrJSON, &wrappedErrToMarshalInMehRepresentation.Details)
				if err != nil {
					// Some weird representation. We to the same as above and set the native
					// representation.
					e.Details["_native"] = fmt.Sprintf("%+v", e.WrappedErr)
				}
			}
			// Final representation of our foreign meh error to marshal is ready. Marshal it.
			wrappedErrToMarshal = wrappedErrToMarshalInMehRepresentation
		}
		wrappedErrJSON, err = json.Marshal(wrappedErrToMarshal)
		if err != nil {
			return nil, NewInternalErrFromErr(err, "marshal wrapped error", nil)
		}
	}
	// Marshal final JSON representation.
	eJSON := jsonError{
		Code:                  e.Code,
		WrappedErr:            wrappedErrJSON,
		WrappedErrPassThrough: e.WrappedErrPassThrough,
		Message:               e.Message,
		Details:               e.Details,
		Trace:                 e.Trace,
	}
	return json.Marshal(eJSON)
}

// UnmarshalJSON unmarshals an Error from JSON representation.
func (e *Error) UnmarshalJSON(data []byte) error {
	var eJSON jsonError
	err := json.Unmarshal(data, &eJSON)
	if err != nil {
		return NewInternalErrFromErr(err, "unmarshal", nil)
	}
	e.Code = eJSON.Code
	e.WrappedErr = nil
	e.WrappedErrPassThrough = eJSON.WrappedErrPassThrough
	e.Message = eJSON.Message
	e.Details = eJSON.Details
	e.Trace = eJSON.Trace
	// Unmarshal wrapped error.
	if eJSON.WrappedErr != nil && len(eJSON.WrappedErr) > 0 {
		var wrappedMehErr Error
		err = json.Unmarshal(eJSON.WrappedErr, &wrappedMehErr)
		if err != nil {
			// No meh error. Extract into regular error.
			e.WrappedErr = plainError{string(eJSON.WrappedErr)}
		} else {
			e.WrappedErr = &wrappedMehErr
		}
	}
	return nil
}

type plainError struct {
	message string
}

func (e plainError) Error() string {
	return e.message
}

// Error is used for implementing the error interface and printing the error
// string by unwrapping errors. The error string will not contain the error code
// or any further details but only messages.
func (e *Error) Error() string {
	segments := make([]string, 0)
	// Add each message if not empty.
	for it := NewErrorUnwrapper(e); it.Next(); {
		var message string
		// We cannot use the normal cast here, because it adds an extra wrapper to
		// non-meh errors.
		if cast, ok := it.Current().(*Error); ok {
			message = cast.Message
		} else {
			message = it.Current().Error()
		}
		// Skip empty messages.
		if message == "" {
			continue
		}
		segments = append(segments, message)
	}
	// Concatenate to classic Go-style error messages using colons.
	return strings.Join(segments, ": ")
}

// ErrorCode returns the first non ErrNeutral Code for the given error.
func ErrorCode(err error) Code {
	for it := NewErrorUnwrapper(err); it.Next(); {
		if c := Cast(it.Current()).Code; c != ErrNeutral {
			return c
		}
	}
	return ErrNeutral
}

// ApplyStackTrace applies the current stack trace to the given error.
func ApplyStackTrace(err error) error {
	e := Cast(err)
	e.Trace = genStackTrace(err)
	return e
}

// stackTracer is the interface for providing an errors.StackTrace that is used
// in the errors package.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// StackTrace returns the lowest-level errors.StackTrace for the Error if one
// was set.
func (e *Error) StackTrace() errors.StackTrace {
	var trace errors.StackTrace
	for it := NewErrorUnwrapper(e); it.Next(); {
		if stackTrace := Cast(it.Current()).Trace.StackTrace; stackTrace != nil {
			trace = stackTrace
		}
	}
	return trace
}

// Wrap wraps the given error with an ErrNeutral, the message and details. If no
// details should be added, set them nil. However, the parameter is mandatory in
// order to enforce providing as much detail as possible.
//
// Warning: Wrap will NOT add a stack trace unlike errors.Wrap as this seems
// inconvenient when calling Wrap while bubbling!
func Wrap(toWrap error, message string, details Details) error {
	return &Error{
		Code:       ErrNeutral,
		WrappedErr: toWrap,
		Message:    message,
		Details:    details,
	}
}

// ApplyCode wraps the given error with the Code.
func ApplyCode(err error, code Code) error {
	return &Error{
		Code:       code,
		WrappedErr: err,
	}
}

// ApplyDetails wraps the given error with an ErrNeutral and the given Details.
func ApplyDetails(err error, details Details) error {
	return &Error{
		Code:       ErrNeutral,
		WrappedErr: err,
		Details:    details,
	}
}

// Cast tries to Cast the given error to *Error. In case of failure, a new
// ErrUnexpected is created, wrapping the original error.
func Cast(err error) *Error {
	eRef, ok := err.(*Error)
	if ok {
		return eRef
	}
	return fromErr(err)
}

// fromErr creates a new Error from the given one. This wraps the error with
// ErrUnexpected and if a nil error is provided adds an error message as this
// should not happen.
func fromErr(err error) *Error {
	var errMessage string
	if err == nil {
		errMessage = "fromErr with nil error"
	}
	return &Error{
		Code:       ErrUnexpected,
		WrappedErr: err,
		Message:    errMessage,
	}
}

// Field names for usage in ToMap.
const (
	MapFieldErrorCode    = "x_code"
	MapFieldErrorMessage = "x_err_message"
)

// ToMap returns the details of the given error as a key-value map with appended
// enhanced information regarding the error itself (Error.Code to
// MapFieldErrorCode and the Error.Error-message to MapFieldErrorMessage).
func ToMap(err error) map[string]interface{} {
	e := Cast(err)
	m := make(map[string]interface{})
	// First, we add all details from the highest level to the lowest one.
	for it := NewErrorUnwrapper(err); it.Next(); {
		for k, v := range Cast(it.Current()).Details {
			m[fmt.Sprintf("%d/%s", it.Level(), k)] = v
		}
	}
	// Then we add all metadata.
	m[MapFieldErrorCode] = ErrorCode(e)
	m[MapFieldErrorMessage] = e.Error()
	return m
}

// NilOrWrap returns nil if the given error is nil or calls meh.Wrap on it
// otherwise.
func NilOrWrap(err error, message string, details Details) error {
	if err != nil {
		return Wrap(err, message, details)
	}
	return nil
}

// Finalize alters the given error for handing it off to other libraries. If the
// given error is nil, nil is returned. If it contains a pass-through, indicated
// by Error.WrappedErrPassThrough being true, the wrapped error is returned from
// the first one, where pass-through is set. Otherwise, the error is returned as
// with Cast.
func Finalize(err error) error {
	if err == nil {
		return nil
	}
	// Check if pass-through found.
	e, ok := err.(*Error)
	for e != nil {
		if !ok || e.WrappedErr == nil {
			break
		}
		if e.WrappedErr != nil && e.WrappedErrPassThrough {
			return e.WrappedErr
		}
		e, ok = e.WrappedErr.(*Error)
	}
	// No pass-through found.
	return Cast(err)
}

// ClearPassThrough returns the given error with all pass-through-fields (from
// wrapped errors as well).
func ClearPassThrough(err error) error {
	if err == nil {
		return nil
	}
	var clearedRoot *Error
	var cleared *Error
	e, ok := err.(*Error)
	for {
		if !ok && cleared == nil {
			return err
		}
		if !ok && cleared != nil {
			return clearedRoot
		}
		wrappedErr := &Error{
			Code:                  e.Code,
			WrappedErr:            e.WrappedErr,
			WrappedErrPassThrough: false,
			Message:               e.Message,
			Details:               e.Details,
			Trace:                 e.Trace,
		}
		if cleared == nil {
			clearedRoot = wrappedErr
			cleared = clearedRoot
		} else {
			cleared.WrappedErr = wrappedErr
			cleared = wrappedErr
		}
		if e.WrappedErr == nil {
			return clearedRoot
		}
		e, ok = e.WrappedErr.(*Error)
	}
}
