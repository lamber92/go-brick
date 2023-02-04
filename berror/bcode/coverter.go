package bcode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var codeConverter CodeConverter = newDefaultCodeConverter()

// ReplaceCodeConverter replace the default converter with a custom one
func ReplaceCodeConverter(conv CodeConverter) {
	codeConverter = conv
}

// ToGRPCCode internal business error code --> grpc code
func ToGRPCCode(code Code) codes.Code {
	return codeConverter.ToGRPCCode(code)
}

// FromGRPCCode grpc code --> Internal business error code
func FromGRPCCode(code codes.Code) Code {
	return codeConverter.FromGRPCCode(code)
}

// ToHTTPStatusCode internal error code --> http-status-code
func ToHTTPStatusCode(code Code) int {
	return codeConverter.ToHTTPStatusCode(code)
}

// FromHTTPStatusCode http-status-code --> internal error code
func FromHTTPStatusCode(code int) Code {
	return codeConverter.FromHTTPStatusCode(code)
}

// RegisterMapToGRPCCode register the mapping relationship from custom error code to grpc error code
func RegisterMapToGRPCCode(code Code, grpcCode codes.Code) {
	codeConverter.RegisterMapToGRPCCode(code, grpcCode)
}

// RegisterMapFromGRPCCode register the mapping relationship from grpc error code to custom error code
func RegisterMapFromGRPCCode(grpcCode codes.Code, code Code) {
	codeConverter.RegisterMapFromGRPCCode(grpcCode, code)
}

// RegisterMapToHTTPStatusCode register the mapping relationship form custom error code to http-status-code
func RegisterMapToHTTPStatusCode(code Code, statusCode int) {
	codeConverter.RegisterMapToHTTPStatusCode(code, statusCode)
}

// RegisterMapFromHTTPStatusCode register the mapping relationship form http-status-code to internal error code
func RegisterMapFromHTTPStatusCode(statusCode int, code Code) {
	codeConverter.RegisterMapFromHTTPStatusCode(statusCode, code)
}

// =======================================
// ----- built-in default converter  -----
// =======================================

// defaultCodeConverter
type defaultCodeConverter struct {
	// internalCodeToGRPCCode mapping relationship between internal business error codes and gRPC error codes
	internalCodeToGRPCCode map[Code]codes.Code
	// gRPCCodeToInternalCode mapping relationship between gRPC error codes and internal business error codes
	gRPCCodeToInternalCode map[codes.Code]Code
	// internalCodeToHTTPStatusCode mapping relationship between internal error codes and http-status-codes
	internalCodeToHTTPStatusCode map[Code]int
	// httpStatusCodeToInternalCode mapping relationship between http-status-codes and internal error codes
	httpStatusCodeToInternalCode map[int]Code
}

func newDefaultCodeConverter() *defaultCodeConverter {
	c := &defaultCodeConverter{
		internalCodeToGRPCCode: map[Code]codes.Code{
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
		},
		gRPCCodeToInternalCode: map[codes.Code]Code{},
		internalCodeToHTTPStatusCode: map[Code]int{
			Unknown:            http.StatusInternalServerError,
			OK:                 http.StatusOK,
			InvalidArgument:    http.StatusBadRequest,
			Unauthorized:       http.StatusUnauthorized,
			Forbidden:          http.StatusForbidden,
			NotFound:           http.StatusNotFound,
			RequestTimeout:     http.StatusRequestTimeout,
			ClientClosed:       499,
			InternalError:      http.StatusInternalServerError,
			ServiceUnavailable: http.StatusServiceUnavailable,
			GatewayTimeout:     http.StatusGatewayTimeout,
			AlreadyExists:      614,
		},
		httpStatusCodeToInternalCode: map[int]Code{
			http.StatusOK:                  OK,
			http.StatusBadRequest:          InvalidArgument,
			http.StatusUnauthorized:        Unauthorized,
			http.StatusForbidden:           Forbidden,
			http.StatusNotFound:            NotFound,
			http.StatusRequestTimeout:      RequestTimeout,
			499:                            ClientClosed,
			http.StatusInternalServerError: InternalError,
			http.StatusServiceUnavailable:  ServiceUnavailable,
			http.StatusGatewayTimeout:      GatewayTimeout,
			614:                            AlreadyExists,
		},
	}
	c.initMapBizCodeFromGRPCCode()
	return c
}

// initMapBizCodeFromGRPCCode Recalculate the mapping from gRPC error codes to internal error codes
func (d *defaultCodeConverter) initMapBizCodeFromGRPCCode() {
	for k, v := range d.internalCodeToGRPCCode {
		if _, exist := d.gRPCCodeToInternalCode[v]; !exist {
			d.gRPCCodeToInternalCode[v] = k
		}
	}
}

func (d *defaultCodeConverter) ToGRPCCode(code Code) codes.Code {
	if c, ok := d.internalCodeToGRPCCode[code]; ok {
		return c
	}
	return codes.Unknown
}

func (d *defaultCodeConverter) FromGRPCCode(code codes.Code) Code {
	if c, ok := d.gRPCCodeToInternalCode[code]; ok {
		return c
	}
	return Unknown
}

func (d *defaultCodeConverter) ToHTTPStatusCode(code Code) int {
	if c, ok := d.internalCodeToHTTPStatusCode[code]; ok {
		return c
	}
	return code.ToInt()
}

func (d *defaultCodeConverter) FromHTTPStatusCode(code int) Code {
	if c, ok := d.httpStatusCodeToInternalCode[code]; ok {
		return c
	}
	return defaultCode(code)
}

func (d *defaultCodeConverter) RegisterMapToGRPCCode(code Code, grpcCode codes.Code) {
	d.internalCodeToGRPCCode[code] = grpcCode
}

func (d *defaultCodeConverter) RegisterMapFromGRPCCode(grpcCode codes.Code, code Code) {
	d.gRPCCodeToInternalCode[grpcCode] = code
}

func (d *defaultCodeConverter) RegisterMapToHTTPStatusCode(code Code, statusCode int) {
	d.internalCodeToHTTPStatusCode[code] = statusCode
}

func (d *defaultCodeConverter) RegisterMapFromHTTPStatusCode(statusCode int, code Code) {
	d.httpStatusCodeToInternalCode[statusCode] = code
}
