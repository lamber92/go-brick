package bslice_test

import (
	"fmt"
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

func TestCombine(t *testing.T) {
	uint1 := []uint{1, 2, 3, 4, 5}
	uint2 := []uint{1, 3, 5, 7, 9}
	uint3 := []uint{2, 4, 6, 8, 10}
	uintSlice := [][]uint{uint1, uint2, uint3}
	assert.Equal(t,
		[]uint{1, 2, 3, 4, 5, 1, 3, 5, 7, 9, 2, 4, 6, 8, 10},
		bslice.Combine(uintSlice))

	str1 := []string{"1111", "2222", "3333"}
	str2 := []string{"2a", "3b", "4c"}
	str3 := []string{"555d", "666e", "777f"}
	strSlice := [][]string{str1, str2, str3}
	assert.Equal(t,
		[]string{"1111", "2222", "3333", "2a", "3b", "4c", "555d", "666e", "777f"},
		bslice.Combine(strSlice))

	type TestStruct struct {
		Field1 int
		Field2 string
	}
	st1 := []TestStruct{{Field1: 11111, Field2: "xx"}, {Field1: 22222, Field2: "yy"}}
	st2 := []TestStruct{{Field1: 33333, Field2: "zz"}}
	st3 := []TestStruct{{Field1: 44444, Field2: "hello"}, {Field1: 55555, Field2: "world"}}
	stSlice := [][]TestStruct{st1, st2, st3}
	assert.Equal(t, []TestStruct{
		{Field1: 11111, Field2: "xx"},
		{Field1: 22222, Field2: "yy"},
		{Field1: 33333, Field2: "zz"},
		{Field1: 44444, Field2: "hello"},
		{Field1: 55555, Field2: "world"},
	}, bslice.Combine(stSlice))
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
	assert.Equal(t, []int64{2, 2, 11, 11, 99, 1111, 12312}, bslice.SortNumbers(s1))
	assert.Equal(t, []int64{12312, 1111, 99, 11, 11, 2, 2}, bslice.SortNumbers(s1, true))
}

func TestSortSting(t *testing.T) {
	s1 := []string{"aaa", "bbb", "a", "b", "xxxx", "ccc", "abc"}
	assert.Equal(t, []string{"a", "aaa", "abc", "b", "bbb", "ccc", "xxxx"}, bslice.SortStings(s1))
	assert.Equal(t, []string{"xxxx", "ccc", "bbb", "b", "abc", "aaa", "a"}, bslice.SortStings(s1, true))
}

type Temp struct {
	Field1 uint
	Field2 string
}

func (Temp) CanConvert() bool { return true }

func (t Temp) String() string {
	return fmt.Sprintf("{Field1:%d Field2:%s}", t.Field1, t.Field2)
}

func TestGetFieldMap_StructObj(t *testing.T) {
	var test1 = []Temp{
		{Field1: 10, Field2: "10"},
		{Field1: 10, Field2: "100"},
		{Field1: 20, Field2: "20"},
		{Field1: 30, Field2: "30"},
		{Field1: 40, Field2: "40"},
		{Field1: 20, Field2: "2000"},
		{Field1: 20, Field2: "20"},
	}
	s1, err := bslice.GetFieldMap[uint, Temp](test1, "Field1")
	t.Logf("%+v %+v\n", s1, err)
	// map[10:[{Field1:10 Field2:10} {Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:2000} {Field1:20 Field2:20}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>

	s2, err := bslice.GetFieldMap[string, Temp](test1, "Field2")
	t.Logf("%+v %+v\n", s2, err)
	// map[10:[{Field1:10 Field2:10}] 100:[{Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:20}] 2000:[{Field1:20 Field2:2000}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>
}

func TestGetFieldMap_StructPtr(t *testing.T) {
	var test1 = []*Temp{
		{Field1: 10, Field2: "10"},
		{Field1: 10, Field2: "100"},
		{Field1: 20, Field2: "20"},
		{Field1: 30, Field2: "30"},
		{Field1: 40, Field2: "40"},
		{Field1: 20, Field2: "2000"},
		{Field1: 20, Field2: "20"},
	}
	s1, err := bslice.GetFieldMap[uint, *Temp](test1, "Field1")
	t.Logf("%+v %+v\n", s1, err)
	// map[10:[{Field1:10 Field2:10} {Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:2000} {Field1:20 Field2:20}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>

	s2, err := bslice.GetFieldMap[string, *Temp](test1, "Field2")
	t.Logf("%+v %+v\n", s2, err)
	// map[10:[{Field1:10 Field2:10}] 100:[{Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:20}] 2000:[{Field1:20 Field2:2000}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>
}

func TestGetFieldValues_StructObj(t *testing.T) {
	var test1 = []Temp{
		{Field1: 10, Field2: "10"},
		{Field1: 10, Field2: "100"},
		{Field1: 20, Field2: "20"},
		{Field1: 30, Field2: "30"},
		{Field1: 40, Field2: "40"},
		{Field1: 20, Field2: "2000"},
		{Field1: 20, Field2: "20"},
	}
	s1, err := bslice.GetFieldValues[Temp, uint](test1, "Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 30, 40, 20, 20}, s1)

	s2, err := bslice.GetFieldValues[Temp, string](test1, "Field2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"10", "100", "20", "30", "40", "2000", "20"}, s2)
}

func TestGetFieldValues_StructPtr(t *testing.T) {
	var test1 = []*Temp{
		{Field1: 10, Field2: "10"},
		{Field1: 10, Field2: "100"},
		{Field1: 20, Field2: "20"},
		{Field1: 30, Field2: "30"},
		{Field1: 40, Field2: "40"},
		{Field1: 20, Field2: "2000"},
		{Field1: 20, Field2: "20"},
	}
	s1, err := bslice.GetFieldValues[*Temp, uint](test1, "Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 30, 40, 20, 20}, s1)

	s2, err := bslice.GetFieldValues[*Temp, string](test1, "Field2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"10", "100", "20", "30", "40", "2000", "20"}, s2)
}

type Temp2 struct {
	StField1 Temp
	StField2 *Temp
}

func (Temp2) CanConvert() bool { return true }

func TestGetFieldValuesEx(t *testing.T) {
	var test1 = []*Temp2{
		{StField1: Temp{Field1: 10, Field2: "10"}, StField2: &Temp{Field1: 100, Field2: "100"}},
		{StField1: Temp{Field1: 10, Field2: "100"}, StField2: &Temp{Field1: 200, Field2: "1000"}},
		{StField1: Temp{Field1: 20, Field2: "20"}, StField2: &Temp{Field1: 200, Field2: "200"}},
		{StField1: Temp{Field1: 20, Field2: "20"}, StField2: &Temp{Field1: 200, Field2: "200"}},
		{StField1: Temp{Field1: 40, Field2: "40"}, StField2: &Temp{Field1: 400, Field2: "400"}},
		{StField1: Temp{Field1: 20, Field2: "2000"}, StField2: &Temp{Field1: 200, Field2: "20000"}},
	}
	s1, err := bslice.GetFieldValuesEx[*Temp2, uint](test1, "StField1.Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 20, 40, 20}, s1)

	s2, err := bslice.GetFieldValuesEx[*Temp2, string](test1, "StField2.Field2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"100", "1000", "200", "200", "400", "20000"}, s2)
}
