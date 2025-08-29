package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
}

// Field represents a log field
type Field struct {
	Key   string
	Value interface{}
}

// ZapLogger wraps zap logger
type ZapLogger struct {
	logger *zap.Logger
}

// LogrusLogger wraps logrus logger
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewZapLogger creates a new zap logger instance
func NewZapLogger(level string, format string) (Logger, error) {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: logger}, nil
}

// NewLogrusLogger creates a new logrus logger instance
func NewLogrusLogger(level string, format string) Logger {
	logger := logrus.New()

	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set format
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	}

	logger.SetOutput(os.Stdout)

	return &LogrusLogger{logger: logger}
}

// ZapLogger methods
func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.zapFields(fields...)...)
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.zapFields(fields...)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.zapFields(fields...)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, l.zapFields(fields...)...)
}

func (l *ZapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, l.zapFields(fields...)...)
}

func (l *ZapLogger) WithFields(fields ...Field) Logger {
	return &ZapLogger{logger: l.logger.With(l.zapFields(fields...)...)}
}

func (l *ZapLogger) zapFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// LogrusLogger methods
func (l *LogrusLogger) Debug(msg string, fields ...Field) {
	l.logger.WithFields(l.logrusFields(fields...)).Debug(msg)
}

func (l *LogrusLogger) Info(msg string, fields ...Field) {
	l.logger.WithFields(l.logrusFields(fields...)).Info(msg)
}

func (l *LogrusLogger) Warn(msg string, fields ...Field) {
	l.logger.WithFields(l.logrusFields(fields...)).Warn(msg)
}

func (l *LogrusLogger) Error(msg string, fields ...Field) {
	l.logger.WithFields(l.logrusFields(fields...)).Error(msg)
}

func (l *LogrusLogger) Fatal(msg string, fields ...Field) {
	l.logger.WithFields(l.logrusFields(fields...)).Fatal(msg)
}

func (l *LogrusLogger) WithFields(fields ...Field) Logger {
	return &LogrusLogger{logger: l.logger.WithFields(l.logrusFields(fields...)).Logger}
}

func (l *LogrusLogger) logrusFields(fields ...Field) logrus.Fields {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}

// SetOutput sets the output destination for logrus logger
func (l *LogrusLogger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

// F is a convenience function to create a Field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
