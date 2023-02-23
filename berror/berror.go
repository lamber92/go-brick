package berror

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  --proto_path=. grpc_status_detail.proto

// Error Provide the interface for feeding back business error info
type Error interface {
	Error() string  // Error output error information in string format.
	Status() Status // Status get main status.

	Tracking(depth ...int) []*TraceInfo // Tracking list the error tracking information that has been collected.

	Wrap(Error)    // Wrap nest the specified error into error chain.
	Unwrap() error // Unwrap returns the next error in the error chain.
}

// Status Carrier of business error info
type Status interface {
	Code() int           // error code
	Reason() string      // error description
	Detail() interface{} // error extension
}

// TraceInfo Basic unit of position information when an error occurs
type TraceInfo struct {
	Func string // function name
	File string // file name
	Line int    // line
}
