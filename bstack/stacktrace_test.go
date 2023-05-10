package bstack_test

import (
	"go-brick/bstack"
	"testing"
)

func TestTakeStack(t *testing.T) {
	stack := bstack.TakeStack(0, bstack.StacktraceFull)
	t.Logf("%#v", stack)
}
