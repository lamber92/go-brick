package blog_test

import (
	"errors"
	"go-brick/bcontext"
	"go-brick/berror"
	"go-brick/blog"
	"go-brick/btrace"
	"testing"
)

func TestAccessLog(t *testing.T) {
	err := errors.New("first layer")
	err = berror.NewInternalError(err, "second layer", struct {
		TestField string
	}{
		TestField: "xxx",
	})
	err = berror.NewInternalError(err, "third layer")

	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Access.WithContext(ctx).Debug("test debug message.")
	blog.Access.WithContext(ctx).Info("test info message.")
	blog.Access.WithContext(ctx).WithError(err).WithStack(err).Warn("test warn message.")
	blog.Access.WithContext(ctx).WithError(err).WithStack(err).Error("test error message.")
	// blog.Access.WithContext(ctx).Panic("test panic message.")
}

func TestAccessPanicLog(t *testing.T) {
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Access.WithContext(ctx).Panic("test panic message.")
}
