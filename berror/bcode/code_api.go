package bcode

import (
	"google.golang.org/grpc/codes"
)

type Code interface {
	ToInt() int
	ToString() string
	Is(target any) bool // compare the target code value with the current error code value
}

type CodeConverter interface {
	ToGRPCCode(code Code) codes.Code   // internal business error code --> grpc code
	FromGRPCCode(code codes.Code) Code // grpc code --> internal business error code
	ToHTTPStatusCode(code Code) int    // internal business error code --> http response body code
	FromHTTPStatusCode(code int) Code  // http response body code --> internal business error code

	RegisterMapToGRPCCode(code Code, grpcCode codes.Code)    // register the mapping relationship from custom error code to grpc error code
	RegisterMapFromGRPCCode(grpcCode codes.Code, code Code)  // register the mapping relationship from grpc error code to custom error code
	RegisterMapToHTTPStatusCode(code Code, statusCode int)   // register the mapping relationship form custom error code to http-status-code
	RegisterMapFromHTTPStatusCode(statusCode int, code Code) // register the mapping relationship form http-status-code to custom error code
}
