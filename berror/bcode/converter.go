package bcode

import (
	"google.golang.org/grpc/codes"
)

type CodeConverter interface {
	ToGRPCCode(code int) codes.Code   // Internal business error code --> grpc code
	FromGRPCCode(code codes.Code) int // grpc code --> Internal business error code
	ToHTTPRspBodyCode(code int) int   // Internal business error code --> http response body code
	FromHTTPRspBodyCode(code int) int // http response body code --> Internal business error code
}

var codeConverter CodeConverter = newDefaultCodeConverter()

// =======================================
// -------- Functional Interface ---------
// =======================================

// ReplaceCodeConverter Replace the default converter with a custom one
func ReplaceCodeConverter(conv CodeConverter) {
	codeConverter = conv
}

// ToGRPCCode Internal business error code --> grpc code
func ToGRPCCode(code int) codes.Code {
	return codeConverter.ToGRPCCode(code)
}

// FromGRPCCode grpc code --> Internal business error code
func FromGRPCCode(code codes.Code) int {
	return codeConverter.FromGRPCCode(code)
}

// ToHTTPRspBodyCode Internal business error code --> http resp body code
func ToHTTPRspBodyCode(code int) int {
	return codeConverter.ToHTTPRspBodyCode(code)
}

// FromHTTPRspBodyCode http resp body code --> Internal business error code
func FromHTTPRspBodyCode(code int) int {
	return codeConverter.FromHTTPRspBodyCode(code)
}

// =======================================
// ----- built-in default converter  -----
// =======================================

// defaultCodeConverter
type defaultCodeConverter struct{}

func newDefaultCodeConverter() *defaultCodeConverter {
	c := &defaultCodeConverter{}
	c.initMapBizCodeFromGRPCCode()
	return c
}

var (
	// internalCodeToGRPCCode Mapping relationship between internal business error codes and gRPC error codes
	internalCodeToGRPCCode = map[int]codes.Code{
		Unknown:            codes.Unknown,
		OK:                 codes.OK,
		InvalidArgument:    codes.InvalidArgument,
		Unauthorized:       codes.Unauthenticated,
		Forbidden:          codes.Aborted,
		NotFound:           codes.NotFound,
		RequestTimeout:     codes.DeadlineExceeded,
		ClientClosed:       codes.Canceled,
		InternalError:      codes.Internal,
		ServiceUnavailable: codes.Unavailable,
		GatewayTimeout:     codes.DeadlineExceeded,
		AlreadyExists:      codes.AlreadyExists,
	}

	// gRPCCodeToInternalCode Mapping relationship between gRPC error codes and internal business error codes
	gRPCCodeToInternalCode = map[codes.Code]int{}
)

// RegisterMapToGRPCCode
// Register the mapping relationship between custom error codes and grpc error codes
func RegisterMapToGRPCCode(code int, grpcCode codes.Code) {
	internalCodeToGRPCCode[code] = grpcCode
}

// initMapBizCodeFromGRPCCode Recalculate the mapping from gRPC error codes to internal error codes
func (*defaultCodeConverter) initMapBizCodeFromGRPCCode() {
	for k, v := range internalCodeToGRPCCode {
		if _, exist := gRPCCodeToInternalCode[v]; !exist {
			gRPCCodeToInternalCode[v] = k
		}
	}
}

func (*defaultCodeConverter) ToGRPCCode(code int) codes.Code {
	if c, ok := internalCodeToGRPCCode[code]; ok {
		return c
	}
	return codes.Unknown
}

func (*defaultCodeConverter) FromGRPCCode(code codes.Code) int {
	if c, ok := gRPCCodeToInternalCode[code]; ok {
		return c
	}
	return Unknown
}

var (
	// internalCodeToHTTPRspBodyCode Mapping relationship between internal business error codes and http body error codes
	internalCodeToHTTPRspBodyCode = map[int]int{}

	// httpRspBodyCodeToInternalCode Mapping relationship between http body error codes and internal business error codes
	httpRspBodyCodeToInternalCode = map[int]int{}
)

func (*defaultCodeConverter) ToHTTPRspBodyCode(code int) int {
	if c, ok := internalCodeToHTTPRspBodyCode[code]; ok {
		return c
	}
	return code
}

func (*defaultCodeConverter) FromHTTPRspBodyCode(code int) int {
	if c, ok := httpRspBodyCodeToInternalCode[code]; ok {
		return c
	}
	return code
}
