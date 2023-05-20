package btrace_test

import (
	"context"
	"go-brick/blog"
	"go-brick/btrace"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testMod1 btrace.Module = "test_mod 1"
	testMod2 btrace.Module = "test_mod 2"
)

func TestDefaultMD(t *testing.T) {
	md := btrace.NewMD(testMod1, "test metadata 1")
	assert.Equal(t, testMod1, md.Module())
	assert.Contains(t, md.String(), "test metadata 1")
}

func TestDefaultChain(t *testing.T) {
	md1 := btrace.NewMD(testMod1, "test metadata 1")
	md2 := btrace.NewMD(testMod2, "test metadata 2")

	chain := btrace.NewChain()
	// append
	chain.Append(md1, md2)
	// get
	_ = chain.Get()
	mds1 := chain.Get()
	// clear
	chain.Clear()
	mds2 := chain.Get()

	assert.Equal(t, 2, len(mds1))
	assert.Equal(t, md1, mds1[0])
	assert.Equal(t, md2, mds1[1])
	assert.Equal(t, 0, len(mds2))

	// string
	chain.Append(md1, md2)
	t.Log(chain)
}

func TestMetadataList_MarshalLogArray(t *testing.T) {
	md1 := btrace.NewMD(testMod1, "test metadata 1")
	md2 := btrace.NewMD(testMod2, "test metadata 2")
	chain := btrace.NewChain()
	chain.Append(md1, md2)
	trace := chain.Get()

	blog.Infow(context.Background(), "test chain printout format", blog.Any("trace", trace))
	// {"level":"INFO","time":"2023-05-20T16:44:54+08:00","type":"BIZ","func":"go-brick/btrace_test.TestMetadataList_MarshalLogArray","msg":"test chain printout format","trace_id":"","trace":[{"module":"test_mod 1","value":"test metadata 1"},{"module":"test_mod 2","value":"test metadata 2"}]}
}
