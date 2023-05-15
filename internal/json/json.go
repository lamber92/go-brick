package json

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// jsonStdIter output standard json iterators according to the official package
var jsonStdIter = jsoniter.ConfigCompatibleWithStandardLibrary

func MarshalToString(src interface{}) (string, error) {
	s, err := jsonStdIter.MarshalToString(src)
	if err != nil {
		return "", fmt.Errorf("json marshal failed. err: %s, source: %+v", err.Error(), src)
	}
	return s, nil
}

func Marshal(src interface{}) ([]byte, error) {
	b, err := jsonStdIter.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed. err: %s, source: %+v", err.Error(), src)
	}
	return b, nil
}

// MarshalIndent
// @n number of spaces for indentation
func MarshalIndent(src interface{}, n uint) ([]byte, error) {
	var (
		indent      = strings.Builder{}
		i      uint = 0
	)
	for i = 0; i < n; i++ {
		indent.WriteByte(' ')
	}
	b, err := jsonStdIter.MarshalIndent(src, "", indent.String())
	if err != nil {
		return nil, fmt.Errorf("json marshal failed. err: %s, source: %+v", err.Error(), src)
	}
	return b, nil
}

func UnmarshalFromString(src string, dest interface{}) error {
	if err := jsonStdIter.UnmarshalFromString(src, dest); err != nil {
		return fmt.Errorf("json unmarshal failed. err: %s, source: %+v", err.Error(), src)
	}
	return nil
}

func Unmarshal(src []byte, dest interface{}) error {
	if err := jsonStdIter.Unmarshal(src, dest); err != nil {
		return fmt.Errorf("json unmarshal failed. err: %s, source: %+v", err.Error(), src)
	}
	return nil
}

// Get get a value
func Get(data []byte, path ...interface{}) jsoniter.Any {
	return jsonStdIter.Get(data, path...)
}

func NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return jsonStdIter.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return jsonStdIter.NewDecoder(reader)
}

// Valid 判断入参是否是合法的json结构
func Valid(data []byte) bool {
	return jsonStdIter.Valid(data)
}
