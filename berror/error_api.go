package berror

import "go-brick/berror/bstatus"

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  --proto_path=. grpc_status_detail.proto

// Error Provide the interface for feeding back business error info
type Error interface {
	Error() string          // Error output error information in string format.
	Status() bstatus.Status // Status get main status.

	Tracking(depth ...int) []*TraceInfo // Tracking list the error tracking information that has been collected.

	Wrap(err error) error // Wrap nest the specified error into error chain.
	Unwrap() error        // Unwrap returns the next error in the error chain.
}
