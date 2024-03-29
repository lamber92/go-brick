package benv

import (
	"fmt"
	"os"
	"strings"

	"github.com/lamber92/go-brick/berror"
	"github.com/spf13/viper"
)

const (
	// _ENV_VAR_KEY_ environment variable of environment name
	_ENV_VAR_KEY_ = "GO_ENV_NAME"
)

var (
	// allowDebugType define which environments can be debugged
	allowDebugType = map[Type]struct{}{
		DEV: {},
		FAT: {},
		SIT: {},
	}
)

type env struct {
	typ  Type
	name string // environment name
	kv   *viper.Viper
}

// Get return environment info
func Get() (Env, error) {
	name := os.Getenv(_ENV_VAR_KEY_)
	if len(name) == 0 {
		return nil, berror.NewNotFound(nil, fmt.Sprintf("cannot find environment variable [%s]", _ENV_VAR_KEY_))
	}

	res := &env{
		name: strings.ToLower(name),
		kv:   viper.New(),
	}
	switch res.name {
	case PRO.ToString():
		res.typ = PRO
	default:
		if strings.HasPrefix(res.name, DEV.ToString()) { // dev
			res.typ = DEV
		} else if strings.HasPrefix(res.name, FAT.ToString()) { // test
			res.typ = FAT
		} else if strings.HasPrefix(res.name, SIT.ToString()) { // pre
			res.typ = SIT
		} else if strings.HasPrefix(res.name, UAT.ToString()) { // pre
			res.typ = UAT
		} else {
			return nil, berror.NewInvalidArgument(nil, fmt.Sprintf("[%s] is invalid", _ENV_VAR_KEY_))
		}
	}
	return res, nil
}

func (e *env) GetType() Type {
	return e.typ
}

func (e *env) GetName() string {
	return e.name
}

func (e *env) Get(key string, fromCache ...bool) (string, error) {
	if len(fromCache) > 0 && fromCache[0] {
		if res := e.kv.GetString(key); len(res) > 0 {
			return res, nil
		}
	}
	if err := e.kv.BindEnv(key); err != nil {
		return "", berror.NewNotFound(err, fmt.Sprintf("cannot find environment variable [%s]", key))
	}
	res := e.kv.GetString(key)
	if len(res) == 0 {
		return "", berror.NewNotFound(nil, fmt.Sprintf("environment variable [%s] is empty", key))
	}
	return res, nil
}

func (e *env) AllowDebug() bool {
	_, ok := allowDebugType[e.typ]
	return ok
}
