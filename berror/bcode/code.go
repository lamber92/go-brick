package bcode

import (
	"net/http"
	"strconv"
)

type defaultCode int

// New create a defaultCode object
func New(code int) Code {
	return defaultCode(code)
}

// internal Error Code preset value
// nb. it is filled with the value of http-status-code, which has nothing to do with the purpose of http-status-code
const (
	Unknown            defaultCode = -1                             // -1:Unknown error
	OK                 defaultCode = 0                              // 200:Success
	InvalidArgument    defaultCode = http.StatusBadRequest          // 400:Invalid parameters
	Unauthorized       defaultCode = http.StatusUnauthorized        // 401:User login authentication failed
	Forbidden          defaultCode = http.StatusForbidden           // 403:Request denied
	NotFound           defaultCode = http.StatusNotFound            // 404:Resource not found
	RequestTimeout     defaultCode = http.StatusRequestTimeout      // 408:Request timeout
	ClientClosed       defaultCode = 499                            // 499:Client connection closed
	InternalError      defaultCode = http.StatusInternalServerError // 500:Internal server error
	ServiceUnavailable defaultCode = http.StatusServiceUnavailable  // 503:Service unavailable
	GatewayTimeout     defaultCode = http.StatusGatewayTimeout      // 504:Gateway timeout
	AlreadyExists      defaultCode = 614                            // 614:Resource already exists
)

func (c defaultCode) ToInt() int {
	return int(c)
}

func (c defaultCode) ToString() string {
	return strconv.Itoa(int(c))
}

// Is compare the target code value with the current error code value
// only Code, Integer and String types are supported
func (c defaultCode) Is(target any) bool {
	switch tmp := target.(type) {
	case Code:
		return c.ToInt() == tmp.ToInt()
	case int:
		return c.ToInt() == tmp
	case int8:
		return c.ToInt() == int(tmp)
	case int32:
		return c.ToInt() == int(tmp)
	case int64:
		return c.ToInt() == int(tmp)
	case uint:
		return c.ToInt() == int(tmp)
	case uint8:
		return c.ToInt() == int(tmp)
	case uint32:
		return c.ToInt() == int(tmp)
	case uint64:
		return c.ToInt() == int(tmp)
	case string:
		return c.ToString() == tmp
	}
	return false
}
