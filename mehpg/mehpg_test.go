package mehpg

import (
	"errors"
	"github.com/lefinal/meh"
	"github.com/stretchr/testify/suite"
	"testing"
)

// NewQueryDBErrSuite tests NewQueryDBErr.
type NewQueryDBErrSuite struct {
	suite.Suite
}

func (suite *NewQueryDBErrSuite) TestOK() {
	originalErr := errors.New("yo")
	message := "Hello World!"
	query := "SELECT *"
	err := NewQueryDBErr(originalErr, message, query).(*meh.Error)
	wrapped := meh.Cast(err).WrappedErr
	suite.Equal(meh.ErrInternal, meh.ErrorCode(err), "should have set error code to internal")
	suite.Equal(originalErr, wrapped, "should have applied the original error")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(query, err.Details["query"], "should have applied the query to details")
}

func TestNewQueryDBErr(t *testing.T) {
	suite.Run(t, new(NewQueryDBErrSuite))
}

// NewScanRowsErrSuite tests NewScanRowsErr.
type NewScanRowsErrSuite struct {
	suite.Suite
}

func (suite *NewScanRowsErrSuite) TestOK() {
	originalErr := errors.New("yo")
	message := "Hello World!"
	query := "SELECT *"
	err := NewScanRowsErr(originalErr, message, query).(*meh.Error)
	suite.Equal(meh.ErrInternal, err.Code, "should have set error code to internal")
	suite.Equal(originalErr, err.WrappedErr, "should have applied the original error")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(query, err.Details["query"], "should have applied the query to details")
}

func TestNewScanRowsErr(t *testing.T) {
	suite.Run(t, new(NewScanRowsErrSuite))
}
