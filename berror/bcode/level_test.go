package bcode_test

import (
	"testing"

	"github.com/lamber92/go-brick/berror/bcode"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLevel(t *testing.T) {
	assert.Equal(t, bcode.LvInfo, bcode.OK.GetLevel())
	assert.Equal(t, bcode.LvNotice, bcode.NotFound.GetLevel())
	assert.Equal(t, bcode.LvWarning, bcode.RequestTimeout.GetLevel())
	assert.Equal(t, bcode.LvCritical, bcode.Unknown.GetLevel())
}

func TestCustomizedLevel(t *testing.T) {
	const (
		Lv1 bcode.Level = iota + 1
		Lv2
		Lv3
		Lv4
	)
	var myCodeToLevel = map[bcode.Code]bcode.Level{
		bcode.Unknown:         Lv3,
		bcode.OK:              Lv1,
		bcode.InvalidArgument: Lv2,
		bcode.Unauthorized:    Lv2,
		bcode.New(99999):      Lv1,
	}
	bcode.ReplaceCodeLevelMapping(myCodeToLevel, Lv4)

	assert.Equal(t, Lv3, bcode.Unknown.GetLevel())
	assert.Equal(t, Lv1, bcode.OK.GetLevel())
	assert.Equal(t, Lv2, bcode.InvalidArgument.GetLevel())
	assert.Equal(t, Lv2, bcode.Unauthorized.GetLevel())
	assert.Equal(t, Lv1, bcode.GetLevel(bcode.New(99999)))
	//
	assert.Equal(t, Lv4, bcode.RequestTimeout.GetLevel())
}
