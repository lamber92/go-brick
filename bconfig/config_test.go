package bconfig_test

import (
	"github.com/stretchr/testify/assert"
	"go-brick/bconfig"
	"go-brick/bconfig/bstorage"
	"go-brick/bcontext"
	"go-brick/btrace"
	"testing"
)

func TestYamlConfig(t *testing.T) {
	bconfig.Init(bconfig.Option{
		Type:      bstorage.YAML,
		ConfigDir: "./bstorage/yaml/config_test",
	})
	ctx := bcontext.New()
	v, err := bconfig.Dynamic.Load(ctx, "TestKey.E", "test")
	assert.Equal(t, nil, err)
	e1 := v.GetString("E1")
	assert.Equal(t, "iii", e1)

	trace, ok := btrace.GetMDFromCtx(ctx)
	assert.Equal(t, true, ok)
	t.Log(trace)
}
