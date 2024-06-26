package yaml

import (
	"time"

	"github.com/lamber92/go-brick/bconfig/bstorage"
	"github.com/lamber92/go-brick/internal/json"
	"github.com/spf13/viper"
)

type defaultValue struct {
	data *viper.Viper
}

func newDefaultValue(v *viper.Viper) bstorage.Value {
	return &defaultValue{data: v}
}

func (d *defaultValue) Sub(key string) bstorage.Value {
	return &defaultValue{data: d.data.Sub(key)}
}

func (d *defaultValue) GetInt(key string) int {
	return d.data.GetInt(key)
}

func (d *defaultValue) GetUint(key string) uint {
	return d.data.GetUint(key)
}

func (d *defaultValue) GetString(key string) string {
	return d.data.GetString(key)
}

func (d *defaultValue) GetBool(key string) bool {
	return d.data.GetBool(key)
}

func (d *defaultValue) GetDuration(key string) time.Duration {
	return d.data.GetDuration(key)
}

func (d *defaultValue) GetIntSlice(key string) []int {
	return d.data.GetIntSlice(key)
}

func (d *defaultValue) GetStringSlice(key string) []string {
	return d.data.GetStringSlice(key)
}

func (d *defaultValue) GetStringMap(key string) map[string]any {
	return d.data.GetStringMap(key)
}

func (d *defaultValue) Unmarshal(rawVal any) error {
	return d.data.Unmarshal(rawVal)
}

func (d *defaultValue) String() string {
	if d.data == nil {
		return "<nil>"
	}
	tmp, _ := json.MarshalToString(d.data.AllSettings())
	return tmp
}
