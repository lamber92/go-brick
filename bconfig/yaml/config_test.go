package yaml_test

import (
	"go-brick/bconfig/yaml"
	"go-brick/bcontext"
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

func TestNewDynamic(t *testing.T) {

}
