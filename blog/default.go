package blog

import (
	"context"
	"fmt"
	"go-brick/btrace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// =======================================
// ----- Default Logger Engine IMPL ------
// =======================================

func newDefaultLogger(typ LoggerType) *defaultLogger {
	// TODO：support load config file
	conf := &defaultConfig{
		Level:      "debug",
		Debug:      true,
		Stacktrace: "warn",
		Encoding:   "json",
		Output:     []string{"stdout"},
	}

	core := zapcore.NewTee(
		zapcore.NewCore(conf.getEncoder(), conf.getWriterSyncer(), conf.getZapLogLevel(conf.Level)),
	)
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(conf.getZapLogLevel(conf.Stacktrace)),
	}
	if conf.Debug {
		options = append(options, zap.Development())
	}

	return &defaultLogger{
		engine: zap.New(core, options...).
			With(zap.String("type", string(typ))),
	}
}

type defaultLogger struct {
	engine *zap.Logger
}

// WithContext parse the built-in information of the infrastructure in the context into log.
func (d *defaultLogger) WithContext(ctx context.Context) Logger {
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	return &defaultLogger{
		engine: d.engine.With(zap.String("trace_id", btrace.GetTraceID(ctx))),
	}
}

// WithError parse the built-in information of the error into log.
func (d *defaultLogger) WithError(err error) Logger {
	// TODO implement me
	panic("implement me")
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (d *defaultLogger) With(fields ...Field) Logger {
	// A new pointer object must be used to store the engine
	// to prevent polluting the original engine
	return &defaultLogger{
		engine: d.engine.With(defaultFields{fields}.Release()...),
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
	d.engine.Debug(msg, defaultFields{fields}.Release()...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Infow(msg string, fields ...Field) {
	d.engine.Info(msg, defaultFields{fields}.Release()...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Warnw(msg string, fields ...Field) {
	d.engine.Warn(msg, defaultFields{fields}.Release()...)

}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (d *defaultLogger) Errorw(msg string, fields ...Field) {
	d.engine.Error(msg, defaultFields{fields}.Release()...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (d *defaultLogger) Panicw(msg string, fields ...Field) {
	d.engine.Panic(msg, defaultFields{fields}.Release()...)
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
// ------ Default Logger Field IMPL ------
// =======================================

type defaultFields []Field

func (fs defaultFields) Release() []zap.Field {
	out := make([]zap.Field, 0, len(fs))
	for _, v := range fs {
		out = append(out, v.(zap.Field))
	}
	return out
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return zap.Any(key, val)
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
