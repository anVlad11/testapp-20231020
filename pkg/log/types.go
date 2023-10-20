package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Err   Level = "error"
)

var (
	levelToZap = map[Level]zapcore.Level{
		Debug: zap.DebugLevel,
		Info:  zap.InfoLevel,
		Warn:  zap.WarnLevel,
		Err:   zap.ErrorLevel,
	}
	defaultLevel = zap.InfoLevel
)

type Level string

type Option func(*ZapLogger)

type ZapLogger struct {
	zapLogger *zap.SugaredLogger
}

// NewTestLogger return instance of Logger that discards all output.
func NewTestLogger() Logger {
	return &ZapLogger{
		zapLogger: zap.NewNop().Sugar(),
	}
}

type Logger interface {
	Info(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	AddField(name, value string) Logger
}
