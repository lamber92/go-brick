package bconfig_test

import (
	"os"
	"testing"

	"github.com/lamber92/go-brick/bconfig"
	"github.com/lamber92/go-brick/bconfig/bstorage"
	"github.com/lamber92/go-brick/bcontext"
	"github.com/lamber92/go-brick/btrace"
	"github.com/stretchr/testify/assert"
)

func TestYamlConfig(t *testing.T) {
	bconfig.Init(bconfig.Option{
		Type:      bstorage.YAML,
		ConfigDir: "./bstorage/yaml/config_test",
	})
	ctx := bcontext.New()
	v, err := bconfig.Dynamic().Load(ctx, "TestKey.E", "test")
	if err != nil {
		t.Fatal(err)
	}
	e1 := v.GetString("E1")
	assert.Equal(t, "iii", e1)

	trace, ok := btrace.GetMDFromCtx(ctx)
	assert.Equal(t, true, ok)
	t.Log(trace)
}

func TestApolloConfig(t *testing.T) {
	if err := os.Setenv("GO_ENV_NAME", "dev"); err != nil {
		panic(err)
	}
	bconfig.Init(bconfig.Option{
		Type:      bstorage.APOLLO,
		ConfigDir: "./bstorage/apollo/config_test",
	})
	ctx := bcontext.New()
	v, err := bconfig.Dynamic().Load(ctx, "TestKey", "dev_apollo_config")
	if err != nil {
		t.Fatal(err)
	}
	e1 := v.Sub("E").GetString("E1")
	assert.Equal(t, "iii", e1)

	trace, ok := btrace.GetMDFromCtx(ctx)
	assert.Equal(t, true, ok)
	t.Log(trace)
}
