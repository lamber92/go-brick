package berror

import (
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
	"sync"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var defConv Converter = newDefaultConverter()

// ReplaceConverter replace built-in Converter impl
func ReplaceConverter(c Converter) {
	defConv = c
}

// RegisterConvHook register hook function.
// nb. if you need this Hook, call it when you initialize the program.
func RegisterConvHook(hook HookFunc) {
	defConv.Hook(hook)
}

// Convert according to the built-in rules, convert the incoming error to Error.
func Convert(err error, reason string, detail ...any) error {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	return defConv.Convert(err, reason, d)
}

// ConvertWithOption according to the built-in rules, convert the incoming error to Error.
func ConvertWithOption(err error, reason string, detail any, options ...ConvOption) error {
	return defConv.Convert(err, reason, detail, options...)
}

type defaultConverter struct {
	once sync.Once
	hook HookFunc
}

func newDefaultConverter() *defaultConverter {
	return &defaultConverter{}
}

// Convert when all error types are encountered, they are automatically wrapped as Error types.
// if you don't want to deal with the Error type, use the IgnoreWrapError option.
func (dc *defaultConverter) Convert(err error, reason string, detail any, options ...ConvOption) error {
	if err == nil {
		return NewWithSkip(err, bstatus.New(bcode.OK, reason, detail), 1)
	}
	// check it's berror.Error or not
	if orig, ok := err.(Error); ok {
		// check ignore warping.
		for _, v := range options {
			if v2, ok2 := v.(*defaultOption); ok2 && v2.Code == ignoreWrap {
				return err
			}
		}
		// keep the original Code and generate a new Error to wrap err.
		code := orig.Status().Code()
		return NewWithSkip(err, bstatus.New(code, reason, detail), 1)
	}
	if dc.hook != nil {
		return dc.hook(err, reason, detail, options...)
	}

	// !!! built-in processing logic !!!

	// check gorm/redis error
	switch err {
	case gorm.ErrRecordNotFound, redis.Nil:
		return NewWithSkip(err, bstatus.New(bcode.NotFound, reason, detail), 1)
	}
	// check it's grpc error or not
	if gerr, ok := status.FromError(err); ok && gerr != nil {
		code := bcode.FromGRPCCode(gerr.Code())
		return NewWithSkip(err, bstatus.New(code, reason, detail), 1)
	}
	// unknown error
	return NewWithSkip(err, bstatus.New(bcode.Unknown, reason, detail), 1)
}

func (dc *defaultConverter) Hook(f HookFunc) {
	dc.once.Do(func() {
		dc.hook = f
	})
}

const (
	ignoreWrap = 1
)

type defaultOption struct {
	Code int
}

// IgnoreWrapError if original error type is berror.Error,
// do not wrap the original error, return directly.
// only valid for berror.Error
func IgnoreWrapError() ConvOption {
	return &defaultOption{Code: ignoreWrap}
}
