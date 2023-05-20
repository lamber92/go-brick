package logger

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
	// Any takes a key and an arbitrary value and chooses the best way to represent
	// them as a field, falling back to a reflection-based approach only if
	// necessary.
	//
	// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
	// them. To minimize surprises, []byte values are treated as binary blobs, byte
	// values are treated as uint8, and runes are always treated as integers.
	Any(key string, value any) Field
}
