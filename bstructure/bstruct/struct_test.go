package bstruct_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-brick/bstructure/bstruct"
	"testing"
)

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
	s1, err := bstruct.GetFieldMap[uint, Temp](test1, "Field1")
	t.Logf("%+v %+v\n", s1, err)
	// map[10:[{Field1:10 Field2:10} {Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:2000} {Field1:20 Field2:20}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>

	s2, err := bstruct.GetFieldMap[string, Temp](test1, "Field2")
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
	s1, err := bstruct.GetFieldMap[uint, *Temp](test1, "Field1")
	t.Logf("%+v %+v\n", s1, err)
	// map[10:[{Field1:10 Field2:10} {Field1:10 Field2:100}] 20:[{Field1:20 Field2:20} {Field1:20 Field2:2000} {Field1:20 Field2:20}] 30:[{Field1:30 Field2:30}] 40:[{Field1:40 Field2:40}]] <nil>

	s2, err := bstruct.GetFieldMap[string, *Temp](test1, "Field2")
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
	s1, err := bstruct.GetFieldValues[Temp, uint](test1, "Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 30, 40, 20, 20}, s1)

	s2, err := bstruct.GetFieldValues[Temp, string](test1, "Field2")
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
	s1, err := bstruct.GetFieldValues[*Temp, uint](test1, "Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 30, 40, 20, 20}, s1)

	s2, err := bstruct.GetFieldValues[*Temp, string](test1, "Field2")
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
	s1, err := bstruct.GetFieldValuesEx[*Temp2, uint](test1, "StField1.Field1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []uint{10, 10, 20, 20, 40, 20}, s1)

	s2, err := bstruct.GetFieldValuesEx[*Temp2, string](test1, "StField2.Field2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"100", "1000", "200", "200", "400", "20000"}, s2)
}
