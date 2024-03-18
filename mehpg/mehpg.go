// Package mehpg provides error functionality regarding postgres-errors.
package mehpg

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	pgconnv5 "github.com/jackc/pgx/v5/pgconn"
	"github.com/lefinal/meh"
	"strings"
)

// Prefixes for PostgreSQL error codes.
//
// See: https://www.postgresql.org/docs/13/errcodes-appendix.html.
const (
	ErrCodePrefixDataException                  = "22"
	ErrCodePrefixIntegrityConstraintViolation   = "23"
	ErrCodePrefixSyntaxErrOrAccessRuleViolation = "42"
)

// NewQueryDBErr creates a new meh.Error with the given error and message and
// sets a field in details to the provided query. If the error is related to
// constraint violation or data exceptions, a meh.ErrBadInput will be returned.
// Otherwise, meh.ErrInternal.
func NewQueryDBErr(err error, message string, query string, args ...any) error {
	var finalDetailedErr error
	details := make(meh.Details)
	details["query"] = query
	details["args"] = args
	// Check if postgres error.
	var pgErr *pgconn.PgError
	var pgErrV5 *pgconnv5.PgError
	if errors.As(err, &pgErr) {
		details["pg_err"] = *pgErr
		details["sqlstate"] = pgErr.Code
		// Check for certain prefixes.
		if strings.HasPrefix(pgErr.Code, ErrCodePrefixIntegrityConstraintViolation) {
			// Constraint violation.
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "constraint violation",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErr.Code, ErrCodePrefixDataException) {
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "data exception",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErr.Code, ErrCodePrefixSyntaxErrOrAccessRuleViolation) {
			// Syntax error
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				Message:    "syntax error",
				WrappedErr: err,
			}
		} else {
			// Otherwise, probably internal error.
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				WrappedErr: err,
			}
		}
	} else if errors.As(err, &pgErrV5) {
		details["pg_err"] = *pgErrV5
		details["sqlstate"] = pgErrV5.Code
		// Check for certain prefixes.
		if strings.HasPrefix(pgErrV5.Code, ErrCodePrefixIntegrityConstraintViolation) {
			// Constraint violation.
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "constraint violation",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErrV5.Code, ErrCodePrefixDataException) {
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "data exception",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErrV5.Code, ErrCodePrefixSyntaxErrOrAccessRuleViolation) {
			// Syntax error
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				Message:    "syntax error",
				WrappedErr: err,
			}
		} else {
			// Otherwise, probably internal error.
			finalDetailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				WrappedErr: err,
			}
		}
	} else if errors.Is(err, sql.ErrTxDone) {
		finalDetailedErr = &meh.Error{
			Code:       meh.ErrInternal,
			Message:    "tx done",
			WrappedErr: err,
		}
	} else if errors.Is(err, sql.ErrConnDone) {
		finalDetailedErr = &meh.Error{
			Code:       meh.ErrInternal,
			Message:    "connection done",
			WrappedErr: err,
		}
	}
	if finalDetailedErr != nil && meh.ErrorCode(finalDetailedErr) != meh.ErrNeutral {
		return meh.Wrap(finalDetailedErr, message, details)
	}
	// Any other internal error.
	return meh.NewInternalErrFromErr(err, message, details)
}

// NewScanRowsErr creates a new meh.ErrInternal with the given error and message
// and saves the provided query to details. The used logic is currently the same
// as NewQueryDBErr but placed in a separate method for future adjustments.
func NewScanRowsErr(err error, message string, query string, args ...any) error {
	return NewQueryDBErr(err, message, query, args...)
}

// NewQueryAndScanRowsErr is used for errors returned from QueryRow with Scan.
// Further logic might be added in the future.
func NewQueryAndScanRowsErr(err error, message string, query string) error {
	return NewQueryDBErr(err, message, query)
}
