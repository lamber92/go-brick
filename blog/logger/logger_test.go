package logger_test

import (
	"errors"
	"go-brick/bcontext"
	"go-brick/berror"
	"go-brick/blog/logger"
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
	logger.Access.WithContext(ctx).Debug("test debug message.")
	logger.Access.WithContext(ctx).Info("test info message.")
	logger.Access.WithContext(ctx).WithError(err).WithStack(err).Warn("test warn message.")
	logger.Access.WithContext(ctx).WithError(err).WithStack(err).Error("test error message.")
	// blog.Access.WithContext(ctx).Panic("test panic message.")
}

func TestAccessPanicLog(t *testing.T) {
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	logger.Access.WithContext(ctx).Panic("test panic message.")
}

func TestLogFormat(t *testing.T) {
	type Name struct {
		FirstName string
		LastName  string
	}
	name := Name{FirstName: "Lamber", LastName: "Chen"}
	logger.Infra.Debugf("xxx: %+v, yyy: %d", name, 18)
	// {"level":"DEBUG","time":"2023-05-15T15:41:49+08:00","type":"INFRA","func":"go-brick/blog/logger_test.TestLogFormat","msg":"xxx: {FirstName:Lamber LastName:Chen}, yyy: 18"}
}

func TestWithField(t *testing.T) {
	logger.Biz.Debugw("test with field",
		logger.NewField().Binary("b", []byte{97, 98}), // base64encode
		logger.NewField().ByteString("bstr", []byte{97, 98}),
		logger.NewField().ByteStrings("bstrs", [][]byte{{97}, {98}}),
		logger.NewField().Any("any", []byte{97, 98}))
	// {"level":"DEBUG","time":"2023-05-15T15:39:40+08:00","type":"BIZ","func":"go-brick/blog/logger_test.TestWithField","msg":"test with field","b":"YWI=","bstr":"ab","bstrs":["a","b"],"any":"YWI="}
	logger.Biz.Infow("test with field", logger.NewField().Any("any", string([]byte{97, 98})))
	// {"level":"INFO","time":"2023-05-15T15:39:40+08:00","type":"BIZ","func":"go-brick/blog/logger_test.TestWithField","msg":"test with field","any":"ab"}
	type Name struct {
		FirstName string
		LastName  string
	}
	name := Name{FirstName: "Lamber", LastName: "Chen"}
	logger.Biz.Infow("test with field", logger.NewField().Any("name", name))
	// {"level":"INFO","time":"2023-05-15T15:39:40+08:00","type":"BIZ","func":"go-brick/blog/logger_test.TestWithField","msg":"test with field","name":{"FirstName":"Lamber","LastName":"Chen"}}

	logger.Biz.Infow("test with empty field", logger.NewField())
	// {"level":"INFO","time":"2023-05-15T15:39:40+08:00","type":"BIZ","func":"go-brick/blog/logger_test.TestWithField","msg":"test with empty field"}
}
