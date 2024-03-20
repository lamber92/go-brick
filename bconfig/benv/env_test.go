package benv_test

import (
	"go-brick/bconfig/benv"
	"testing"
)

func TestEnv(t *testing.T) {
	e, err := benv.Get()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e.GetType())
	t.Log(e.GetName())
	t.Log(e.AllowDebug())
	t.Log(e.Get("_test_", true))
}
