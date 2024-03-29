package blog

import (
	"context"

	"github.com/lamber92/go-brick/blog/logger"
)

var _biz = logger.Biz.WithOptions(logger.AddCallerSkip(1))

func Debug(ctx context.Context, msg string) {
	_biz.WithContext(ctx).Debug(msg)
}

func Info(ctx context.Context, msg string) {
	_biz.WithContext(ctx).Info(msg)
}

func Warn(ctx context.Context, err error, msg string) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Warn(msg)
}

func Error(ctx context.Context, err error, msg string) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Error(msg)
}

func Panic(ctx context.Context, err error, msg string) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Panic(msg)
}

func Debugf(ctx context.Context, format string, a ...any) {
	_biz.WithContext(ctx).Debugf(format, a...)
}

func Infof(ctx context.Context, format string, a ...any) {
	_biz.WithContext(ctx).Infof(format, a...)
}

func Warnf(ctx context.Context, err error, format string, a ...any) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Warnf(format, a...)
}

func Errorf(ctx context.Context, err error, format string, a ...any) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Errorf(format, a...)
}

func Panicf(ctx context.Context, err error, format string, a ...any) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Panicf(format, a...)
}

func Debugw(ctx context.Context, msg string, fields ...logger.Field) {
	_biz.WithContext(ctx).Debugw(msg, fields...)
}

func Infow(ctx context.Context, msg string, fields ...logger.Field) {
	_biz.WithContext(ctx).Infow(msg, fields...)
}

func Warnw(ctx context.Context, err error, msg string, fields ...logger.Field) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Warnw(msg, fields...)
}

func Errorw(ctx context.Context, err error, msg string, fields ...logger.Field) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Errorw(msg, fields...)
}

func Panicw(ctx context.Context, err error, msg string, fields ...logger.Field) {
	_biz.WithContext(ctx).WithError(err).WithStack(err).Panicw(msg, fields...)
}
