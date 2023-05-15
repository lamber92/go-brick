package berror

import (
	"go-brick/berror/bstatus"
	"go-brick/bstack"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  --proto_path=. grpc_status_detail.proto

// Error Provide the interface for feeding back business error info
type Error interface {
	// Error output error information in string format.
	Error() string
	// Status get main status.
	Status() bstatus.Status
	// Stack tracking list the error tracking information that has been collected.
	Stack() bstack.StackList
}

type Warp interface {
	// Wrap nest the specified error into error chain.
	Wrap(error) error
	// Unwrap returns the next error in the error chain.
	Unwrap() error
}
