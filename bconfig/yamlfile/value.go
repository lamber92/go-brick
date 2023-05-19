package yamlfile

import (
	"go-brick/bconfig"
	"time"

	"github.com/spf13/viper"
)

type defaultValue struct {
	v *viper.Viper
}

func newDefaultValue(v *viper.Viper) bconfig.Value {
	return &defaultValue{v: v}
}

func (d *defaultValue) Sub(key string) bconfig.Value {
	return &defaultValue{v: d.v.Sub(key)}
}

func (d *defaultValue) GetInt(key string) int {
	return d.v.GetInt(key)
}

func (d *defaultValue) GetUint(key string) uint {
	return d.v.GetUint(key)
}

func (d *defaultValue) GetString(key string) string {
	return d.v.GetString(key)
}

func (d *defaultValue) GetBool(key string) bool {
	return d.v.GetBool(key)
}

func (d *defaultValue) GetDuration(key string) time.Duration {
	return d.v.GetDuration(key)
}

func (d *defaultValue) GetIntSlice(key string) []int {
	return d.v.GetIntSlice(key)
}

func (d *defaultValue) GetStringSlice(key string) []string {
	return d.v.GetStringSlice(key)
}

func (d *defaultValue) GetStringMap(key string) map[string]any {
	return d.v.GetStringMap(key)
}

func (d *defaultValue) Unmarshal(rawVal any) error {
	return d.v.Unmarshal(rawVal)
}
