package bconfig

import (
	"context"
	"time"
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
}

type Config interface {
	// StaticLoad load configuration once.
	// throughout the lifetime, the configuration is read only once, and the value is cached.
	// calling again will fetch the data in the cache
	StaticLoad(key string, namespace ...string) (Value, error)
	// DynamicLoad load real-time configuration values, but allow for slight delays.
	DynamicLoad(ctx context.Context, key string, namespace ...string) (Value, error)
}