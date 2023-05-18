package berror

type HookFunc func(err error, reason string, detail any, options ...ConvOption) error

type Converter interface {
	// Convert according to the built-in rules, convert the incoming error to Error.
	Convert(err error, reason string, detail any, options ...ConvOption) error
	// Hook custom error conversion hook function.
	// nb. if you need this Hook, call it when you initialize the program.
	Hook(f HookFunc)
}

type ConvOption interface{}
