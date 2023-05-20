package logger

import (
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

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func (f *defaultField) Any(key string, value any) Field {
	f.Field = zap.Any(key, value)
	return f
}
