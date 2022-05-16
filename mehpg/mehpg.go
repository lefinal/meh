package mehpg

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
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
func NewQueryDBErr(err error, message string, query string) error {
	detailedErr := err
	details := make(meh.Details)
	details["query"] = query
	// Check if postgres error.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		details["pg_err"] = *pgErr
		details["sqlstate"] = pgErr.Code
		// Check for certain prefixes.
		if strings.HasPrefix(pgErr.Code, ErrCodePrefixIntegrityConstraintViolation) {
			// Constraint violation.
			detailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "constraint violation",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErr.Code, ErrCodePrefixDataException) {
			detailedErr = &meh.Error{
				Code:       meh.ErrBadInput,
				Message:    "data exception",
				WrappedErr: err,
			}
		} else if strings.HasPrefix(pgErr.Code, ErrCodePrefixSyntaxErrOrAccessRuleViolation) {
			// Syntax error
			detailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				Message:    "syntax error",
				WrappedErr: err,
			}
		} else {
			// Otherwise, probably internal error.
			detailedErr = &meh.Error{
				Code:       meh.ErrInternal,
				WrappedErr: err,
			}
		}
	} else if errors.Is(err, sql.ErrTxDone) {
		detailedErr = &meh.Error{
			Code:       meh.ErrInternal,
			Message:    "tx done",
			WrappedErr: err,
		}
	} else if errors.Is(err, sql.ErrConnDone) {
		detailedErr = &meh.Error{
			Code:       meh.ErrInternal,
			Message:    "connection done",
			WrappedErr: err,
		}
	}
	// Any other internal error.
	return &meh.Error{
		Code:       meh.ErrInternal,
		WrappedErr: meh.Wrap(detailedErr, "exec db query", details),
		Message:    message,
	}
}

// NewScanRowsErr creates a new meh.ErrInternal with the given error and message
// and saves the provided query to details.
func NewScanRowsErr(err error, message string, query string) error {
	return &meh.Error{
		Code:       meh.ErrInternal,
		WrappedErr: err,
		Message:    message,
		Details: meh.Details{
			"query": query,
		},
	}
}
