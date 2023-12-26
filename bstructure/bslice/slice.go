package bslice

import (
	"errors"
	"fmt"
	"go-brick/btype"
	"reflect"
	"sort"
)

// Join concat two slices
func Join[T any](first []T, second []T) []T {
	r := make([]T, 0, len(first)+len(second))
	r = append(r, first...)
	r = append(r, second...)
	return r
}

// Joins concat multiple slices
func Joins[T any](src ...[]T) []T {
	length := 0
	for _, v := range src {
		length += len(v)
	}
	r := make([]T, 0, length)
	for _, v := range src {
		r = append(r, v...)
	}
	return r
}

// RemoveDuplicates Remove duplicate items
func RemoveDuplicates[T btype.Number | ~string](src []T) []T {
	r := make([]T, 0, len(src))
	check := make(map[T]struct{})
	for _, v := range src {
		_, exist := check[v]
		if !exist {
			check[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// SortNumbers Sort numeric slice
func SortNumbers[T btype.Number](src []T, desc ...bool) []T {
	if len(src) == 0 {
		return src
	}
	if desc != nil && desc[0] {
		sort.Slice(src, func(i, j int) bool {
			if src[i] > src[j] {
				return true
			}
			return false
		})
	} else {
		sort.Slice(src, func(i, j int) bool {
			if src[i] < src[j] {
				return true
			}
			return false
		})
	}
	return src
}

// descStringSlice attaches the methods of Interface to []string, sorting in decreasing order.
type descStringSlice []string

func (x descStringSlice) Len() int           { return len(x) }
func (x descStringSlice) Less(i, j int) bool { return x[i] > x[j] }
func (x descStringSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// SortStings Sort string slice
func SortStings(src []string, desc ...bool) []string {
	if len(src) == 0 {
		return src
	}
	if desc != nil && desc[0] {
		sort.Sort(descStringSlice(src))
	} else {
		sort.Sort(sort.StringSlice(src))
	}
	return src
}

// GetFieldMap specify a field of the structure as the key and convert the structure slice into a mapping table.
// the source slice can be a structure object slice or a structure pointer slice,
// other types will return err.
func GetFieldMap[Tk btype.Number | ~string, Tv btype.Struct](src []Tv, fieldName string) (map[Tk][]Tv, error) {
	if len(src) == 0 {
		return map[Tk][]Tv{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make(map[Tk][]Tv)
	for _, elem := range src {
		val := reflect.ValueOf(elem)

		var field reflect.Value
		switch val.Kind() {
		case reflect.Struct:
			field = val.FieldByName(fieldName)
		case reflect.Pointer:
			val = val.Elem()
			if val.Kind() != reflect.Struct {
				return nil, errors.New("source slice element is not a struct")
			}
			field = val.FieldByName(fieldName)
		default:
			return nil, errors.New("source slice element is not a struct")
		}
		if !field.IsValid() {
			return nil, errors.New(fmt.Sprintf("source slice element has not field name '%s'", fieldName))
		}
		switch key := field.Interface().(type) {
		case Tk:
			s, exist := r[key]
			if !exist {
				s = make([]Tv, 0)
			}
			r[key] = append(s, elem)
		default:
			return nil, errors.New(fmt.Sprintf("source slice element field should not be key for map '%+v'", key))
		}
	}
	return r, nil
}

// GetFieldValues extract the specified structure field value and return it in a slice.
// the source slice can be a structure object slice or a structure pointer slice,
// other types will return err.
func GetFieldValues[T btype.Struct, Tr any](src []T, fieldName string) ([]Tr, error) {
	if len(src) == 0 {
		return []Tr{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make([]Tr, 0, len(src))
	for _, elem := range src {
		val := reflect.ValueOf(elem)

		var field reflect.Value
		switch val.Kind() {
		case reflect.Struct:
			field = val.FieldByName(fieldName)
		case reflect.Pointer:
			val = val.Elem()
			if val.Kind() != reflect.Struct {
				return nil, errors.New("source slice element is not a struct")
			}
			field = val.FieldByName(fieldName)
		default:
			return nil, errors.New("source slice element is not a struct")
		}
		if !field.IsValid() {
			return nil, errors.New(fmt.Sprintf("source slice element has not field name '%s'", fieldName))
		}
		r = append(r, field.Interface().(Tr))
	}
	return r, nil
}
