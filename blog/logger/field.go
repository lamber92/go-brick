package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// =======================================
// ------ Default Field Field IMPL -------
// =======================================

type defaultFields []Field

func (df defaultFields) Release() []zapcore.Field {
	out := make([]zap.Field, 0, len(df))
	for _, v := range df {
		field := v.(*defaultField).Field
		if field.Type == zapcore.UnknownType {
			continue
		}
		out = append(out, v.(*defaultField).Field)
	}
	return out
}

func NewField() Field {
	return &defaultField{}
}

type defaultField struct {
	zapcore.Field
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func (f *defaultField) Binary(key string, val []byte) Field {
	f.Field = zap.Binary(key, val)
	return f
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func (f *defaultField) ByteString(key string, val []byte) Field {
	f.Field = zap.ByteString(key, val)
	return f
}

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func (f *defaultField) ByteStrings(key string, val [][]byte) Field {
	f.Field = zap.ByteStrings(key, val)
	return f
}

// Bool constructs a field that carries a bool.
func (f *defaultField) Bool(key string, val bool) Field {
	f.Field = zap.Bool(key, val)
	return f
}

// Int constructs a field with the given key and value.
func (f *defaultField) Int(key string, val int) Field {
	f.Field = zap.Int64(key, int64(val))
	return f
}

// Ints constructs a field that carries a slice of integers.
func (f *defaultField) Ints(key string, nums []int) Field {
	f.Field = zap.Ints(key, nums)
	return f
}

// Uint constructs a field with the given key and value.
func (f *defaultField) Uint(key string, val uint) Field {
	f.Field = zap.Uint64(key, uint64(val))
	return f
}

// Uints constructs a field that carries a slice of unsigned integers.
func (f *defaultField) Uints(key string, nums []uint) Field {
	f.Field = zap.Uints(key, nums)
	return f
}

// String constructs a field with the given key and value.
func (f *defaultField) String(key string, val string) Field {
	f.Field = zap.String(key, val)
	return f
}

// Strings constructs a field that carries a slice of strings.
func (f *defaultField) Strings(key string, ss []string) Field {
	f.Field = zap.Strings(key, ss)
	return f
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func (f *defaultField) Time(key string, val time.Time) Field {
	f.Field = zap.Time(key, val)
	return f
}

// Times constructs a field that carries a slice of time.Times.
func (f *defaultField) Times(key string, ts []time.Time) Field {
	f.Field = zap.Times(key, ts)
	return f
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func (f *defaultField) Duration(key string, val time.Duration) Field {
	f.Field = zap.Duration(key, val)
	return f
}

// Durations constructs a field that carries a slice of time.Durations.
func (f *defaultField) Durations(key string, ds []time.Duration) Field {
	f.Field = zap.Durations(key, ds)
	return f
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func (f *defaultField) Any(key string, value any) Field {
	var field zap.Field

	switch val := value.(type) {
	case zapcore.ObjectMarshaler:
		field = zap.Object(key, val)
	case zapcore.ArrayMarshaler:
		field = zap.Array(key, val)
	case bool:
		field = zap.Bool(key, val)
	case *bool:
		field = zap.Boolp(key, val)
	case []bool:
		field = zap.Bools(key, val)
	case complex128:
		field = zap.Complex128(key, val)
	case *complex128:
		field = zap.Complex128p(key, val)
	case []complex128:
		field = zap.Complex128s(key, val)
	case complex64:
		field = zap.Complex64(key, val)
	case *complex64:
		field = zap.Complex64p(key, val)
	case []complex64:
		field = zap.Complex64s(key, val)
	case float64:
		field = zap.Float64(key, val)
	case *float64:
		field = zap.Float64p(key, val)
	case []float64:
		field = zap.Float64s(key, val)
	case float32:
		field = zap.Float32(key, val)
	case *float32:
		field = zap.Float32p(key, val)
	case []float32:
		field = zap.Float32s(key, val)
	case int:
		field = zap.Int(key, val)
	case *int:
		field = zap.Intp(key, val)
	case []int:
		field = zap.Ints(key, val)
	case int64:
		field = zap.Int64(key, val)
	case *int64:
		field = zap.Int64p(key, val)
	case []int64:
		field = zap.Int64s(key, val)
	case int32:
		field = zap.Int32(key, val)
	case *int32:
		field = zap.Int32p(key, val)
	case []int32:
		field = zap.Int32s(key, val)
	case int16:
		field = zap.Int16(key, val)
	case *int16:
		field = zap.Int16p(key, val)
	case []int16:
		field = zap.Int16s(key, val)
	case int8:
		field = zap.Int8(key, val)
	case *int8:
		field = zap.Int8p(key, val)
	case []int8:
		field = zap.Int8s(key, val)
	case string:
		field = zap.String(key, val)
	case *string:
		field = zap.Stringp(key, val)
	case []string:
		field = zap.Strings(key, val)
	case uint:
		field = zap.Uint(key, val)
	case *uint:
		field = zap.Uintp(key, val)
	case []uint:
		field = zap.Uints(key, val)
	case uint64:
		field = zap.Uint64(key, val)
	case *uint64:
		field = zap.Uint64p(key, val)
	case []uint64:
		field = zap.Uint64s(key, val)
	case uint32:
		field = zap.Uint32(key, val)
	case *uint32:
		field = zap.Uint32p(key, val)
	case []uint32:
		field = zap.Uint32s(key, val)
	case uint16:
		field = zap.Uint16(key, val)
	case *uint16:
		field = zap.Uint16p(key, val)
	case []uint16:
		field = zap.Uint16s(key, val)
	case uint8:
		field = zap.Uint8(key, val)
	case *uint8:
		field = zap.Uint8p(key, val)
	case []byte:
		field = zap.Binary(key, val)
	case uintptr:
		field = zap.Uintptr(key, val)
	case *uintptr:
		field = zap.Uintptrp(key, val)
	case []uintptr:
		field = zap.Uintptrs(key, val)
	case time.Time:
		field = zap.Time(key, val)
	case *time.Time:
		field = zap.Timep(key, val)
	case []time.Time:
		field = zap.Times(key, val)
	case time.Duration:
		field = zap.Duration(key, val)
	case *time.Duration:
		field = zap.Durationp(key, val)
	case []time.Duration:
		field = zap.Durations(key, val)
	case error:
		field = zap.NamedError(key, val)
	case []error:
		field = zap.Errors(key, val)
	case fmt.Stringer:
		field = zap.Stringer(key, val)
	default:
		field = zap.Reflect(key, val)
	}

	f.Field = field
	return f
}
