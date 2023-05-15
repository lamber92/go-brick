package logger

import (
	"context"
	"fmt"
	"go-brick/berror"
	"go-brick/blog/config"
	"go-brick/bstack"
	"go-brick/btrace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// =======================================
// ----- Default Logger Engine IMPL ------
// =======================================

func newDefaultLogger(typ Type) *defaultLogger {
	// TODO：support load config file
	conf := config.NewDefault()
	core := zapcore.NewTee(
		zapcore.NewCore(
			conf.GetEncoder(),
			conf.GetWriterSyncer(),
			conf.GetLogLevel(conf.Level),
		),
	)
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		// zap.AddStacktrace(conf.GetLogLevel(conf.Stacktrace)),
	}
	if conf.Debug {
		options = append(options, zap.Development())
	}

	return &defaultLogger{
		loggerType: typ,
		engine:     zap.New(core, options...).Named(string(typ)),
	}
}

type defaultLogger struct {
	loggerType Type
	engine     *zap.Logger
}

// WithContext parse the built-in information of the infrastructure in the context into log.
func (d *defaultLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return d
	}
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	return &defaultLogger{
		engine: d.engine.With(zap.String("trace_id", btrace.GetTraceID(ctx))),
	}
}

// WithError parse the built-in information of the error into log.
func (d *defaultLogger) WithError(err error) Logger {
	if err == nil {
		return d
	}
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	switch tmp := err.(type) {
	case zapcore.ObjectMarshaler:
		return &defaultLogger{
			engine: d.engine.With(zap.Object("err", tmp)),
		}
	default:
		return &defaultLogger{
			engine: d.engine.With(zap.String("err", err.Error())),
		}
	}
}

// WithStack parse the built-in information of the stack into log.
func (d *defaultLogger) WithStack(source any) Logger {
	if source == nil {
		return d
	}
	var stack bstack.StackList
	switch tmp := source.(type) {
	case berror.Error:
		stack = tmp.Stack()
	default:
		stack = bstack.TakeStack(1, bstack.StacktraceMax)
	}
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	return &defaultLogger{
		engine: d.engine.With(zap.Array("stack", stack)),
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (d *defaultLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return d
	}
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	return &defaultLogger{
		engine: d.engine.With(defaultFields(fields).Release()...),
	}
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (d *defaultLogger) WithOptions(options ...Option) Logger {
	return &defaultLogger{
		engine: d.engine.WithOptions(defaultOptions(options).Release()...),
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (d *defaultLogger) Debug(msg string) {
	d.engine.Debug(msg)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (d *defaultLogger) Info(msg string) {
	d.engine.Info(msg)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (d *defaultLogger) Warn(msg string) {
	d.engine.Warn(msg)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (d *defaultLogger) Error(msg string) {
	d.engine.Error(msg)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (d *defaultLogger) Panic(msg string) {
	d.engine.Panic(msg)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (d *defaultLogger) Debugf(format string, a ...any) {
	d.engine.Debug(getMessage(format, a))
}

// Infof uses fmt.Sprintf to log a templated message.
func (d *defaultLogger) Infof(format string, a ...any) {
	d.engine.Info(getMessage(format, a))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (d *defaultLogger) Warnf(format string, a ...any) {
	d.engine.Warn(getMessage(format, a))
}

// Errorf uses fmt.Sprintf to log a templated message.
func (d *defaultLogger) Errorf(format string, a ...any) {
	d.engine.Error(getMessage(format, a))
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (d *defaultLogger) Panicf(format string, a ...any) {
	d.engine.Panic(getMessage(format, a))
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Debugw(msg string, fields ...Field) {
	d.engine.Debug(msg, defaultFields(fields).Release()...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Infow(msg string, fields ...Field) {
	d.engine.Info(msg, defaultFields(fields).Release()...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Warnw(msg string, fields ...Field) {
	d.engine.Warn(msg, defaultFields(fields).Release()...)

}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Errorw(msg string, fields ...Field) {
	d.engine.Error(msg, defaultFields(fields).Release()...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (d *defaultLogger) Panicw(msg string, fields ...Field) {
	d.engine.Panic(msg, defaultFields(fields).Release()...)
}

// Close close logger engine
func (d *defaultLogger) Close() error {
	if err := d.engine.Sync(); err != nil {
		return berror.NewInternalError(err, fmt.Sprintf("failed to close log engine [%s]", d.loggerType))
	}
	return nil
}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}
	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}
	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

// =======================================
// ------ Default Field Field IMPL -------
// =======================================

type defaultFields []Field

func (df defaultFields) Release() []zapcore.Field {
	out := make([]zap.Field, 0, len(df))
	for _, v := range df {
		out = append(out, v.(zapcore.Field))
	}
	return out
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return zap.Binary(key, val)
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, val [][]byte) Field {
	return zap.ByteStrings(key, val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value any) Field {
	return zap.Any(key, value)
}

// =======================================
// ------ Default Option Field IMPL ------
// =======================================

type defaultOptions []Option

func (do defaultOptions) Release() []zap.Option {
	out := make([]zap.Option, 0, len(do))
	for _, v := range do {
		out = append(out, v.(zap.Option))
	}
	return out
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option).
func AddCallerSkip(skip int) Option {
	return zap.AddCallerSkip(skip)
}
