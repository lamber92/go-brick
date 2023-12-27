package bstruct

import (
	"errors"
	"fmt"
	"go-brick/btype"
	"reflect"
	"strings"
)

// GetFieldMap specify a field of the structure as the key and convert the structure slice into a mapping table.
// the source slice can be a structure object slice or a structure pointer slice,
// other types will return err.
func GetFieldMap[Tk btype.Number | ~string, Tv btype.Struct](src []Tv, fieldName string) (map[Tk][]Tv, error) {
	if len(src) == 0 || len(fieldName) == 0 {
		return map[Tk][]Tv{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make(map[Tk][]Tv)
	for _, elem := range src {
		field, err := getStructFieldValue(reflect.ValueOf(elem), fieldName)
		if err != nil {
			return nil, err
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
	if len(src) == 0 || len(fieldName) == 0 {
		return []Tr{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}

	r := make([]Tr, 0, len(src))
	for _, elem := range src {
		field, err := getStructFieldValue(reflect.ValueOf(elem), fieldName)
		if err != nil {
			return nil, err
		}
		v, ok := field.Interface().(Tr)
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
func GetFieldValuesEx[T btype.Struct, Tr any](src []T, fieldName string) ([]Tr, error) {
	if len(src) == 0 || len(fieldName) == 0 {
		return []Tr{}, nil
	}
	if !src[0].CanConvert() {
		return nil, errors.New("source slice can not convert")
	}
	fieldNames := strings.Split(fieldName, ".")
	if len(fieldNames) == 1 {
		return GetFieldValues[T, Tr](src, fieldName)
	}

	r := make([]Tr, 0, len(src))
	for _, elem := range src {
		field, err := getNestedStructFieldValue(reflect.ValueOf(elem), fieldNames)
		if err != nil {
			return nil, err
		}
		v, ok := field.Interface().(Tr)
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
