package errs

// NewStackTrace creates a new empty stack trace `ST`.
func NewStackTrace() *ST {
	return &ST{
		frames: make([]Frame, 0),
	}
}

// NewStackTrace creates a new empty stack trace `ST` and captures a call frame.
// See `ST.Capture()` for additional information about capturing cal frames.
func NewStackTraceCapture(skip uint) *ST {
	st := NewStackTrace()
	st.Capture(skip)

	return st
}

// ST is a stack of Frames from innermost (newest) to outermost (oldest).
type ST struct {
	frames []Frame
}

// Capture a call frame and prepend it to the stack trace.
func (st *ST) Capture(skip uint) (frame Frame, ok bool) {
	frame, ok = GetFrame(skip + 1)
	if ok {
		st.frames = append(st.frames, frame)
	}
	return
}

// CaptureMultiple captures multiple call frames and prepends them to the stack trace.
func (st *ST) CaptureMultiple(skip uint, amount uint, includeRuntime bool) []Frame {
	frames := make([]Frame, amount)

	var i uint
	for i = 0; i < amount; i++ {
		frame, ok := GetFrame(skip + 1)
		if !ok || (!includeRuntime && isRuntimeCall(frame)) {
			break
		}

		st.frames = prepend(st.frames, frame)
		frames = append(frames, frame)
		skip++
	}

	return frames
}

// Len returns the amount of frames within this stack trace.
func (st ST) Len() uint {
	return uint(len(st.frames))
}

// Frames returns all captured frames.
func (st ST) Frames() []Frame {
	return st.frames
}

func prepend(slice []Frame, frame Frame) []Frame {
	slice = append(slice, Frame{})
	copy(slice[1:], slice)
	slice[0] = frame

	return slice
}
