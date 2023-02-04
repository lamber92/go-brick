package bstatus

import (
	"fmt"
	"go-brick/berror/bcode"

	"google.golang.org/grpc/codes"
)

// Business Status preset value
// If you need to add a custom status, you can get it through the NewStatus() method
var (
	Unknown            = &defaultStatus{bcode.Unknown, "Unknown Error", nil}
	OK                 = &defaultStatus{bcode.OK, "Success", nil}
	InvalidArgument    = &defaultStatus{bcode.InvalidArgument, "Invalid Parameters", nil}
	Unauthorized       = &defaultStatus{bcode.Unauthorized, "Not Logged In", nil}
	Forbidden          = &defaultStatus{bcode.Forbidden, "Request Denied", nil}
	NotFound           = &defaultStatus{bcode.NotFound, "Resource Not Found", nil}
	RequestTimeout     = &defaultStatus{bcode.RequestTimeout, "Request Timeout", nil}
	ClientClosed       = &defaultStatus{bcode.ClientClosed, "Client Connection Closed", nil}
	InternalError      = &defaultStatus{bcode.InternalError, "Internal Server Error", nil}
	ServiceUnavailable = &defaultStatus{bcode.ServiceUnavailable, "Service Unavailable", nil}
	GatewayTimeout     = &defaultStatus{bcode.GatewayTimeout, "Gateway Timeout", nil}
	AlreadyExists      = &defaultStatus{bcode.AlreadyExists, "Resource Already Exists", nil}
)

// =======================================
// ---- Default Status Interface IMPL ----
// =======================================

// defaultStatus
// Provide built-in error code carrier
type defaultStatus struct {
	code   bcode.Code  // error code
	reason string      // Error reasons with business attributes. Usually used to return to the client as a prompt.
	detail interface{} // A collection of real underlying error cause information
}

// New create a defaultStatus object pointer
func New(code bcode.Code, reason string, detail interface{}) Status {
	return &defaultStatus{
		code:   code,
		reason: reason,
		detail: detail,
	}
}

func (c *defaultStatus) Code() bcode.Code {
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
var internalCodeMapToStatus = map[bcode.Code]Status{
	bcode.Unknown:            Unknown,
	bcode.OK:                 OK,
	bcode.InvalidArgument:    InvalidArgument,
	bcode.Unauthorized:       Unauthorized,
	bcode.Forbidden:          Forbidden,
	bcode.NotFound:           NotFound,
	bcode.RequestTimeout:     RequestTimeout,
	bcode.ClientClosed:       ClientClosed,
	bcode.InternalError:      InternalError,
	bcode.ServiceUnavailable: ServiceUnavailable,
	bcode.GatewayTimeout:     GatewayTimeout,
	bcode.AlreadyExists:      AlreadyExists,
}

// RegisterMapFromCode
// Register the relationship between custom error codes and error states
func RegisterMapFromCode(code bcode.Code, status Status) {
	internalCodeMapToStatus[code] = status
}

// GetByCode
// Get the registered internal error status by error code.
// the unregistered one returns "unknown status"
func GetByCode(code bcode.Code) Status {
	e, ok := internalCodeMapToStatus[code]
	if ok {
		return e
	}
	return Unknown
}

// RegisterInvalidArgument
// Register custom error status/code for invalid parameters.
// or overwrite the existing error code mapping relationship.
// Notice: The http response body error code does not need to be mapped,
// because it is the http client that needs to obtain this custom error code.
func RegisterInvalidArgument(code bcode.Code, message string, detail interface{}) Status {
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
func RegisterNotFound(code bcode.Code, message string, detail interface{}) Status {
	bcode.RegisterMapToGRPCCode(code, codes.NotFound)
	st := &defaultStatus{
		code:   code,
		reason: message,
		detail: detail,
	}
	internalCodeMapToStatus[code] = st
	return st
}

// RegisterAlreadyExists
// Register custom error status/code for resource already exists.
// or overwrite the existing error code mapping relationship.
// Notice: The http response body error code does not need to be mapped,
// because it is the http client that needs to obtain this custom error code.
func RegisterAlreadyExists(code bcode.Code, message string, detail interface{}) Status {
	bcode.RegisterMapToGRPCCode(code, codes.AlreadyExists)
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
func RegisterInternalError(code bcode.Code, message string, detail interface{}) Status {
	bcode.RegisterMapToGRPCCode(code, codes.Internal)
	st := &defaultStatus{
		code:   code,
		reason: message,
		detail: detail,
	}
	internalCodeMapToStatus[code] = st
	return st
}
