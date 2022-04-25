package meh

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// newErrorUnwrapperSuite tests NewErrorUnwrapper.
type newErrorUnwrapperSuite struct {
	suite.Suite
}

func (suite *newErrorUnwrapperSuite) TestNilErr() {
	ew := NewErrorUnwrapper(nil)
	suite.False(ew.awaitFirstNext, "await first should be false")
	suite.Nil(ew.current, "current error should not be set")
}

func (suite *newErrorUnwrapperSuite) TestOK() {
	ew := NewErrorUnwrapper(&Error{})
	suite.True(ew.awaitFirstNext, "await first should be true")
	suite.NotNil(ew.current, "current error should be set")
}

func TestNewErrorWrapper(t *testing.T) {
	suite.Run(t, new(newErrorUnwrapperSuite))
}

// errorUnwrapperSuite tests ErrorUnwrapper.Current.
type errorUnwrapperCurrentSuite struct {
	suite.Suite
}

func (suite *errorUnwrapperCurrentSuite) TestAwaitFirstNext() {
	it := NewErrorUnwrapper(NewInternalErr("sad life", nil))
	suite.True(it.awaitFirstNext, "should await first next")
	suite.Nil(it.Current(), "should return nil error from current")
}

func (suite *errorUnwrapperCurrentSuite) TestOK() {
	e := NewInternalErr("sad life", nil)
	it := NewErrorUnwrapper(e)
	it.Next()
	suite.False(it.awaitFirstNext, "should not await next")
	suite.Equal(e, it.Current(), "should return expected error")
}

func TestErrorUnwrapper_Current(t *testing.T) {
	suite.Run(t, new(errorUnwrapperCurrentSuite))
}

// TestErrorUnwrapper_Next tests ErrorUnwrapper.Next.
func TestErrorUnwrapper_Next(t *testing.T) {
	originalErr := errors.New("sad life")
	inner := Wrap(originalErr, "inner", nil)
	outer := Wrap(inner, "outer", nil)
	uw := NewErrorUnwrapper(outer)
	assert.Nil(t, uw.Current(), "should return nil error before first next call")
	assert.True(t, uw.Next(), "should have errors after first next call")
	assert.Equal(t, outer, uw.Current(), "should return outer error after first next")
	assert.True(t, uw.Next(), "should have errors after second next call")
	assert.Equal(t, inner, uw.Current(), "should return inner error after second next")
	assert.True(t, uw.Next(), "should have errors after third next call")
	assert.Equal(t, originalErr, uw.Current(), "should return original error after third next")
	assert.False(t, uw.Next(), "should have no more errors after fourth next call")
}
