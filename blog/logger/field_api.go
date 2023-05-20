package logger

import "time"

// design borrowed from https://github.com/uber-go/zap

type Field interface {
	// Binary constructs a field that carries an opaque binary blob.
	//
	// Binary data is serialized in an encoding-appropriate format. For example,
	// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
	// use ByteString.
	Binary(key string, val []byte) Field

	// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
	// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
	// Binary.
	ByteString(key string, val []byte) Field

	// ByteStrings constructs a field that carries a slice of []byte, each of which
	// must be UTF-8 encoded text.
	ByteStrings(key string, val [][]byte) Field

	// Bool constructs a field that carries a bool.
	Bool(key string, val bool) Field

	// Int constructs a field with the given key and value.
	Int(key string, val int) Field

	// Ints constructs a field that carries a slice of integers.
	Ints(key string, nums []int) Field

	// Uint constructs a field with the given key and value.
	Uint(key string, val uint) Field

	// Uints constructs a field that carries a slice of unsigned integers.
	Uints(key string, nums []uint) Field

	// String constructs a field with the given key and value.
	String(key string, val string) Field

	// Strings constructs a field that carries a slice of strings.
	Strings(key string, ss []string) Field

	// Time constructs a Field with the given key and value. The encoder
	// controls how the time is serialized.
	Time(key string, val time.Time) Field

	// Times constructs a field that carries a slice of time.Times.
	Times(key string, ts []time.Time) Field

	// Duration constructs a field with the given key and value. The encoder
	// controls how the duration is serialized.
	Duration(key string, val time.Duration) Field

	// Durations constructs a field that carries a slice of time.Durations.
	Durations(key string, ds []time.Duration) Field

	// Any takes a key and an arbitrary value and chooses the best way to represent
	// them as a field, falling back to a reflection-based approach only if
	// necessary.
	//
	// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
	// them. To minimize surprises, []byte values are treated as binary blobs, byte
	// values are treated as uint8, and runes are always treated as integers.
	Any(key string, value any) Field
}
