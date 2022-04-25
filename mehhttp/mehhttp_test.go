package mehhttp

import (
	"errors"
	"github.com/lefinal/meh/meh"
	"github.com/lefinal/zaprec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

// responseWriterStub mocks http.ResponseWriter.
type responseWriterStub struct {
	mock.Mock
}

func (s *responseWriterStub) Header() http.Header {
	return s.Called().Get(0).(http.Header)
}

func (s *responseWriterStub) Write(bytes []byte) (int, error) {
	args := s.Called(bytes)
	return args.Int(0), args.Error(1)
}

func (s *responseWriterStub) WriteHeader(statusCode int) {
	s.Called(statusCode)
}

// LogAndRespondErrorSuite tests LogAndRespondError.
type LogAndRespondErrorSuite struct {
	suite.Suite
	req    *http.Request
	rr     *httptest.ResponseRecorder
	logger *zap.Logger
	rec    *zaprec.RecordStore
}

func (suite *LogAndRespondErrorSuite) SetupTest() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	suite.Require().Nil(err, "creating request should not fail")
	suite.req = req
	suite.rr = httptest.NewRecorder()
	suite.logger, suite.rec = zaprec.NewRecorder(nil)
}

func (suite *LogAndRespondErrorSuite) TestLog() {
	LogAndRespondError(suite.logger, suite.rr, suite.req, &meh.Error{})
	suite.Len(suite.rec.Records(), 1, "should have been logged")
}

func (suite *LogAndRespondErrorSuite) TestAppliedDetails() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	suite.Require().Nil(err, "creating request should not fail")
	req.Header.Set("User-Agent", "007")
	req.RemoteAddr = "addr"
	req.Host = "host"
	LogAndRespondError(suite.logger, suite.rr, req, &meh.Error{})
	suite.Require().Len(suite.rec.Records(), 1, "should have been logged")
	record := suite.rec.Records()[0]
	suite.Contains(record.Fields, zap.Any("0/http_req_url", "http://localhost:8080"),
		"request url should have been added to details")
	suite.Contains(record.Fields, zap.Any("0/http_req_host", "host"),
		"request host should have been added to details")
	suite.Contains(record.Fields, zap.Any("0/http_req_method", "GET"),
		"request method should have been added to details")
	suite.Contains(record.Fields, zap.Any("0/http_req_user_agent", "007"),
		"request user should have been added to details")
	suite.Contains(record.Fields, zap.Any("0/http_req_remote_addr", "addr"),
		"request remote address should have been added to details")
}

// TestNotFoundError assures that meh.ErrNotFound is mapped to
// http.StatusInternalServerError per default.
func (suite *LogAndRespondErrorSuite) TestNotFoundError() {
	LogAndRespondError(suite.logger, suite.rr, suite.req, &meh.Error{
		Code:    meh.ErrNotFound,
		Message: "hidden",
	})
	suite.Equal(http.StatusInternalServerError, suite.rr.Code, "should return correct code")
	suite.NotContains(suite.rr.Body.String(), "hidden", "should mask error message")
}

// TestNotFoundError assures that meh.ErrBadInput is mapped to
// http.StatusInternalServerError per default.
func (suite *LogAndRespondErrorSuite) TestBadInputError() {
	LogAndRespondError(suite.logger, suite.rr, suite.req, &meh.Error{
		Code:    meh.ErrBadInput,
		Message: "hidden",
	})
	suite.Equal(http.StatusInternalServerError, suite.rr.Code, "should return correct code")
	suite.NotContains(suite.rr.Body.String(), "hidden", "should mask error message")
}

// TestNotFoundError assures that meh.ErrInternal is mapped to
// http.StatusInternalServerError per default.
func (suite *LogAndRespondErrorSuite) TestInternalError() {
	LogAndRespondError(suite.logger, suite.rr, suite.req, &meh.Error{
		Code:    meh.ErrInternal,
		Message: "hidden",
	})
	suite.Equal(http.StatusInternalServerError, suite.rr.Code, "should return correct code")
	suite.NotContains(suite.rr.Body.String(), "hidden", "should mask error message")
}

// TestUnexpectedError assures that for other errors always a
// http.StatusInternalServerError is responded.
func (suite *LogAndRespondErrorSuite) TestUnexpectedError() {
	LogAndRespondError(suite.logger, suite.rr, suite.req, &meh.Error{
		Message: "hidden",
	})
	suite.Equal(http.StatusInternalServerError, suite.rr.Code, "should return correct code")
	suite.NotContains(suite.rr.Body.String(), "hidden", "should mask error message")
}

// TestWriteFail assures that write errors are logged.
func (suite *LogAndRespondErrorSuite) TestWriteFail() {
	responseWriterStub := &responseWriterStub{}
	responseWriterStub.On("Header").Return(http.Header{})
	responseWriterStub.On("WriteHeader", mock.Anything)
	responseWriterStub.On("Write", mock.Anything).Return(0, errors.New("sad life"))
	defer responseWriterStub.AssertExpectations(suite.T())
	LogAndRespondError(suite.logger, responseWriterStub, suite.req, &meh.Error{
		Code: meh.ErrInternal,
	})
	suite.Len(suite.rec.Records(), 2, "should have logged both entries")
}

func Test_LogAndRespondError(t *testing.T) {
	suite.Run(t, new(LogAndRespondErrorSuite))
}

// TestSetHTTPStatusCodeMapping tests SetHTTPStatusCodeMapping.
func TestSetHTTPStatusCodeMapping(t *testing.T) {
	newMapping := func(_ meh.Code) int {
		return http.StatusTeapot
	}
	SetHTTPStatusCodeMapping(newMapping)
	req, err := http.NewRequest(http.MethodGet, "http://meow", nil)
	require.Nil(t, err, "create request should not fail")
	rr := httptest.NewRecorder()
	LogAndRespondError(zap.NewNop(), rr, req, &meh.Error{Code: meh.ErrInternal})
	assert.Equal(t, http.StatusTeapot, rr.Code, "should return correct code")
}
