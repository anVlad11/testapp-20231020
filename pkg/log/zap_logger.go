package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(service string, logLevel Level) (*ZapLogger, error) {
	zapLevel := getZapLevel(logLevel)

	cfg := newZapConfig(service, zapLevel)
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		zapLogger: zl.Sugar(),
	}, nil
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Infow(msg, fieldsToInterface(fields)...)
}

func (l *ZapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panicw(msg, fieldsToInterface(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Errorw(msg, fieldsToInterface(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debugw(msg, fieldsToInterface(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warnw(msg, fieldsToInterface(fields)...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.zapLogger.Infof(template, args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.zapLogger.Errorf(template, args...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.zapLogger.Debugf(template, args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.zapLogger.Warnf(template, args...)
}

func (l *ZapLogger) AddField(name, value string) Logger {
	return &ZapLogger{
		zapLogger: l.zapLogger.With(name, value),
	}
}

func newZapConfig(service string, level zapcore.Level) zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "lvl",
			TimeKey:        "ts",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochNanosTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"source": service,
		},
	}
}

func getZapLevel(l Level) zapcore.Level {
	out, ok := levelToZap[l]
	if !ok {
		return defaultLevel
	}
	return out
}
