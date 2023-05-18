package yamlfile

import (
	"go-brick/bconfig"
	"time"

	"github.com/spf13/viper"
)

type defaultValue struct {
	v *viper.Viper
}

func (d defaultValue) Sub(key string) bconfig.Value {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetInt(key string) int {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetUint(key string) uint {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetString(key string) string {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetBool(key string) bool {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetDuration(key string) time.Duration {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetIntSlice(key string) []int {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetStringSlice(key string) []string {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) GetStringMap(key string) map[string]any {
	// TODO implement me
	panic("implement me")
}

func (d defaultValue) Unmarshal(rawVal any) error {
	// TODO implement me
	panic("implement me")
}
