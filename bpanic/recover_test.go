package bpanic_test

import (
	"context"
	"go-brick/blog"
	"go-brick/bpanic"
	"testing"
)

func TestRecover(t *testing.T) {
	defer bpanic.Recover(func(err error) {
		blog.Warn(context.Background(), err, "test recover")
	})
	f := func() {
		panic("xxx")
	}
	f()
	// {"level":"WARN","time":"2023-05-16T14:18:10+08:00","type":"BIZ","func":"go-brick/bpanic_test.TestRecover.func1","msg":"test recover","trace_id":"","err":{"code":500,"reason":"recover","detail":".(type)=string","next":"xxx"},"stack":[{"func":"go-brick/bpanic_test.TestRecover.func2","file":"D:/GitHub/go-brick/bpanic/recover_test.go:15"},{"func":"go-brick/bpanic_test.TestRecover","file":"D:/GitHub/go-brick/bpanic/recover_test.go:17"},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go:1446"}]}
}
