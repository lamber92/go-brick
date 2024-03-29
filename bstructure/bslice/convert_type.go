package bslice

import (
	"github.com/lamber92/go-brick/btype"
)

// ToUint64 Unsigned integer slice cast
// usually used to dock proto messages
func ToUint64[T btype.Unsigned](src []T) []uint64 {
	r := make([]uint64, 0, len(src))
	for _, v := range src {
		r = append(r, uint64(v))
	}
	return r
}

// ToInt64 signed integer slice cast
// usually used to dock proto messages
func ToInt64[T btype.Signed](src []T) []int64 {
	r := make([]int64, 0, len(src))
	for _, v := range src {
		r = append(r, int64(v))
	}
	return r
}

// TODO: 不同数据类型的强制类型转换
