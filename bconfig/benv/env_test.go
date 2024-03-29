package benv_test

import (
	"go-brick/bconfig/benv"
	"go-brick/berror"
	"go-brick/berror/bcode"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	if err := os.Setenv("GO_ENV_NAME", "dev_apollo_config"); err != nil {
		panic(err)
	}
	t.Run()
}

func TestEnv_GetSystemENV(t *testing.T) {
	e, err := benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, benv.DEV, e.GetType())
	assert.Equal(t, "dev_apollo_config", e.GetName())
	assert.Equal(t, true, e.AllowDebug())
}

func TestEnv_Get(t *testing.T) {
	if err := os.Setenv("_test_", "something"); err != nil {
		panic(err)
	}
	e, err := benv.Get()
	if err != nil {
		t.Fatal(err)
	}

	v, err := e.Get("_test_")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "something", v)

	v, err = e.Get("_test_2_")
	assert.Equal(t, true, berror.IsCode(err, bcode.NotFound))
}

func TestEnv_AllowDebug(t *testing.T) {
	// dev
	e, err := benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, e.AllowDebug())

	// fat
	if err = os.Setenv("GO_ENV_NAME", "fat_apollo_config"); err != nil {
		t.Fatal(err)
	}
	e, err = benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, e.AllowDebug())

	// sit
	if err = os.Setenv("GO_ENV_NAME", "sit_apollo_config"); err != nil {
		t.Fatal(err)
	}
	e, err = benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, e.AllowDebug())

	// uat
	if err = os.Setenv("GO_ENV_NAME", "uat_apollo_config"); err != nil {
		t.Fatal(err)
	}
	e, err = benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, e.AllowDebug())

	// pro
	if err = os.Setenv("GO_ENV_NAME", "pro"); err != nil {
		t.Fatal(err)
	}
	e, err = benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, e.AllowDebug())
}
