package bslice_test

import (
	"go-brick/bstructure/bslice"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	uint1 := []uint{1, 2, 3, 4, 5}
	uint2 := []uint{1, 3, 5, 7, 9}
	assert.Equal(t, []uint{1, 2, 3, 4, 5, 1, 3, 5, 7, 9}, bslice.Join(uint1, uint2))

	str1 := []string{"1111", "2222", "3333"}
	str2 := []string{"2a", "3b", "4c"}
	assert.Equal(t, []string{"1111", "2222", "3333", "2a", "3b", "4c"}, bslice.Join(str1, str2))

	type TestStruct struct {
		Field1 int
		Field2 string
	}
	st1 := []TestStruct{{Field1: 11111, Field2: "xx"}, {Field1: 22222, Field2: "yy"}}
	st2 := []TestStruct{{Field1: 33333, Field2: "zz"}}
	assert.Equal(t, []TestStruct{
		{Field1: 11111, Field2: "xx"},
		{Field1: 22222, Field2: "yy"},
		{Field1: 33333, Field2: "zz"},
	}, bslice.Join(st1, st2))
}

func TestJoins(t *testing.T) {
	uint1 := []uint{1, 2, 3, 4, 5}
	uint2 := []uint{1, 3, 5, 7, 9}
	uint3 := []uint{2, 4, 6, 8, 10}
	assert.Equal(t,
		[]uint{1, 2, 3, 4, 5, 1, 3, 5, 7, 9, 2, 4, 6, 8, 10},
		bslice.Joins(uint1, uint2, uint3))

	str1 := []string{"1111", "2222", "3333"}
	str2 := []string{"2a", "3b", "4c"}
	str3 := []string{"555d", "666e", "777f"}
	assert.Equal(t,
		[]string{"1111", "2222", "3333", "2a", "3b", "4c", "555d", "666e", "777f"},
		bslice.Joins(str1, str2, str3))

	type TestStruct struct {
		Field1 int
		Field2 string
	}
	st1 := []TestStruct{{Field1: 11111, Field2: "xx"}, {Field1: 22222, Field2: "yy"}}
	st2 := []TestStruct{{Field1: 33333, Field2: "zz"}}
	st3 := []TestStruct{{Field1: 44444, Field2: "hello"}, {Field1: 55555, Field2: "world"}}
	assert.Equal(t, []TestStruct{
		{Field1: 11111, Field2: "xx"},
		{Field1: 22222, Field2: "yy"},
		{Field1: 33333, Field2: "zz"},
		{Field1: 44444, Field2: "hello"},
		{Field1: 55555, Field2: "world"},
	}, bslice.Joins(st1, st2, st3))
}

func TestRemoveDuplicates(t *testing.T) {
	s1 := []string{"1", "xxx", "yyyyy", "yy", "xxx", "x", "x", "x", "yy"}
	assert.Equal(t,
		[]string{"1", "xxx", "yyyyy", "yy", "x"},
		bslice.RemoveDuplicates(s1))

	s2 := []int64{1, 2, 3, 4, 99, 99, 100, 101, 100, 99}
	assert.Equal(t,
		[]int64{1, 2, 3, 4, 99, 100, 101},
		bslice.RemoveDuplicates(s2))

	s3 := []float64{1, 2.00000000000001, 3, 4.123, 99, 99, 100, 101, 100, 2.00000000000001, 2.00000000000001}
	assert.Equal(t,
		[]float64{1, 2.00000000000001, 3, 4.123, 99, 100, 101},
		bslice.RemoveDuplicates(s3))
}

func TestSortNumber(t *testing.T) {
	s1 := []int64{99, 1111, 12312, 11, 2, 11, 2}
	assert.Equal(t, []int64{2, 2, 11, 11, 99, 1111, 12312}, bslice.SortNumber(s1))
	assert.Equal(t, []int64{12312, 1111, 99, 11, 11, 2, 2}, bslice.SortNumber(s1, true))
}

func TestSortSting(t *testing.T) {
	s1 := []string{"aaa", "bbb", "a", "b", "xxxx", "ccc", "abc"}
	assert.Equal(t, []string{"a", "aaa", "abc", "b", "bbb", "ccc", "xxxx"}, bslice.SortSting(s1))
	assert.Equal(t, []string{"xxxx", "ccc", "bbb", "b", "abc", "aaa", "a"}, bslice.SortSting(s1, true))
}
