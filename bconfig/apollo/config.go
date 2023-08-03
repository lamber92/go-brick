package apollo

import (
	"context"
	"errors"
	"fmt"
	"go-brick/bconfig"
	"go-brick/berror"
	"go-brick/btrace"
	"sync"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/perror"
)

const (
	defaultApplication = "application"
)

type Config struct {
	Host        string
	AppID       string
	Cluster     string
	Namespace   string
	IsBackup    bool
	Secret      string
	Label       string
	SyncTimeout int
}

type apolloConfig struct {
	client agollo.Client
	sync.Mutex
}

func New(conf *Config, logger ...log.LoggerInterface) (bconfig.Config, error) {
	lgr := newDefaultLogger()
	if len(logger) > 0 {
		lgr = logger[0]
	}
	log.InitLogger(lgr)
	return newConfig(conf)
}

func newConfig(conf *Config) (*apolloConfig, error) {
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		appConfig := &config.AppConfig{
			AppID:             conf.AppID,
			Cluster:           conf.Cluster,
			NamespaceName:     conf.Namespace,
			IP:                conf.Host,
			IsBackupConfig:    conf.IsBackup,
			Secret:            conf.Secret,
			Label:             conf.Label,
			SyncServerTimeout: conf.SyncTimeout,
		}
		return appConfig, nil
	})
	if err != nil {
		return nil, berror.Convert(err, "init apollo-client failed")
	}

	return &apolloConfig{
		client: client,
	}, nil
}

func (a *apolloConfig) Load(ctx context.Context, key string, namespace ...string) (out bconfig.Value, err error) {
	ns := defaultApplication
	if len(namespace) > 0 {
		ns = namespace[0]
	}
	conf := a.client.GetConfig(ns)
	if conf == nil {
		err = berror.NewNotFound(nil, fmt.Sprintf("cannot find key in Apollo. namespace: %s | key: %s", ns, key))
		return
	}
	v, err := conf.GetCache().Get(key)
	if err != nil {
		if errors.Is(err, perror.ErrNotFound) {
			err = berror.NewNotFound(err, fmt.Sprintf("cannot find key in Apollo. namespace: %s | key: %s", ns, key))
			return
		}
		return
	}
	out = newDefaultValue(v)
	defer func() {
		if err == nil {
			btrace.AppendMDIntoCtx(ctx, newMetadata(ns, key, out))
		}
	}()
	return
}

func (a *apolloConfig) RegisterOnChange(changeFunc bconfig.OnChangeFunc) {
	hook := newDefaultListener(changeFunc)
	a.client.AddChangeListener(hook)
}

func (a *apolloConfig) Close() {
	a.Lock()
	if a.client != nil {
		a.client.Close()
	}
	a.client = nil
	a.Unlock()
}
