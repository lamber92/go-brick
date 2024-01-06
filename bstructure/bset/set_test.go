package bset_test

import (
	"go-brick/bstructure/bset"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestGetFromSlice(t *testing.T) {
	var s1 = []uint{1, 2, 3, 4, 5, 5, 1}
	m1 := bset.GetFromSlice(s1)
	assert.Equal(t,
		map[uint]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
		}, m1)

	var s2 = []string{"1", "2", "3", "4", "5", "5", "1"}
	m2 := bset.GetFromSlice(s2)
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
	m3 := bset.GetFromSlice(s3)
	assert.Equal(t,
		map[uintptr]struct{}{
			uintptr(unsafe.Pointer(pa1)): {},
			uintptr(unsafe.Pointer(pa2)): {},
			uintptr(unsafe.Pointer(pa3)): {},
			uintptr(unsafe.Pointer(pa4)): {},
		}, m3)
}
