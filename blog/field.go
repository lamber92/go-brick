package blog

import (
	"go-brick/blog/logger"
	"time"
)

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) logger.Field {
	return logger.NewField().Binary(key, val)
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) logger.Field {
	return logger.NewField().ByteString(key, val)
}

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, val [][]byte) logger.Field {
	return logger.NewField().ByteStrings(key, val)
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) logger.Field {
	return logger.NewField().Bool(key, val)
}

// Int constructs a field with the given key and value.
func Int(key string, val int) logger.Field {
	return logger.NewField().Int(key, val)
}

// Ints constructs a field that carries a slice of integers.
func Ints(key string, nums []int) logger.Field {
	return logger.NewField().Ints(key, nums)
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) logger.Field {
	return logger.NewField().Uint(key, val)
}

// Uints constructs a field that carries a slice of unsigned integers.
func Uints(key string, nums []uint) logger.Field {
	return logger.NewField().Uints(key, nums)
}

// String constructs a field with the given key and value.
func String(key string, val string) logger.Field {
	return logger.NewField().String(key, val)
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) logger.Field {
	return logger.NewField().Strings(key, ss)
}

// Time constructs a logger.Field {} with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) logger.Field {
	return logger.NewField().Time(key, val)
}

// Times constructs a field that carries a slice of time.Times.
func Times(key string, ts []time.Time) logger.Field {
	return logger.NewField().Times(key, ts)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) logger.Field {
	return logger.NewField().Duration(key, val)
}

// Durations constructs a field that carries a slice of time.Durations.
func Durations(key string, ds []time.Duration) logger.Field {
	return logger.NewField().Durations(key, ds)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value any) logger.Field {
	return logger.NewField().Any(key, value)
}
