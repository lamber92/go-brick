package yamlfile

import (
	"context"
	"fmt"
	"go-brick/bconfig"
	"go-brick/berror"
	"sync"

	"github.com/spf13/viper"
)

const defaultNamespace = "config"

// TODO: 加载自定义配置文件目录路径

var _config = newDefault()

func newDefault() *yamlFile {
	return &yamlFile{
		static:  sync.Map{},
		dynamic: sync.Map{},
	}
}

type yamlFile struct {
	static  sync.Map
	dynamic sync.Map
}

func (f *yamlFile) StaticLoad(key string, namespace ...string) (bconfig.Value, error) {
	var filename = defaultNamespace
	if len(namespace) > 0 {
		filename = namespace[0]
	}

	conf, ok := f.static.Load(key)
	if ok {
		return conf.(bconfig.Value), nil
	}
	v, err := f.loadFile("", filename)
	if err != nil {
		return nil, err
	}
	v = v.Sub(key)
	if v == nil {
		return nil, berror.NewNotFound(nil, fmt.Sprintf("Cannot find key[%s]", key))
	}
	f.static.Store(key, v)

	return defaultValue{v: v}, nil
}

func (f *yamlFile) DynamicLoad(ctx context.Context, key string, namespace ...string) (bconfig.Value, error) {
	// TODO implement me
	panic("implement me")
}

func (f *yamlFile) loadFile(dir, filename string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, berror.Convert(err, fmt.Sprintf("Failed to load config[%s:%s]", dir, filename), nil)
	}
	return v, nil
}
