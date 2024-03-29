package bstruct

import (
	"errors"
	"fmt"
	"github.com/lamber92/go-brick/btype"
	"reflect"
	"strings"
)

// GetFieldMap specify a field of the structure as the key and convert the structure slice into a mapping table.
// the source slice can be a structure object slice or a structure pointer slice,
// other types will return err.
func GetFieldMap[K comparable, V btype.Struct](src []V, fieldName string) (map[K][]V, error) {
	if len(src) == 0 || len(fieldName) == 0 {
		return map[K][]V{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make(map[K][]V)
	for _, elem := range src {
		field, err := getStructFieldValue(reflect.ValueOf(elem), fieldName)
		if err != nil {
			return nil, err
		}
		switch key := field.Interface().(type) {
		case K:
			s, exist := r[key]
			if !exist {
				s = make([]V, 0)
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
func GetFieldValues[T btype.Struct, V any](src []T, fieldName string) ([]V, error) {
	if len(src) == 0 || len(fieldName) == 0 {
		return []V{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make([]V, 0, len(src))
	for _, elem := range src {
		field, err := getStructFieldValue(reflect.ValueOf(elem), fieldName)
		if err != nil {
			return nil, err
		}
		v, ok := field.Interface().(V)
		if !ok {
			return nil, errors.New(fmt.Sprintf("the field type is not match. now: '%T'", field.Interface()))
		}
		r = append(r, v)
	}
	return r, nil
}

func getStructFieldValue(elem reflect.Value, fieldName string) (field reflect.Value, err error) {
	switch elem.Kind() {
	case reflect.Struct:
		field = elem.FieldByName(fieldName)
	case reflect.Pointer:
		elem = elem.Elem()
		if elem.Kind() != reflect.Struct {
			err = errors.New("source slice element is not a struct")
			return
		}
		field = elem.FieldByName(fieldName)
	default:
		err = errors.New("source slice element is not a struct")
		return
	}
	if !field.IsValid() {
		err = errors.New(fmt.Sprintf("source slice element has not field name '%s'", fieldName))
		return
	}
	return
}

// GetFieldValuesEx it is an enhanced version of GetFieldValue,
// that supports probing the fieldName value of nested structures.
// uses the symbol `.` to separate the field names of each layer of the structure.
func GetFieldValuesEx[T btype.Struct, V any](src []T, fieldName string) ([]V, error) {
	if len(src) == 0 || len(fieldName) == 0 {
		return []V{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}
	fieldNames := strings.Split(fieldName, ".")
	if len(fieldNames) == 1 {
		return GetFieldValues[T, V](src, fieldName)
	}

	r := make([]V, 0, len(src))
	for _, elem := range src {
		field, err := getNestedStructFieldValue(reflect.ValueOf(elem), fieldNames)
		if err != nil {
			return nil, err
		}
		v, ok := field.Interface().(V)
		if !ok {
			return nil, errors.New(fmt.Sprintf("the field type is not match. now: '%T'", field.Interface()))
		}
		r = append(r, v)
	}
	return r, nil
}

func getNestedStructFieldValue(elem reflect.Value, fieldNames []string) (field reflect.Value, err error) {
	if len(fieldNames) > 1 {
		elem, err = getStructFieldValue(elem, fieldNames[0])
		if err != nil {
			return
		}
		return getNestedStructFieldValue(elem, fieldNames[1:])
	} else {
		return getStructFieldValue(elem, fieldNames[0])
	}
}
