package meh

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// ErrorErrorSuite tests Error.Error.
type ErrorErrorSuite struct {
	suite.Suite
}

func (suite *ErrorErrorSuite) TestSingleWithoutWrapped() {
	e := &Error{Message: "woof"}
	suite.Equal("woof", e.Error(), "should return correct message")
}

func (suite *ErrorErrorSuite) TestWrappers() {
	e := &Error{
		Message: "outer",
		WrappedErr: &Error{
			Message:    "inner",
			WrappedErr: errors.New("original"),
		},
	}
	suite.Equal("outer: inner: original", e.Error(), "should return correct message")
}

func TestError_Error(t *testing.T) {
	suite.Run(t, new(ErrorErrorSuite))
}

// TestErrorCode tests ErrorCode.
func TestErrorCode(t *testing.T) {
	e := &Error{
		Code: ErrNeutral,
		WrappedErr: &Error{
			Code: ErrNeutral,
			WrappedErr: &Error{
				Code: ErrInternal,
				WrappedErr: &Error{
					Code: ErrBadInput,
					WrappedErr: &Error{
						Code: ErrNeutral,
						WrappedErr: &Error{
							Code: ErrNotFound,
						},
					},
				},
			},
		},
	}
	assert.Equal(t, ErrInternal, ErrorCode(e), "should return correct code")
}

// WrapSuite tests Wrap.
type WrapSuite struct {
	suite.Suite
}

func (suite *WrapSuite) TestEmptyMessage() {
	d := Details{"Hello": "World!"}
	originalErr := &Error{}
	e := Wrap(originalErr, "", d).(*Error)
	suite.Equal(ErrNeutral, e.Code, "should have set correct code")
	suite.Empty(e.Message, "should have have set empty message")
	suite.Equal(d, e.Details, "should have set details")
	suite.Equal(originalErr, e.WrappedErr, "should have kept original error")
}

func (suite *WrapSuite) TestNilDetails() {
	originalErr := &Error{Details: Details{"woof": "meow"}}
	e := Wrap(originalErr, "baa", nil).(*Error)
	suite.Equal(ErrNeutral, e.Code, "should have set correct code")
	suite.Equal("baa", e.Message, "should have set correct message")
	suite.Equal(originalErr, e.WrappedErr, "should have kept original error")
}

func TestWrap(t *testing.T) {
	suite.Run(t, new(WrapSuite))
}

// ApplyDetailsSuite tests ApplyDetails.
type ApplyDetailsSuite struct {
	suite.Suite
}

func (suite *ApplyDetailsSuite) TestOldDetailsNil() {
	originalErr := &Error{}
	e := ApplyDetails(originalErr, Details{"hello": "world"}).(*Error)
	suite.Equal(originalErr, e.WrappedErr, "should have kept error")
	suite.Require().NotEmpty(e.Details, "should have set details")
	suite.Equal("world", e.Details["hello"], "should have set correct details")
}

func (suite *ApplyDetailsSuite) TestNewDetailsNil() {
	originalErr := &Error{Details: Details{"hello": "world"}}
	e := ApplyDetails(originalErr, nil).(*Error)
	suite.Equal(originalErr, e.WrappedErr, "should have kept original error")
	suite.Empty(e.Details, "should not have set details")
}

func TestApplyDetails(t *testing.T) {
	suite.Run(t, new(ApplyDetailsSuite))
}

// CastSuite tests Cast.
type CastSuite struct {
	suite.Suite
}

func (suite *CastSuite) TestOK() {
	e := &Error{
		Code:       ErrInternal,
		WrappedErr: errors.New("moo"),
		Message:    "Hello World!",
		Details:    Details{"oink": "cluck"},
	}
	suite.Equal(e, Cast(e), "should cast correctly")
}

func (suite *CastSuite) TestNoMehError() {
	err := errors.New("chirp")
	e := Cast(err)
	suite.IsType(&Error{}, e, "should have created error with correct type")
	suite.Equal(fromErr(err), e, "should have created a new meh error from given")
}

func TestCast(t *testing.T) {
	suite.Run(t, new(CastSuite))
}

// fromErrSuite tests fromErr.
type fromErrSuite struct {
	suite.Suite
}

func (suite *fromErrSuite) TestNilError() {
	e := fromErr(nil)
	suite.Equal(ErrUnexpected, e.Code, "should have set correct code")
	suite.NotEmpty(e.Message, "should not have empty message")
}

func (suite *fromErrSuite) TestOK() {
	originalErr := errors.New("chirp")
	e := fromErr(originalErr)
	suite.Equal(ErrUnexpected, e.Code, "should have set correct code")
	suite.Equal(originalErr, e.WrappedErr, "should have wrapped original error")
	suite.Empty(e.Message, "should not have applied any error message")
}

func Test_fromErr(t *testing.T) {
	suite.Run(t, new(fromErrSuite))
}

// ToMapSuite tests ToMap.
type ToMapSuite struct {
	suite.Suite
}

func (suite *ToMapSuite) TestOK() {
	d := Details{"roar": "meow"}
	e := &Error{
		Code:       ErrInternal,
		WrappedErr: errors.New("woof"),
		Message:    "chirp",
		Details:    d,
	}
	ed := ToMap(e)
	suite.Equal(ErrInternal, ed[MapFieldErrorCode], "should have set correct code")
	suite.Equal(e.Error(), ed[MapFieldErrorMessage], "should have set correct error message")
	suite.Equal("meow", ed["0/roar"], "should have kept old details")
	// Assure copied and not set in place.
	_, ok := d[MapFieldErrorCode]
	suite.False(ok, "should not have touched original details")
}

func (suite *ToMapSuite) TestNoMehError() {
	err := errors.New("woof")
	ed := ToMap(err)
	suite.Equal(ErrUnexpected, ed[MapFieldErrorCode], "should set error code correctly")
	suite.Equal(err.Error(), ed[MapFieldErrorMessage], "should set error message correctly")
}

func (suite *ToMapSuite) TestDuplicateDetailKeys() {
	e := &Error{
		Code: ErrInternal,
		Details: Details{
			"detail_0":         "meow",
			"detail_to_mask_1": "woof",
		},
		WrappedErr: &Error{
			Code: ErrNeutral,
			Details: Details{
				"detail_to_mask_2": "chirp",
				"detail_to_mask_1": "ola",
				"detail_3":         "wow",
				"x_masked":         []string{"a", "b", "c"},
			},
			WrappedErr: &Error{
				Code: ErrNotFound,
				Details: Details{
					"detail_to_mask_2": "oh no",
					"detail_to_mask_1": "cluck",
				},
			},
		},
	}
	suite.Equal(map[string]interface{}{
		"0/detail_0":         "meow",
		"0/detail_to_mask_1": "woof",
		"1/detail_to_mask_2": "chirp",
		"1/detail_to_mask_1": "ola",
		"1/detail_3":         "wow",
		"1/x_masked":         []string{"a", "b", "c"},
		"2/detail_to_mask_2": "oh no",
		"2/detail_to_mask_1": "cluck",
		MapFieldErrorMessage: "",
		MapFieldErrorCode:    ErrInternal,
	}, ToMap(e), "should set correct details")
}

func TestToMap(t *testing.T) {
	suite.Run(t, new(ToMapSuite))
}

// TestApplyCode tests ApplyCode.
func TestApplyCode(t *testing.T) {
	originalErr := &Error{
		Code:    ErrNotFound,
		Message: "sad life",
	}
	e := ApplyCode(originalErr, ErrInternal).(*Error)
	assert.Equal(t, originalErr, e.WrappedErr, "should have kept original error")
	assert.Equal(t, ErrInternal, ErrorCode(e), "should have applied correct error code")
}

// NilOrWrapSuite tests NilOrWrap.
type NilOrWrapSuite struct {
	suite.Suite
}

func (suite *NilOrWrapSuite) TestNil() {
	e := NilOrWrap(nil, "meow", nil)
	suite.Nil(e, "should return nil")
}

func (suite *NilOrWrapSuite) TestNotNil() {
	originalErr := NewInternalErr("sad life", Details{"meow": "woof"})
	e := NilOrWrap(originalErr, "hello", Details{"hello": "world"})
	suite.Equal(Wrap(originalErr, "hello", Details{"hello": "world"}), e, "should have been wrapped")
}

func TestNilOrWrap(t *testing.T) {
	suite.Run(t, new(NilOrWrapSuite))
}
