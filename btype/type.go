// Package btype
// internal generic type definition
// refer to the definition from https://pkg.go.dev/golang.org/x/exp/constraints
// because the `x` package is still in the iteration process and may be unstable,
// only some definitions are extracted here for use and modification.
package btype

// Ordered data types that can be compared and sorted
type Ordered interface {
	Number | String
}

type Number interface {
	Integer | Float
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type String interface {
	~string
}

type Struct interface {
	CanConvert() bool
}
