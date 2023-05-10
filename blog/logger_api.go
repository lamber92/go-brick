package blog

import (
	"context"
)

type LoggerType string

const (
	TypeAccess LoggerType = "ACCESS"
	TypeBiz    LoggerType = "BIZ"
	TypeInfra  LoggerType = "INFRA"
)

type Logger interface {
	// WithContext parse the built-in information of the infrastructure in the context into log.
	WithContext(ctx context.Context) Logger
	// WithError parse the built-in information of the error into log.
	WithError(err error) Logger
	// WithStack parse the built-in information of the stack into log.
	WithStack(source any) Logger
	// With creates a child logger and adds structured context to it. Fields added
	// to the child don't affect the parent, and vice versa.
	With(fields ...Field) Logger

	// Debug logs a message at DebugLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Debug(msg string)
	// Info logs a message at InfoLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Info(msg string)
	// Warn logs a message at WarnLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Warn(msg string)
	// Error logs a message at ErrorLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Error(msg string)
	// Panic logs a message at PanicLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(msg string)

	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
	Panicf(format string, a ...any)

	// Debugw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Debugw(msg string, fields ...Field)
	Infow(msg string, fields ...Field)
	Warnw(msg string, fields ...Field)
	Errorw(msg string, fields ...Field)
	Panicw(msg string, fields ...Field)

	Close() error
}

type Field interface{}
