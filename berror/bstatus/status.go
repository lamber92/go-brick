package bstatus

import (
	"fmt"
	"go-brick/berror"
	"go-brick/berror/bcode"

	"google.golang.org/grpc/codes"
)

// Business Status preset value
// If you need to add a custom status, you can get it through the NewStatus() method
var (
	StatusUnknown            = &defaultStatus{bcode.Unknown, "Unknown Error", nil}
	StatusOK                 = &defaultStatus{bcode.OK, "Success", nil}
	StatusInvalidArgument    = &defaultStatus{bcode.InvalidArgument, "Invalid Parameters", nil}
	StatusUnauthorized       = &defaultStatus{bcode.Unauthorized, "Not Logged In", nil}
	StatusForbidden          = &defaultStatus{bcode.Forbidden, "Request Denied", nil}
	StatusNotFound           = &defaultStatus{bcode.NotFound, "Resource Not Found", nil}
	StatusRequestTimeout     = &defaultStatus{bcode.RequestTimeout, "Request Timeout", nil}
	StatusClientClosed       = &defaultStatus{bcode.ClientClosed, "Client Connection Closed", nil}
	StatusInternalError      = &defaultStatus{bcode.InternalError, "Internal Server Error", nil}
	StatusServiceUnavailable = &defaultStatus{bcode.ServiceUnavailable, "Service Unavailable", nil}
	StatusGatewayTimeout     = &defaultStatus{bcode.GatewayTimeout, "Gateway Timeout", nil}
	StatusAlreadyExists      = &defaultStatus{bcode.AlreadyExists, "Resource Already Exists", nil}
)

// =======================================
// ---- Default Status Interface IMPL ----
// =======================================

// defaultStatus
// Provide built-in error code carrier
type defaultStatus struct {
	code   int         // error code
	reason string      // Error reasons with business attributes. Usually used to return to the client as a prompt.
	detail interface{} // A collection of real underlying error cause information
}

// New
// Create a defaultStatus object pointer
func New(code int, reason string, detail interface{}) berror.Status {
	return &defaultStatus{code, reason, detail}
}

func (c *defaultStatus) Code() int {
	return c.code
}

func (c *defaultStatus) Reason() string {
	return c.reason
}

func (c *defaultStatus) Detail() interface{} {
	return c.detail
}

func (c *defaultStatus) String() string {
	if c.detail != nil {
		return fmt.Sprintf("[%d]:%s. (detail: %v)", c.code, c.reason, c.detail)
	}
	if c.reason != "" {
		return fmt.Sprintf("[%d]:%s", c.code, c.reason)
	}
	return fmt.Sprintf("[%d]", c.code)
}

// =======================================
// ----- Default Internal Status Hub -----
// =======================================

// internalCodeMapToStatus
// Registered internal error status container.
// First register the internal error state as a preset value.
var internalCodeMapToStatus = map[int]berror.Status{
	bcode.Unknown:            StatusUnknown,
	bcode.OK:                 StatusOK,
	bcode.InvalidArgument:    StatusInvalidArgument,
	bcode.Unauthorized:       StatusUnauthorized,
	bcode.Forbidden:          StatusForbidden,
	bcode.NotFound:           StatusNotFound,
	bcode.RequestTimeout:     StatusRequestTimeout,
	bcode.ClientClosed:       StatusClientClosed,
	bcode.InternalError:      StatusInternalError,
	bcode.ServiceUnavailable: StatusServiceUnavailable,
	bcode.GatewayTimeout:     StatusGatewayTimeout,
	bcode.AlreadyExists:      StatusAlreadyExists,
}

// RegisterMapFromCode
// Register the relationship between custom error codes and error states
func RegisterMapFromCode(code int, status berror.Status) {
	internalCodeMapToStatus[code] = status
}

// GetByCode
// Get the registered internal error status by error code.
// the unregistered one returns "unknown status"
func GetByCode(code int) berror.Status {
	e, ok := internalCodeMapToStatus[code]
	if ok {
		return e
	}
	return StatusUnknown
}

// RegisterInvalidArgument
// Register custom error status/code for invalid parameters.
// or overwrite the existing error code mapping relationship.
// Notice: The http response body error code does not need to be mapped,
// because it is the http client that needs to obtain this custom error code.
func RegisterInvalidArgument(code int, message string, detail interface{}) berror.Status {
	bcode.RegisterMapToGRPCCode(code, codes.InvalidArgument)
	st := &defaultStatus{
		code:   code,
		reason: message,
		detail: detail,
	}
	internalCodeMapToStatus[code] = st
	return st
}

// RegisterNotFound
// Register custom error status/code for resource not found.
// or overwrite the existing error code mapping relationship.
// Notice: The http response body error code does not need to be mapped,
// because it is the http client that needs to obtain this custom error code.
func RegisterNotFound(code int, message string, detail interface{}) berror.Status {
	bcode.RegisterMapToGRPCCode(code, codes.NotFound)
	st := &defaultStatus{
		code:   code,
		reason: message,
		detail: detail,
	}
	internalCodeMapToStatus[code] = st
	return st
}

// RegisterInternalError
// Register custom error status/code for internal error.
// or overwrite the existing error code mapping relationship.
// Notice: The http response body error code does not need to be mapped,
// because it is the http client that needs to obtain this custom error code.
func RegisterInternalError(code int, message string, detail interface{}) berror.Status {
	bcode.RegisterMapToGRPCCode(code, codes.Internal)
	st := &defaultStatus{
		code:   code,
		reason: message,
		detail: detail,
	}
	internalCodeMapToStatus[code] = st
	return st
}
