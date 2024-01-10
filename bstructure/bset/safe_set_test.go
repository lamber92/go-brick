package bset_test

import (
	"go-brick/bstructure/bset"
	"go-brick/bstructure/bslice"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_defaultSet = bset.NewSafeSet([]string{"1", "2", "3", "4", "5", "abc", "hhhhhhh"}...)
)

func TestSafeSet_Clone(t *testing.T) {
	set := _defaultSet.Clone()
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "abc", "hhhhhhh"}, bslice.SortStrings(set.ToSlice()))
}

func TestSafeSet_Add(t *testing.T) {
	set := _defaultSet.Clone()
	set.Add("1", "2", "3", "6", "7", "8")
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6", "7", "8", "abc", "hhhhhhh"}, bslice.SortStrings(set.ToSlice()))
}

func TestSafeSet_Delete(t *testing.T) {
	set := _defaultSet.Clone()
	set.Delete("1", "2", "3", "6", "7", "8", "hhh")
	assert.Equal(t, []string{"4", "5", "abc", "hhhhhhh"}, bslice.SortStrings(set.ToSlice()))
}

func TestSafeSet_Has(t *testing.T) {
	set := _defaultSet.Clone()
	assert.Equal(t, true, set.Has("1"))
	assert.Equal(t, false, set.Has("world"))
}

func TestSafeSet_Contains(t *testing.T) {
	set := _defaultSet.Clone()

	v1, ok1 := set.Contains("1", "2", "3", "abc")
	assert.Equal(t, true, ok1)
	assert.Equal(t, []string{"1", "2", "3", "abc"}, bslice.SortStrings(v1.ToSlice()))

	v2, ok2 := set.Contains("1", "2", "3", "abc", "hhh")
	assert.Equal(t, true, ok2)
	assert.Equal(t, []string{"1", "2", "3", "abc"}, bslice.SortStrings(v2.ToSlice()))

	v3, ok3 := set.Contains("xxx", "yyy")
	assert.Equal(t, false, ok3)
	assert.Equal(t, 0, v3.Len())
}

func TestSafeSet_Clear(t *testing.T) {
	set := _defaultSet.Clone()
	set.Clear()
	assert.Equal(t, 0, set.Len())
}

func TestSafeSet_Len(t *testing.T) {
	set := _defaultSet.Clone()
	assert.Equal(t, 7, set.Len())

	set.Add("1", "2", "3", "6", "7", "8")
	assert.Equal(t, 10, set.Len())

	set.Delete("1", "2", "3", "xxx")
	assert.Equal(t, 7, set.Len())
}

func TestSafeSet_IsEmpty(t *testing.T) {
	set := _defaultSet.Clone()
	assert.Equal(t, false, set.IsEmpty())
	set.Clear()
	assert.Equal(t, true, set.IsEmpty())
}

func TestSafeSet_IntersectionSet(t *testing.T) {
	var (
		a = bset.NewSafeSet("111", "222", "333", "444")
		b = bset.NewSafeSet("222", "333", "444", "555")
		c = bset.NewSafeSet("333", "444", "555", "666")
	)
	assert.Equal(t, []string{"333", "444"}, bslice.SortStrings(a.IntersectionSet(b, c).ToSlice()))
}

func TestSafeSet_UnionSet(t *testing.T) {
	var (
		a = bset.NewSafeSet("111", "222", "333", "444")
		b = bset.NewSafeSet("222", "333", "444", "555")
		c = bset.NewSafeSet("333", "444", "555", "666")
	)
	assert.Equal(t, []string{"111", "222", "333", "444", "555", "666"}, bslice.SortStrings(a.UnionSet(b, c).ToSlice()))
}

func TestSafeSet_ComplementSet(t *testing.T) {
	var (
		a = bset.NewSafeSet("111", "222", "333", "444")
		b = bset.NewSafeSet("222", "333", "444", "555")
		c = bset.NewSafeSet("333", "444", "555", "666")
	)
	assert.Equal(t, []string{"555"}, bslice.SortStrings(a.ComplementSet(b).ToSlice()))
	assert.Equal(t, []string{"555", "666"}, bslice.SortStrings(a.ComplementSet(c).ToSlice()))
}
