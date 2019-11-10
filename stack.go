package errs

// ST is a stack of Frames from innermost (newest) to outermost (oldest).
type ST struct {
	frames []Frame
}

// Capture a call frame and prepend it to the stack trace.
func (st *ST) Capture(skip uint) {
	frame, _ := GetFrame(skip + 1)
	st.frames = prepend(st.frames, frame)
}

// Frames returns all captured frames.
func (st ST) Frames() []Frame {
	return st.frames
}

func prepStackTrace(st *ST) *ST {
	st.frames = []Frame{}
	return st
}

// NewStackTrace creates a new stack trace ST.
func NewStackTrace() *ST {
	var st ST
	return prepStackTrace(&st)
}

// prepend Frame to slice.
func prepend(slice []Frame, frame Frame) []Frame {
	slice = append(slice, Frame{})
	copy(slice[1:], slice)
	slice[0] = frame

	return slice
}
