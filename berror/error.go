package berror

import (
	"errors"
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
	"go-brick/bstack"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap/zapcore"
)

var jsonStdIter = jsoniter.ConfigCompatibleWithStandardLibrary

// defaultError
// Provide built-in error status carrier
type defaultError struct {
	err    error            // original error
	status bstatus.Status   // business information
	stack  bstack.StackList // stack information when this object(*defaultError) was created
}

// New create and return an error containing a code and reason.
// If the parameter 'err' is passed in, it will wrap err.
// nb. Nesting again will result in inaccurate stack cheapness,
// if necessary, use NewWithSkip instead.
func New(status bstatus.Status, err ...error) Error {
	e := &defaultError{
		stack:  bstack.TakeStack(1, bstack.StacktraceMax),
		status: status,
	}
	if len(err) > 0 {
		e.err = err[0]
	}
	return e
}

// Error output error information in string format
func (d *defaultError) Error() string {
	if d == nil {
		return ""
	}
	str, _ := jsonStdIter.MarshalToString(d.format())
	return str
}

// Status get main status
func (d *defaultError) Status() bstatus.Status {
	if d == nil {
		return bstatus.Unknown
	}
	return d.status
}

// Stack list the error tracking information that has been collected
func (d *defaultError) Stack() bstack.StackList {
	if d == nil || d.stack == nil || len(d.stack) == 0 {
		return bstack.StackList{}
	}
	return d.stack
}

// Wrap nest the specified error into error chain.
// Notice: will overwrite the original internal error
func (d *defaultError) Wrap(err error) error {
	if d == nil {
		return nil
	}
	d.err = err
	return d
}

// Unwrap returns the next error in the error chain.
func (d *defaultError) Unwrap() error {
	if d == nil {
		return nil
	}
	return d.err
}

// Is reports whether any error in error chain matches target.
func (d *defaultError) Is(target error) bool {
	return errors.Is(d, target)
}

// As finds the first error in error chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func (d *defaultError) As(target any) bool {
	return errors.As(d, target)
}

type summary struct {
	Code   bcode.Code `json:"code"`
	Reason string     `json:"reason"`
	Detail any        `json:"detail"`
	Next   any        `json:"next"`
}

func (d *defaultError) format() *summary {
	if d == nil || d.status == nil {
		return nil
	}
	sum := &summary{
		Code:   d.status.Code(),
		Reason: d.status.Reason(),
		Detail: d.status.Detail(),
	}
	switch next := d.err.(type) {
	case *defaultError:
		sum.Next = next.format()
	default:
		sum.Next = next.Error()
	}
	return sum
}

// MarshalLogObject zapcore.ObjectMarshaler impl
func (d *defaultError) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	// code/reason
	status := d.status
	enc.AddInt("code", status.Code().ToInt())
	enc.AddString("reason", status.Reason())
	// detail
	if status.Detail() != nil {
		if obj, ok := status.Detail().(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject("detail", obj)
		} else {
			_ = enc.AddReflected("detail", status.Detail())
		}
	}
	// nest error
	if d.err == nil {
		return
	}
	if next, ok := d.err.(*defaultError); ok {
		_ = enc.AddObject("next", next)
		return
	}
	enc.AddString("next", d.err.Error())
	return
}

// NewWithSkip
// create and return an error containing the stack trace.
// @offset: offset stack depth
func NewWithSkip(err error, status bstatus.Status, skip int) Error {
	return &defaultError{
		err:    err,
		status: status,
		stack:  bstack.TakeStack(skip+1, bstack.StacktraceMax),
	}
}

// NewInvalidArgument create a invalid argument error
func NewInvalidArgument(err error, reason string, detail ...any) error {
	var ds any = nil
	if len(detail) > 0 {
		ds = detail[0]
	}
	return NewWithSkip(err, bstatus.New(bcode.InvalidArgument, reason, ds), 1)
}

// NewNotFound create a not found error
func NewNotFound(err error, reason string, detail ...any) error {
	var ds any = nil
	if len(detail) > 0 {
		ds = detail[0]
	}
	return NewWithSkip(err, bstatus.New(bcode.NotFound, reason, ds), 1)
}

// NewInternalError create a internal error
func NewInternalError(err error, reason string, detail ...any) error {
	var ds any = nil
	if len(detail) > 0 {
		ds = detail[0]
	}
	return NewWithSkip(err, bstatus.New(bcode.InternalError, reason, ds), 1)
}
