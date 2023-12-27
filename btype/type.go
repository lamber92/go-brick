package btype

type Number interface {
	IntegerEx | Float
}

type IntegerEx interface {
	Integer | uintptr
}

type Integer interface {
	SignedInteger | UnsignedInteger
}

type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

type Struct interface {
	CanConvert() bool
}
