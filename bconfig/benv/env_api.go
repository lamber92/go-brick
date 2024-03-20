package benv

// Type environment type
type Type string

func (t Type) ToString() string {
	return string(t)
}

const (
	DEV Type = "dev" // Development environment
	FAT Type = "fat" // Factory Acceptance Test environment
	SIT Type = "sit" // System Integration Test environment
	UAT Type = "uat" // User Acceptance Test environment
	PRO Type = "pro" // Production environment
)

type Env interface {
	// GetType get environment type
	GetType() Type
	// GetName get environment name
	GetName() string
	// Get get environment value by key
	// if the environment variable does not exist or its value is empty, an error will be returned.
	Get(key string, fromCache ...bool) (string, error)
	// AllowDebug determine whether the current environment can be debugged
	AllowDebug() bool
}
