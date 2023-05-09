package blog

import (
	"context"
)

func Debug(ctx context.Context, msg string) {
	Biz.WithContext(ctx).Debug(msg)
}

func Info(ctx context.Context, msg string) {
	Biz.WithContext(ctx).Info(msg)
}

func Warn(ctx context.Context, err error, msg string) {
	Biz.WithContext(ctx).WithError(err).Warn(msg)
}

func Error(ctx context.Context, err error, msg string) {
	Biz.WithContext(ctx).WithError(err).Error(msg)
}

func Panic(ctx context.Context, err error, msg string) {
	Biz.WithContext(ctx).WithError(err).Panic(msg)
}

func Debugf(ctx context.Context, format string, a ...any) {
	Biz.WithContext(ctx).Debugf(format, a)
}

func Infof(ctx context.Context, format string, a ...any) {
	Biz.WithContext(ctx).Infof(format, a)
}

func Warnf(ctx context.Context, err error, format string, a ...any) {
	Biz.WithContext(ctx).WithError(err).Warnf(format, a)
}

func Errorf(ctx context.Context, err error, format string, a ...any) {
	Biz.WithContext(ctx).WithError(err).Errorf(format, a)
}

func Panicf(ctx context.Context, err error, format string, a ...any) {
	Biz.WithContext(ctx).WithError(err).Panicf(format, a)
}

func Debugw(ctx context.Context, msg string, fields ...Field) {
	Biz.WithContext(ctx).Debugw(msg, fields)
}

func Infow(ctx context.Context, msg string, fields ...Field) {
	Biz.WithContext(ctx).Infow(msg, fields)
}

func Warnw(ctx context.Context, err error, msg string, fields ...Field) {
	Biz.WithContext(ctx).WithError(err).Warnw(msg, fields)
}

func Errorw(ctx context.Context, err error, msg string, fields ...Field) {
	Biz.WithContext(ctx).WithError(err).Errorw(msg, fields)
}

func Panicw(ctx context.Context, err error, msg string, fields ...Field) {
	Biz.WithContext(ctx).WithError(err).Panicw(msg, fields)
}
