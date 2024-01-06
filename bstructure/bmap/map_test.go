package bmap_test

import (
	"go-brick/bstructure/bmap"
	"go-brick/bstructure/bslice"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	type Temp struct {
		Field1 uint
		Field2 string
	}
	var test1 = map[uint]Temp{
		1: {},
		2: {},
		4: {},
	}
	assert.Equal(t, bslice.SortNumbers([]uint{1, 2, 4}), bslice.SortNumbers(bmap.Keys(test1)))
}

func TestValues(t *testing.T) {
	type Temp struct {
		Field1 uint
		Field2 string
	}
	var test1 = map[uint]Temp{
		1: {Field1: 10, Field2: "10"},
		2: {Field1: 20, Field2: "20"},
		4: {Field1: 30, Field2: "30"},
	}
	t.Logf("%+v\n", bmap.Values(test1))
}
