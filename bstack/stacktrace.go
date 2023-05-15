package bstack

import (
	"go-brick/internal/bufferpool"
	"go-brick/internal/json"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap/zapcore"
)

// Referenced from: go.uber.org\zap@v1.24.0\stacktrace.go,
// and customize and modify according to this project

var _stacktracePool = sync.Pool{
	New: func() interface{} {
		return &stacktrace{
			storage: make([]uintptr, 64),
		}
	},
}

type stacktrace struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// StacktraceDepth specifies how deep of a stack trace should be captured.
type StacktraceDepth int

const (
	// StacktraceFull captures the entire call stack, allocating more
	// storage for it if needed.
	StacktraceFull = -1
	// StacktraceFirst captures only the first frame.
	StacktraceFirst StacktraceDepth = 1
	// StacktraceMax captures only the first ten frames.
	StacktraceMax StacktraceDepth = 10
)

// captureStacktrace captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// captureStacktrace.
//
// The caller must call Free on the returned stacktrace after using it.
func captureStacktrace(skip int, depth StacktraceDepth) *stacktrace {
	stack := _stacktracePool.Get().(*stacktrace)

	switch depth {
	case StacktraceFirst:
		stack.pcs = stack.storage[:1]
	case StacktraceFull:
		stack.pcs = stack.storage
	case StacktraceMax:
		if len(stack.storage) > int(StacktraceMax) {
			stack.pcs = stack.storage[:StacktraceMax]
		} else {
			stack.pcs = stack.storage
		}
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	if depth == StacktraceFull {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *stacktrace) Free() {
	st.frames = nil
	st.pcs = nil
	_stacktracePool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *stacktrace) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *stacktrace) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

func TakeStack(skip int, depth StacktraceDepth) StackList {
	stack := captureStacktrace(skip+1, depth)
	defer stack.Free()

	stackFmt := newStackFormatter(stack.Count())
	stackFmt.FormatStack(stack)
	return stackFmt.Stack()
}

// stackFormatter formats a stack trace into a readable string representation.
type stackFormatter struct {
	list StackList
}

// newStackFormatter builds a new stackFormatter.
func newStackFormatter(layer int) stackFormatter {
	return stackFormatter{
		list: make(StackList, 0, layer),
	}
}

func (sf *stackFormatter) Stack() StackList {
	return sf.list
}

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *stackFormatter) FormatStack(stack *stacktrace) {
	// nb. On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

// FormatFrame formats the given frame.
func (sf *stackFormatter) FormatFrame(frame runtime.Frame) {
	sf.list = append(sf.list, &stackInfo{
		Func: frame.Function,
		File: frame.File,
		Line: frame.Line,
	})
}

// TrimmedPath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func (sf *stackFormatter) TrimmedPath(file string) string {
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	//
	// Find the last separator.
	//
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file
	}
	buf := bufferpool.Get()
	// Keep everything after the penultimate separator.
	buf.AppendString(file[idx+1:])
	caller := buf.String()
	buf.Free()
	return caller
}

// stackInfo stack info
type stackInfo struct {
	Func string `json:"func"` // function name
	File string `json:"file"` // file name
	Line int    `json:"line"` // line no
}

func (s *stackInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("func", s.Func)
	enc.AddString("file", s.File+":"+strconv.Itoa(s.Line))
	return nil
}

// StackList stack info list
type StackList []*stackInfo

func (sl StackList) MarshalLogArray(enc zapcore.ArrayEncoder) (err error) {
	for _, s := range sl {
		if err = enc.AppendObject(s); err != nil {
			return err
		}
	}
	return
}

func (sl StackList) Error() string {
	out, _ := json.MarshalToString(sl)
	return out
}
