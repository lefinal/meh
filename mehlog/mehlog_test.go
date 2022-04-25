package mehlog

import (
	"errors"
	"github.com/lefinal/meh/meh"
	"github.com/lefinal/zaprec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

// LogSuite tests Log.
type LogSuite struct {
	suite.Suite
}

// TestNilError assures that nil errors are logged, because of high possibility
// of bad usage.
func (suite *LogSuite) TestNilError() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, nil)
	suite.Len(rec.Records(), 1, "should be logged")
}

// TestNonMehError assures that non-meh errors are still logged.
func (suite *LogSuite) TestNonMehError() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, errors.New("sad life"))
	suite.Len(rec.RecordsByLevel(zapcore.ErrorLevel), 1, "should be logged")
}

// TestNotFoundError assures that errors with code meh.ErrNotFound are still
// logged to zapcore.ErrorLevel if not translated otherwise.
func (suite *LogSuite) TestNotFoundError() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, &meh.Error{Code: meh.ErrNotFound})
	suite.Len(rec.RecordsByLevel(zapcore.ErrorLevel), 1, "should be logged")
}

// TestNotFoundError assures that errors with code meh.ErrInternal are logged to
// zapcore.ErrorLevel if not translated otherwise.
func (suite *LogSuite) TestInternalError() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, &meh.Error{Code: meh.ErrInternal})
	suite.Len(rec.RecordsByLevel(zapcore.ErrorLevel), 1, "should be logged")
}

// TestNotFoundError assures that errors with code meh.ErrUnexpected are logged
// to zapcore.ErrorLevel if not translated otherwise.
func (suite *LogSuite) TestUnexpectedError() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, &meh.Error{})
	suite.Len(rec.RecordsByLevel(zapcore.ErrorLevel), 1, "should be logged")
}

// TestDetails assures that details are logged.
func (suite *LogSuite) TestDetails() {
	logger, rec := zaprec.NewRecorder(nil)
	Log(logger, meh.Wrap(&meh.Error{Details: meh.Details{"hello": "world"}}, "ola",
		meh.Details{"i_love": "cookies"}),
	)
	records := rec.Records()
	suite.Require().Len(records, 1, "should have been logged")
	suite.Contains(records[0].Fields, zap.Any("1/hello", "world"), "should contain details from root error")
	suite.Contains(records[0].Fields, zap.Any("0/i_love", "cookies"), "should contain details from root error")
}

// TestErrorMessage assures that the correct log message is used.
func (suite *LogSuite) TestErrorMessage() {
	logger, rec := zaprec.NewRecorder(nil)
	e := meh.Wrap(&meh.Error{Message: "inner"}, "outer", nil)
	Log(logger, e)
	records := rec.Records()
	suite.Require().Len(records, 1, "should have been logged")
	suite.Equal(e.Error(), records[0].Entry.Message, "should use error message as message")
}

func TestLog(t *testing.T) {
	suite.Run(t, new(LogSuite))
}

// TestWrapAndLog tests WrapAndLog.
func TestWrapAndLog(t *testing.T) {
	expected := meh.Wrap(&meh.Error{Message: "inner"}, "outer", nil).Error()
	logger, rec := zaprec.NewRecorder(nil)
	WrapAndLog(logger, &meh.Error{Message: "inner"}, "outer")
	records := rec.Records()
	require.Len(t, records, 1, "should have been logged")
	assert.Equal(t, expected, records[0].Entry.Message, "should have wrapped")
}

// logToLevelSuite tests logToLevel.
type logToLevelSuite struct {
	suite.Suite
}

// expect a logToLevel call with the given in-level to log an entry with given
// out-level.
func (suite logToLevelSuite) expect(levelIn zapcore.Level, levelOut zapcore.Level) {
	logger, rec := zaprec.NewRecorder(nil)
	logToLevel(logger, levelIn, "meow")
	suite.Len(rec.Records(), 1, "should have been logged")
	suite.Equal(levelOut, rec.Records()[0].Entry.Level, "should have logged to correct level")
}

func (suite *logToLevelSuite) TestDebugLevel() {
	suite.expect(zapcore.DebugLevel, zapcore.DebugLevel)
}

func (suite *logToLevelSuite) TestInfoLevel() {
	suite.expect(zapcore.InfoLevel, zapcore.InfoLevel)
}

func (suite *logToLevelSuite) TestWarnLevel() {
	suite.expect(zapcore.WarnLevel, zapcore.WarnLevel)
}

func (suite *logToLevelSuite) TestErrorLevel() {
	suite.expect(zapcore.ErrorLevel, zapcore.ErrorLevel)
}

func (suite *logToLevelSuite) TestDPanicLevel() {
	suite.expect(zapcore.DPanicLevel, zapcore.DPanicLevel)
}

func (suite *logToLevelSuite) TestPanicLevel() {
	suite.Panics(func() {
		suite.expect(zapcore.PanicLevel, zapcore.PanicLevel)
	}, "should panic")
}

// Cannot test fatal level.

func (suite *logToLevelSuite) TestUnknownLevel() {
	suite.expect(zapcore.Level(123), zapcore.ErrorLevel)
}

func Test_logToLevel(t *testing.T) {
	suite.Run(t, new(logToLevelSuite))
}
