package berror

import (
	"go-brick/berror/bstatus"
	"go-brick/bstack"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  --proto_path=. grpc_status_detail.proto

// Error Provide the interface for feeding back business error info
type Error interface {
	Chain
	// Error output error information in string format.
	Error() string
	// Status get main status.
	Status() bstatus.Status
	// Stack tracking list the error tracking information that has been collected.
	Stack() bstack.StackList
}

type Chain interface {
	// Cause returns the underlying cause of the error, if possible.
	Cause() error
	// Unwrap provides compatibility for Go 1.13 error chains.
	Unwrap() error
}
