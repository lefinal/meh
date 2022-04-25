package meh

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

// NewInternalErrFromErrSuite tests NewInternalErrFromErr.
type NewInternalErrFromErrSuite struct {
	suite.Suite
}

func (suite *NewInternalErrFromErrSuite) TestOK() {
	originalErr := errors.New("yo")
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewInternalErrFromErr(originalErr, message, details).(*Error)

	suite.Equal(ErrInternal, err.Code, "should have set correct code")
	suite.Equal(originalErr, err.WrappedErr, "should have applied the original error")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewInternalErrFromErr(t *testing.T) {
	suite.Run(t, new(NewInternalErrFromErrSuite))
}

// NewInternalErrSuite tests NewInternalErr.
type NewInternalErrSuite struct {
	suite.Suite
}

func (suite *NewInternalErrSuite) TestOK() {
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewInternalErr(message, details).(*Error)

	suite.Equal(ErrInternal, err.Code, "should have set correct error code")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewInternalErr(t *testing.T) {
	suite.Run(t, new(NewInternalErrSuite))
}

// NewBadInputErrFromErrSuite tests NewBadInputErrFromErr.
type NewBadInputErrFromErrSuite struct {
	suite.Suite
}

func (suite *NewBadInputErrFromErrSuite) TestOK() {
	originalErr := errors.New("yo")
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewBadInputErrFromErr(originalErr, message, details).(*Error)

	suite.Equal(ErrBadInput, err.Code, "should have set correct error code")
	suite.Equal(originalErr, err.WrappedErr, "should have applied the original error")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewBadInputErrFromErr(t *testing.T) {
	suite.Run(t, new(NewBadInputErrFromErrSuite))
}

// NewBadInputErrSuite tests NewBadInputErr.
type NewBadInputErrSuite struct {
	suite.Suite
}

func (suite *NewBadInputErrSuite) TestOK() {
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewBadInputErr(message, details).(*Error)

	suite.Equal(ErrBadInput, err.Code, "should have set correct error code")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewBadInputErr(t *testing.T) {
	suite.Run(t, new(NewBadInputErrSuite))
}

// NewNotFoundErrSuite tests NewNotFoundErr.
type NewNotFoundErrSuite struct {
	suite.Suite
}

func (suite *NewNotFoundErrSuite) TestOK() {
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewNotFoundErr(message, details).(*Error)

	suite.Equal(ErrNotFound, err.Code, "should have set correct error code")
	suite.Equal(message, err.Message, "should not have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewNotFoundErr(t *testing.T) {
	suite.Run(t, new(NewNotFoundErrSuite))
}

// NewNotFoundErrFromErrSuite tests NewNotFoundErrFromErr.
type NewNotFoundErrFromErrSuite struct {
	suite.Suite
}

func (suite *NewNotFoundErrFromErrSuite) TestOK() {
	originalErr := errors.New("yo")
	message := "Hello World!"
	details := Details{"hello": "world"}

	err := NewBadInputErrFromErr(originalErr, message, details).(*Error)

	suite.Equal(ErrBadInput, err.Code, "should have set correct error code")
	suite.Equal(originalErr, err.WrappedErr, "should have applied the original error")
	suite.Equal(message, err.Message, "should have applied message")
	suite.Equal(details, err.Details, "should have applied details")
}

func TestNewNotFoundErrFromErr(t *testing.T) {
	suite.Run(t, new(NewNotFoundErrFromErrSuite))
}
