package yamlfile

import (
	"context"
	"fmt"
	"go-brick/bconfig"
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/internal/bufferpool"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	defaultRoot = "config"

	defaultDirStatic  = "static"
	defaultDirDynamic = "dynamic"
	defaultFilename   = "config"
)

var (
	_Root = defaultRoot
)

// TODO: 加载自定义配置文件目录路径

var _config = newDefault()

func newDefault() *yamlFile {
	return &yamlFile{
		static:  sync.Map{},
		dynamic: sync.Map{},
	}
}

type yamlFile struct {
	static  sync.Map // map[filename]config
	dynamic sync.Map // map[filename]config
}

func (f *yamlFile) StaticLoad(key string, filenames ...string) (bconfig.Value, error) {
	var filename = defaultFilename
	if len(filenames) > 0 {
		filename = filenames[0]
	}
	// try to get from cache
	cache, ok := f.static.Load(filename)
	if ok {
		return f.getSub(cache.(*viper.Viper), key)
	}
	// read config file
	newData, err := f.loadFile(filename, true)
	if err != nil {
		return nil, err
	}
	if newData == nil {
		return nil, f.notfoundError(key)
	}
	// cache config
	// do not check key is existing or not
	f.static.Store(filename, newData)

	return f.getSub(newData, key)
}

func (f *yamlFile) DynamicLoad(ctx context.Context, key string, filenames ...string) (bconfig.Value, error) {
	// TODO: 需要将读取配置动作加入到上下文追踪信息

	var filename = defaultFilename
	if len(filenames) > 0 {
		filename = filenames[0]
	}

	// try to get from cache
	cache, ok := f.static.Load(filename)
	if ok {
		return f.getSub(cache.(*viper.Viper), key)
	}

	// TODO: 加原子锁。堵塞等待获取解锁。解锁后再次读缓存。
	// Load()到这里的时间间隙很短，忽略多次读同一个文件的情况。

	// read config file
	newData, err := f.loadFile(filename, true)
	if err != nil {
		return nil, err
	}
	if newData == nil {
		return nil, f.notfoundError(key)
	}

	// run watcher
	newData.WatchConfig()
	newData.OnConfigChange(func(in fsnotify.Event) {
		logger.Infra.Infow("[EVENT] config update", "event", in.String())
	})

	// cache config
	// do not check key is existing or not
	f.static.Store(filename, newData)

	return f.getSub(newData, key)
}

func (f *yamlFile) loadFile(filename string, static bool) (*viper.Viper, error) {
	dir := f.getDir(static)
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, berror.Convert(err, fmt.Sprintf("Failed to load config[%s:%s]", dir, filename), nil)
	}
	return v, nil
}

func (f *yamlFile) getDir(static bool) string {
	tmp := bufferpool.Get()
	tmp.AppendString(_Root)
	if static {
		tmp.AppendByte(os.PathSeparator)
		tmp.AppendString(defaultDirStatic)
	} else {
		tmp.AppendByte(os.PathSeparator)
		tmp.AppendString(defaultDirDynamic)
	}
	return tmp.String()
}

func (f *yamlFile) getSub(v *viper.Viper, key string) (bconfig.Value, error) {
	sub := v.Sub(key)
	if sub == nil {
		return nil, f.notfoundError(key)
	}
	return newDefaultValue(sub), nil
}

func (f *yamlFile) notfoundError(key string) error {
	return berror.NewNotFound(nil, fmt.Sprintf("Cannot find key[%s]", key))
}
