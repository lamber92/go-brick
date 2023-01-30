package bstatus

import "net/http"

// Business error code preset value
// Note: It is filled with the value of http-status-code, which has nothing to do with the purpose of http-status-code
const (
	Unknown            = -1                             // -1:Unknown error
	OK                 = 0                              // 200:Success
	InvalidArgument    = http.StatusBadRequest          // 400:Invalid parameters
	Unauthorized       = http.StatusUnauthorized        // 401:User login authentication failed
	Forbidden          = http.StatusForbidden           // 403:Request denied
	NotFound           = http.StatusNotFound            // 404:Resource not found
	RequestTimeout     = http.StatusRequestTimeout      // 408:Request timeout
	ClientClosed       = 499                            // 499:Client connection closed
	InternalError      = http.StatusInternalServerError // 500:Internal server error
	ServiceUnavailable = http.StatusServiceUnavailable  // 503:Service unavailable
	GatewayTimeout     = http.StatusGatewayTimeout      // 504:Gateway timeout
	AlreadyExists      = 10001                          // 10001:Resource already exists
)
