package blog_test

import (
	"go-brick/bcontext"
	"go-brick/blog"
	"go-brick/btrace"
	"testing"
)

func TestAccessLog(t *testing.T) {
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Access.WithContext(ctx).Debug("test debug message.")
	blog.Access.WithContext(ctx).Info("test info message.")
	blog.Access.WithContext(ctx).Warn("test warn message.")
	blog.Access.WithContext(ctx).Error("test error message.")
	// blog.Access.WithContext(ctx).Panic("test panic message.")
}

func TestAccessPanicLog(t *testing.T) {
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Access.WithContext(ctx).Panic("test panic message.")
}
