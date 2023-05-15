package bstack_test

import (
	"go-brick/bstack"
	"testing"
)

func TestTakeStack(t *testing.T) {
	stack := bstack.TakeStack(0, bstack.StacktraceFull)
	t.Logf("%s", stack)
	// [{"func":"go-brick/bstack_test.TestTakeStack","file":"D:/GitHub/go-brick/bstack/stacktrace_test.go","line":9},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go","line":1446}]
}
