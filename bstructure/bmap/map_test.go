package bmap_test

import (
	"go-brick/bstructure/bmap"
	"go-brick/bstructure/bslice"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	var s1 = []uint{1, 2, 3, 4, 5, 5, 1}
	m1 := bmap.FromSlice(s1)
	assert.Equal(t,
		map[uint]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
		}, m1)

	var s2 = []string{"1", "2", "3", "4", "5", "5", "1"}
	m2 := bmap.FromSlice(s2)
	assert.Equal(t,
		map[string]struct{}{
			"1": {},
			"2": {},
			"3": {},
			"4": {},
			"5": {},
		}, m2)

	var (
		a1, a2, a3, a4                           = 1, 2, 3, 4
		pa1, pa2, pa3, pa4, pa5, pa6             = &a1, &a2, &a3, &a4, &a2, &a4
		uptr1, uptr2, uptr3, uptr4, uptr5, uptr6 = uintptr(unsafe.Pointer(pa1)), uintptr(unsafe.Pointer(pa2)), uintptr(unsafe.Pointer(pa3)), uintptr(unsafe.Pointer(pa4)), uintptr(unsafe.Pointer(pa5)), uintptr(unsafe.Pointer(pa6))
		s3                                       = []uintptr{uptr1, uptr2, uptr3, uptr4, uptr5, uptr6}
	)
	m3 := bmap.FromSlice(s3)
	assert.Equal(t,
		map[uintptr]struct{}{
			uintptr(unsafe.Pointer(pa1)): {},
			uintptr(unsafe.Pointer(pa2)): {},
			uintptr(unsafe.Pointer(pa3)): {},
			uintptr(unsafe.Pointer(pa4)): {},
		}, m3)
}

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
	assert.Equal(t, bslice.SortNumber([]uint{1, 2, 4}), bslice.SortNumber(bmap.Keys(test1)))
}
