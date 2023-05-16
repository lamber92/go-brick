package berror

import (
	"errors"
	"fmt"
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
	"net"
	"os"
	"strings"
)

const recoverReason = "recover"

var _identifyErr = defaultIdentify

func ReplaceRecoverIdentify(f func(r any, hook func(err error))) {
	_identifyErr = f
}

// Recover catching and recovering from panics
func Recover(hook func(error)) {
	if r := recover(); r != nil {
		_identifyErr(r, hook)
	}
}

// defaultIdentify the default processing method for identifying panic reasons
func defaultIdentify(r any, hook func(error)) {
	var (
		err    error
		status bstatus.Status
	)
	switch tmp := r.(type) {
	// borrowed from github.com\gin-gonic\gin@v1.7.1\recovery.go
	// ignore specific network errors
	case *net.OpError:
		err = tmp
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		if se, ok := tmp.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				status = bstatus.New(bcode.ClientClosed, recoverReason, ".(type)=*net.OpError")
				break
			}
		}
		status = bstatus.New(bcode.InternalError, recoverReason, ".(type)=*net.OpError")
	case string:
		err = errors.New(tmp)
		status = bstatus.New(bcode.InternalError, recoverReason, ".(type)=string")
	case error:
		err = tmp
		status = bstatus.New(bcode.InternalError, recoverReason, ".(type)=error")
	default:
		err = fmt.Errorf("%+v", r)
		status = bstatus.New(bcode.InternalError, recoverReason, fmt.Sprintf(".(type)=%T", tmp))
	}
	// convert to internal error and skip stacktrace layer
	err = NewWithSkip(err, status, 3)
	hook(err)
}
