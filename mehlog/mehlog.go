// Package mehlog allows logging of meh.Error to zap.Logger.
package mehlog

import (
	"github.com/lefinal/meh"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var omitErrorMessageField = false
var omitErrorMessageFieldMutex sync.RWMutex

// OmitErrorMessageField sets whether the error message field with key
// meh.MapFieldErrorMessage should be omitted in logs and only output as log
// message in order to improve human readability.
func OmitErrorMessageField(omit bool) {
	omitErrorMessageFieldMutex.Lock()
	defer omitErrorMessageFieldMutex.Unlock()
	omitErrorMessageField = omit
}

var (
	// defaultLevelTranslator is the default LevelTranslator that translates every
	// meh.Code to zapcore.ErrorLevel.
	defaultLevelTranslator LevelTranslator = func(code meh.Code) zapcore.Level {
		return zapcore.ErrorLevel
	}
	// defaultLevelTranslatorMutex locks defaultLevelTranslator.
	defaultLevelTranslatorMutex sync.RWMutex
)

// SetDefaultLevelTranslator sets the LevelTranslator to be used for regular
// Log-calls.
func SetDefaultLevelTranslator(lt LevelTranslator) {
	defaultLevelTranslatorMutex.Lock()
	defer defaultLevelTranslatorMutex.Unlock()
	defaultLevelTranslator = lt
}

// LevelTranslator translates the given meh.Code to zapcore.Level for logging.
type LevelTranslator func(code meh.Code) zapcore.Level

// WrapAndLog calls Log after meh.Wrap with the given error and message.
func WrapAndLog(logger *zap.Logger, err error, message string) {
	Log(logger, meh.Wrap(err, message, nil))
}

// Log the given error using the default level translator that can be set via
// SetDefaultLevelTranslator.
func Log(logger *zap.Logger, err error) {
	e := meh.Cast(err)
	defaultLevelTranslatorMutex.RLock()
	level := defaultLevelTranslator(meh.ErrorCode(e))
	defaultLevelTranslatorMutex.RUnlock()
	LogToLevel(logger, level, err)
}

// LogToLevel logs the given error to the given zapcore.Level.
func LogToLevel(logger *zap.Logger, level zapcore.Level, err error) {
	e := meh.Cast(err)
	// Build fields.
	omitErrorMessageFieldMutex.RLock()
	omitErrorMessageField := omitErrorMessageField
	omitErrorMessageFieldMutex.RUnlock()
	fieldMap := meh.ToMap(e)
	fields := make([]zap.Field, 0, len(fieldMap))
	for k, v := range fieldMap {
		if omitErrorMessageField && k == meh.MapFieldErrorMessage {
			continue
		}
		fields = append(fields, zap.Any(k, v))
	}
	// Log it.
	logToLevel(logger, level, e.Error(), fields...)
}

// logToLevel calls the correct LogToLevel method for the given zap.Logger based on the
// zapcore.Level.
func logToLevel(logger *zap.Logger, level zapcore.Level, message string, fields ...zapcore.Field) {
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(message, fields...)
	case zapcore.InfoLevel:
		logger.Info(message, fields...)
	case zapcore.WarnLevel:
		logger.Warn(message, fields...)
	case zapcore.ErrorLevel:
		logger.Error(message, fields...)
	case zapcore.DPanicLevel:
		logger.DPanic(message, fields...)
	case zapcore.PanicLevel:
		logger.Panic(message, fields...)
	case zapcore.FatalLevel:
		logger.Fatal(message, fields...)
	default:
		logger.Error(message, fields...)
	}
}
