package blog_test

import (
	"errors"
	"go-brick/bcontext"
	"go-brick/berror"
	"go-brick/blog"
	"go-brick/btrace"
	"testing"
)

func TestLog(t *testing.T) {
	err := errors.New("first layer")
	err = berror.NewInternalError(err, "second layer", struct {
		TestField string
	}{
		TestField: "xxx",
	})
	err = berror.NewInternalError(err, "third layer")

	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Debug(ctx, "test debug message.")
	// {"level":"DEBUG","time":"2023-05-15T15:56:47+08:00","type":"BIZ","func":"go-brick/blog_test.TestLog","msg":"test debug message.","trace_id":"33bc0a072e5d4084aeee0dfcea36eaa6"}
	blog.Info(ctx, "test info message.")
	// {"level":"INFO","time":"2023-05-15T15:56:47+08:00","type":"BIZ","func":"go-brick/blog_test.TestLog","msg":"test info message.","trace_id":"33bc0a072e5d4084aeee0dfcea36eaa6"}
	blog.Warn(ctx, err, "test warn message.")
	// {"level":"WARN","time":"2023-05-15T15:56:47+08:00","type":"BIZ","func":"go-brick/blog_test.TestLog","msg":"test warn message.","trace_id":"33bc0a072e5d4084aeee0dfcea36eaa6","err":{"code":500,"reason":"third layer","next":{"code":500,"reason":"second layer","detail":{"TestField":"xxx"},"next":"first layer"}},"stack":[{"func":"go-brick/blog_test.TestLog","file":"D:/GitHub/go-brick/blog/log_test.go:20"},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go:1446"}]}
	blog.Error(ctx, err, "test error message.")
	// {"level":"ERROR","time":"2023-05-15T15:56:47+08:00","type":"BIZ","func":"go-brick/blog_test.TestLog","msg":"test error message.","trace_id":"33bc0a072e5d4084aeee0dfcea36eaa6","err":{"code":500,"reason":"third layer","next":{"code":500,"reason":"second layer","detail":{"TestField":"xxx"},"next":"first layer"}},"stack":[{"func":"go-brick/blog_test.TestLog","file":"D:/GitHub/go-brick/blog/log_test.go:20"},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go:1446"}]}
}

func TestLogFormat(t *testing.T) {
	type Name struct {
		FirstName string
		LastName  string
	}
	name := Name{FirstName: "Lamber", LastName: "Chen"}
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Debugf(ctx, "xxx: %+v, yyy: %d", name, 18)
	// {"level":"DEBUG","time":"2023-05-15T15:45:22+08:00","type":"BIZ","func":"go-brick/blog_test.TestLogFormat","msg":"xxx: {FirstName:Lamber LastName:Chen}, yyy: 18","trace_id":"9895452ea51642c9ad8f69dab5adcc86"}
}

func TestWithField(t *testing.T) {
	ctx := bcontext.New().Set(btrace.KeyTraceID, btrace.GenTraceID())
	blog.Debugw(ctx, "test with field",
		blog.Binary("b", []byte{97, 98}), // base64encode
		blog.ByteString("bstr", []byte{97, 98}),
		blog.ByteStrings("bstrs", [][]byte{{97}, {98}}),
		blog.Any("any", []byte{97, 98}))
	// {"level":"DEBUG","time":"2023-05-15T15:40:12+08:00","type":"BIZ","func":"go-brick/blog_test.TestWithField","msg":"test with field","trace_id":"29e38f4e13cc48a8826323adbc611073","b":"YWI=","bstr":"ab","bstrs":["a","b"],"any":"YWI="}
	blog.Infow(ctx, "test with field", blog.Any("any", string([]byte{97, 98})))
	// {"level":"INFO","time":"2023-05-15T15:40:12+08:00","type":"BIZ","func":"go-brick/blog_test.TestWithField","msg":"test with field","trace_id":"29e38f4e13cc48a8826323adbc611073","any":"ab"}
	type Name struct {
		FirstName string
		LastName  string
	}
	name := Name{FirstName: "Lamber", LastName: "Chen"}
	blog.Infow(ctx, "test with field", blog.Any("name", name))
	// {"level":"INFO","time":"2023-05-15T15:40:12+08:00","type":"BIZ","func":"go-brick/blog_test.TestWithField","msg":"test with field","trace_id":"29e38f4e13cc48a8826323adbc611073","name":{"FirstName":"Lamber","LastName":"Chen"}}
}
