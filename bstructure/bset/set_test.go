package bset_test

import (
	"go-brick/bstructure/bset"
	"go-brick/bstructure/bslice"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	s := map[uint]struct{}{
		1:      {},
		2:      {},
		3:      {},
		999999: {},
	}
	assert.Equal(t,
		map[uint]struct{}{
			1: {}, 2: {}, 3: {}, 999999: {},
		},
		bset.Clone(s))
}

func TestFromSlice(t *testing.T) {
	var s1 = []uint{1, 2, 3, 4, 5, 5, 1}
	m1 := bset.FromSlice(s1)
	assert.Equal(t,
		map[uint]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
		}, m1)

	var s2 = []string{"1", "2", "3", "4", "5", "5", "1"}
	m2 := bset.FromSlice(s2)
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
	m3 := bset.FromSlice(s3)
	assert.Equal(t,
		map[uintptr]struct{}{
			uintptr(unsafe.Pointer(pa1)): {},
			uintptr(unsafe.Pointer(pa2)): {},
			uintptr(unsafe.Pointer(pa3)): {},
			uintptr(unsafe.Pointer(pa4)): {},
		}, m3)
}

func TestToSlice(t *testing.T) {
	var s = []uint{1, 2, 3, 4, 5, 5, 1}
	set := bset.FromSlice(s)
	assert.Equal(t, []uint{1, 2, 3, 4, 5}, bslice.SortNumbers(bset.ToSlice(set)))
}

func TestToSafeSet(t *testing.T) {
	var s = []uint{1, 2, 3, 4, 5, 5, 1}
	ss := bset.ToSafeSet(bset.FromSlice(s))
	assert.Equal(t, []uint{1, 2, 3, 4, 5}, bslice.SortNumbers(ss.ToSlice()))
}

func TestIntersectionSet(t *testing.T) {
	var (
		a = map[string]struct{}{
			"111": {},
			"222": {},
			"333": {},
			"444": {},
		}
		b = map[string]struct{}{
			"222": {},
			"333": {},
			"444": {},
			"555": {},
		}
		c = map[string]struct{}{
			"333": {},
			"444": {},
			"555": {},
			"666": {},
		}
	)
	assert.Equal(t,
		map[string]struct{}{
			"333": {},
			"444": {},
		},
		bset.IntersectionSet(a, b, c))
}

func TestUnionSet(t *testing.T) {
	var (
		a = map[string]struct{}{
			"111": {},
			"222": {},
			"333": {},
			"444": {},
		}
		b = map[string]struct{}{
			"222": {},
			"333": {},
			"444": {},
			"555": {},
		}
		c = map[string]struct{}{
			"333": {},
			"444": {},
			"555": {},
			"666": {},
		}
	)
	assert.Equal(t,
		map[string]struct{}{
			"111": {},
			"222": {},
			"333": {},
			"444": {},
			"555": {},
			"666": {},
		},
		bset.UnionSet(a, b, c))
}

func TestComplementSet(t *testing.T) {
	var (
		a = map[string]struct{}{
			"111": {},
			"222": {},
			"333": {},
			"444": {},
		}
		b = map[string]struct{}{
			"222": {},
			"333": {},
			"444": {},
			"555": {},
		}
		c = map[string]struct{}{
			"333": {},
			"444": {},
			"555": {},
			"666": {},
		}
	)
	assert.Equal(t,
		map[string]struct{}{
			"555": {},
		},
		bset.ComplementSet(a, b))
	assert.Equal(t,
		map[string]struct{}{
			"555": {},
			"666": {},
		},
		bset.ComplementSet(a, c))
}
