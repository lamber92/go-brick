package bstorage

import (
	"context"
	"time"
)

// Type built-in config type
type Type int

const (
	YAML   Type = 1
	APOLLO Type = 2
)

// Value interface.
// borrowed from spf13/viper interface design
type Value interface {
	// Sub returns new Viper instance representing a sub tree of this instance.
	// Sub is case-insensitive for a key.
	Sub(key string) Value

	// GetInt returns the value associated with the key as an integer.
	GetInt(key string) int
	// GetUint returns the value associated with the key as an unsigned integer.
	GetUint(key string) uint
	// GetString returns the value associated with the key as a string.
	GetString(key string) string
	// GetBool returns the value associated with the key as a boolean.
	GetBool(key string) bool
	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) time.Duration
	// GetIntSlice returns the value associated with the key as a slice of int values.
	GetIntSlice(key string) []int
	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) []string
	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) map[string]any

	// Unmarshal unmarshals the config into a Struct. Make sure that the tags
	// on the fields of the structure are properly set.
	Unmarshal(rawVal any) error
	// String format Value printer
	String() string
}

type OnChangeFunc func(event string)

type Config interface {
	// GetType get configuration type
	GetType() Type
	// Load load configuration Value
	Load(ctx context.Context, key string, namespace ...string) (Value, error)
	// RegisterOnChange register callback function for configuration changing notification
	// nb. this function is not thread-safe.
	RegisterOnChange(OnChangeFunc)
	// Close release resources
	Close()
}
