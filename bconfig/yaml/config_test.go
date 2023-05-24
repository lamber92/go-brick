package yaml_test

import (
	"go-brick/bconfig/yaml"
	"go-brick/bcontext"
	"go-brick/btrace"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStatic_GetFromKey(t *testing.T) {
	yaml.ReplaceRootDir("./config_test")
	ctx := bcontext.New()
	static := yaml.NewStatic()
	v, err := static.Load(ctx, "TestKey", "test")
	// double get from cache
	v, err = static.Load(ctx, "TestKey", "test")

	assert.Equal(t, nil, err)

	assert.Equal(t, "123456", v.GetString("A"))
	assert.Equal(t, 123456, v.GetInt("B"))
	assert.Equal(t, []string{"xxxxxxxxxx", "yyyyyyyyyy"}, v.GetStringSlice("C"))
	assert.Equal(t, time.Minute*3, v.GetDuration("D"))
	assert.Equal(t, "iii", v.Sub("E").GetString("E1"))
	assert.Equal(t, "jjj", v.Sub("E").GetString("E2"))

	// nb. A not-so-friendly detail from viper.Viper.
	// GetStringMap() will forcibly convert the key to lowercase.
	// https://github.com/spf13/viper/issues/1431
	e := v.GetStringMap("E")
	assert.Equal(t, "iii", e["e1"])
	assert.Equal(t, "jjj", e["e2"])
}

func TestNewStatic_ParseToStruct(t *testing.T) {
	yaml.ReplaceRootDir("./config_test")
	ctx := bcontext.New()
	static := yaml.NewStatic()
	v, err := static.Load(ctx, "TestKey", "test")

	type TestKey struct {
		Axx string        `mapstructure:"A"`
		Bxx int           `mapstructure:"B"`
		Cxx []string      `mapstructure:"C"`
		Dxx time.Duration `mapstructure:"D"`
		Exx struct {
			E1xx string `mapstructure:"E1"`
			E2xx string `mapstructure:"E2"`
		} `mapstructure:"E"`
	}
	testKey := TestKey{}
	err = v.Unmarshal(&testKey)

	assert.Equal(t, nil, err)

	assert.Equal(t, "123456", testKey.Axx)
	assert.Equal(t, 123456, testKey.Bxx)
	assert.Equal(t, []string{"xxxxxxxxxx", "yyyyyyyyyy"}, testKey.Cxx)
	assert.Equal(t, time.Minute*3, testKey.Dxx)
	assert.Equal(t, "iii", testKey.Exx.E1xx)
	assert.Equal(t, "jjj", testKey.Exx.E2xx)
}

func TestNewStatic_ParseMultiLayer(t *testing.T) {
	yaml.ReplaceRootDir("./config_test")
	ctx := bcontext.New()
	static := yaml.NewStatic()
	v, err := static.Load(ctx, "TestKey.E", "test")

	type TestKeyE struct {
		E1xx string `mapstructure:"E1"`
		E2xx string `mapstructure:"E2"`
	}
	testKey := TestKeyE{}
	err = v.Unmarshal(&testKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, "iii", testKey.E1xx)
	assert.Equal(t, "jjj", testKey.E2xx)
}

func TestNewStatic_GetTraceMD(t *testing.T) {
	yaml.ReplaceRootDir("./config_test")
	ctx := bcontext.New()
	static := yaml.NewStatic()
	_, _ = static.Load(ctx, "TestKey.E", "test")

	trace, ok := btrace.GetMDFromCtx(ctx)
	assert.Equal(t, true, ok)
	t.Log(trace)
}

func TestNewDynamic(t *testing.T) {
	yaml.ReplaceRootDir("./config_test")
	ctx := bcontext.New()
	dynamic := yaml.NewDynamic()
	v, err := dynamic.Load(ctx, "TestKey.E", "test")
	assert.Equal(t, nil, err)

	var (
		iii = "iii"
		kkk = "kkk"
	)

	origValue := v.GetString("E1")
	assert.Equal(t, iii, origValue)

	modifyConfig := func(old, new string) {
		filePath := "./config_test/dynamic/test.yaml" // 文件路径
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal(err)
		}
		newContent := strings.Replace(string(content), old, new, 1)
		if err = os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// change E1 value
	modifyConfig(iii, kkk)
	time.Sleep(time.Second)
	v, err = dynamic.Load(ctx, "TestKey.E", "test")
	// check E1 value is changed or not
	assert.Equal(t, nil, err)
	newValue := v.GetString("E1")
	assert.Equal(t, kkk, newValue)

	// restore E1 value
	modifyConfig(kkk, iii)
	time.Sleep(time.Second)
	v, err = dynamic.Load(ctx, "TestKey.E", "test")
	// check E1 value is restored or not
	assert.Equal(t, nil, err)
	newValue2 := v.GetString("E1")
	assert.Equal(t, iii, newValue2)

	// check load config tracing info
	trace, ok := btrace.GetMDFromCtx(ctx)
	assert.Equal(t, true, ok)
	t.Log(trace)
}
