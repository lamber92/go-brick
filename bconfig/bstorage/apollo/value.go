package apollo

import (
	"time"

	"github.com/lamber92/go-brick/bconfig/bstorage"
	"github.com/lamber92/go-brick/internal/json"
	"github.com/spf13/cast"
)

type defaultValue struct {
	data string
}

func newDefaultValue(v any) bstorage.Value {
	if s, ok := v.(string); ok {
		return &defaultValue{data: s}
	}
	return &defaultValue{data: cast.ToString(v)}
}

func (d *defaultValue) Sub(key string) bstorage.Value {
	return &defaultValue{data: json.Get([]byte(d.data), key).ToString()}
}

func (d *defaultValue) GetInt(key string) int {
	return json.Get([]byte(d.data), key).ToInt()
}

func (d *defaultValue) GetUint(key string) uint {
	return json.Get([]byte(d.data), key).ToUint()
}

func (d *defaultValue) GetString(key string) string {
	return json.Get([]byte(d.data), key).ToString()
}

func (d *defaultValue) GetBool(key string) bool {
	return json.Get([]byte(d.data), key).ToBool()
}

func (d *defaultValue) GetDuration(key string) time.Duration {
	return time.Duration(json.Get([]byte(d.data), key).ToInt64())
}

func (d *defaultValue) GetIntSlice(key string) []int {
	s := make([]int, 0)
	json.Get([]byte(d.data), key).ToVal(&s)
	return s
}

func (d *defaultValue) GetStringSlice(key string) []string {
	s := make([]string, 0)
	json.Get([]byte(d.data), key).ToVal(&s)
	return s
}

func (d *defaultValue) GetStringMap(key string) map[string]any {
	s := make(map[string]any)
	json.Get([]byte(d.data), key).ToVal(&s)
	return s
}

func (d *defaultValue) Unmarshal(rawVal any) error {
	return json.UnmarshalFromString(d.data, rawVal)
}

func (d *defaultValue) String() string {
	return d.data
}
