package bstatus

import "go-brick/berror/bcode"

// Status Carrier of business error info
type Status interface {
	Code() bcode.Code // error code
	Reason() string   // error description
	Detail() any      // error extension
}
