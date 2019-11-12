package errs

// ST is a stack of Frames from innermost (newest) to outermost (oldest).
type ST struct {
	frames []Frame
}

// Capture a call frame and prepend it to the stack trace.
func (st *ST) Capture(skip uint) *Frame {
	frame, ok := GetFrame(skip + 1)
	if ok {
		st.frames = prepend(st.frames, frame)
		return &frame
	}

	return nil
}

// Len returns the amount of frames within this stack trace.
func (st ST) Len() uint {
	return uint(len(st.frames))
}

// Frames returns all captured frames.
func (st ST) Frames() []Frame {
	return st.frames
}

// NewStackTrace creates a new empty stack trace ST.
func NewStackTrace() *ST {
	var st ST
	return prepStackTrace(&st)
}

func prepStackTrace(st *ST) *ST {
	st.frames = []Frame{}
	return st
}

func prepend(slice []Frame, frame Frame) []Frame {
	slice = append(slice, Frame{})
	copy(slice[1:], slice)
	slice[0] = frame

	return slice
}
