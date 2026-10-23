package meh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// stdWrappingSuite tests that meh errors keep working when wrapped the
// standard way, e.g. by libraries using fmt.Errorf with %w or errors.Join.
type stdWrappingSuite struct {
	suite.Suite
}

func (suite *stdWrappingSuite) TestErrorCodeThroughFmtWrap() {
	err := Wrap(fmt.Errorf("call: %w", NewBadInputErr("bad", nil)), "outer", nil)
	suite.Equal(ErrBadInput, ErrorCode(err))
}

// joinedErr wraps several errors like errors.Join, which needs Go 1.20.
type joinedErr struct {
	errs []error
}

func (e joinedErr) Error() string {
	return "joined"
}

func (e joinedErr) Unwrap() []error {
	return e.errs
}

func (suite *stdWrappingSuite) TestErrorCodeThroughJoin() {
	joined := joinedErr{errs: []error{NewNotFoundErr("missing", nil), errors.New("rollback failed")}}
	suite.Equal(ErrNotFound, ErrorCode(Wrap(joined, "outer", nil)))
}

func (suite *stdWrappingSuite) TestErrorCodePlainErrorStaysUnexpected() {
	suite.Equal(ErrUnexpected, ErrorCode(errors.New("plain")))
	suite.Equal(ErrUnexpected, ErrorCode(Wrap(errors.New("plain"), "outer", nil)))
	suite.Equal(ErrUnexpected, ErrorCode(fmt.Errorf("call: %w", errors.New("plain"))))
}

func (suite *stdWrappingSuite) TestErrorCodeNeutralStaysNeutral() {
	suite.Equal(ErrNeutral, ErrorCode(Wrap(&Error{Code: ErrNeutral}, "outer", nil)))
}

func (suite *stdWrappingSuite) TestErrorMessageNotRepeated() {
	err := Wrap(fmt.Errorf("call: %w", NewBadInputErr("bad", nil)), "outer", nil)
	suite.Equal("outer: call: bad", err.Error())
}

func (suite *stdWrappingSuite) TestErrorsAsThroughMehErrors() {
	inner := &Error{Code: ErrNotFound, Message: "missing"}
	err := Wrap(fmt.Errorf("call: %w", Wrap(inner, "middle", nil)), "outer", nil)
	var found *Error
	suite.Require().True(errors.As(fmt.Errorf("top: %w", err), &found))
	suite.True(errors.Is(err, inner))
}

func (suite *stdWrappingSuite) TestErrorsIsSentinelThroughMehErrors() {
	sentinel := errors.New("sentinel")
	err := Wrap(Wrap(sentinel, "inner", nil), "outer", nil)
	suite.True(errors.Is(err, sentinel))
}

func (suite *stdWrappingSuite) TestUnwrapperVisitsStdWrappedLevels() {
	inner := NewBadInputErr("bad", nil)
	wrapped := fmt.Errorf("call: %w", inner)
	err := Wrap(wrapped, "outer", nil)
	levels := make([]error, 0)
	for it := NewErrorUnwrapper(err); it.Next(); {
		levels = append(levels, it.Current())
	}
	suite.Equal([]error{err, wrapped, inner}, levels)
}

func TestStdWrapping(t *testing.T) {
	suite.Run(t, new(stdWrappingSuite))
}
