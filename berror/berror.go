package berror

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  --proto_path=. grpc_status_detail.proto

// Error Provide the interface for feeding back business error info
type Error interface {
	Error() string                    // Output error information in string format
	Status() Status                   // Get main status
	WithSubStatus(Status) Error       // Carrying child status
	SubStatus() Status                // Get child status
	ListTracking(...int) []*TraceInfo // List the error tracking information that has been collected
	Unwrap() error                    // Support Is() and As() functions
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
