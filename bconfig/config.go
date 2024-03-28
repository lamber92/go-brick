package bconfig

import (
	"context"
	"fmt"

	"go-brick/bconfig/benv"
	"go-brick/bconfig/bstorage"
	"go-brick/bconfig/bstorage/apollo"
	"go-brick/bconfig/bstorage/yaml"
	"sync"
)

var (
	_once sync.Once

	Static  bstorage.Config
	Dynamic bstorage.Config
	Env     benv.Env
)

type Option struct {
	Type      bstorage.Type
	ConfigDir string
}

func Init(opt Option) {
	_once.Do(func() {
		var (
			err error
		)
		// load environment info
		Env, err = benv.Get()
		if err != nil {
			panic(err)
		}
		if len(opt.ConfigDir) > 0 {
			yaml.InitRootDir(opt.ConfigDir)
		}
		// init config manager from diff way by Type
		switch opt.Type {
		case bstorage.YAML:
			initFromYAML()
		case bstorage.APOLLO:
			if err = initFromApollo(); err != nil {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("Unsupported Config-Type [%d]", opt.Type))
		}
	})
}

func initFromYAML() {
	Static = yaml.NewStatic()
	Dynamic = yaml.NewDynamic()
}

func initFromApollo() error {
	basic, err := yaml.NewStatic().Load(context.Background(), "Apollo", Env.GetName())
	if err != nil {
		return err
	}

	conf := &apollo.Config{}
	if err = basic.Unmarshal(conf); err != nil {
		return err
	}
	Dynamic, err = apollo.New(conf)
	if err != nil {
		return err
	}
	Static = Dynamic

	return nil
}

func Close() {
	switch Static.GetType() {
	case bstorage.YAML:
		Static.Close()
		Dynamic.Close()
	case bstorage.APOLLO:
		Dynamic.Close()
	}
}
