package errs

import (
	"fmt"
	"runtime"
	"strings"
)

var nilFrame Frame

// Frame is a single step in a stack trace and contains information about the
// function and its package, file and line location.
type Frame struct {
	Path string // Path contains the file path of the function.
	Line int    // Line contains the line number of the called function.
	Func string // Func contains the full name of the called function.
}

// IsEmpty returns true when the Frame is created with all empty fields.
func (f Frame) IsEmpty() bool {
	return f.Path == "" && f.Line == 0 && f.Func == ""
}

// String returns the string representation of the Frame.
func (f Frame) String() string {
	if f.IsEmpty() {
		return ""
	}

	return fmt.Sprintf("%s:%d: %s()", f.Path, f.Line, f.Func)
}

// GetFrame gets a frame from the call stack. Skip indicates the amount of
// frames that have to be skipped before the right frame is to be returned.
// It returns an empty frame with `ok` as `false` when an error occurs.
func GetFrame(skip uint) (frame Frame, ok bool) {
	pc, path, line, ok := runtime.Caller(int(skip + 1))
	if !ok {
		return
	}

	pcFn := runtime.FuncForPC(pc)
	if pcFn == nil {
		return
	}

	frame = Frame{
		Path: path,
		Line: line,
		Func: pcFn.Name(),
	}

	return frame, true
}

func isRuntimeCall(frame Frame) bool {
	if strings.HasPrefix(frame.Func, "runtime.") &&
		strings.Index(frame.Path, "runtime") > 0 {
		return true
	}
	if frame.Func == "testing.tRunner()" &&
		strings.HasSuffix(frame.Path, "testing.go") {
		return true
	}

	return false
}
