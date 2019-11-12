package errs

import (
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

func TestNewStackTrace(t *testing.T) {
	st := NewStackTrace()
	if st.Len() != 0 {
		t.Error(fail.Msg{
			Func: "NewStackTrace",
			Msg:  "should return a pointer to a new empty stack trace",
		})
	}
}

func TestST_Capture(t *testing.T) {
	st := NewStackTrace()

	frame := st.Capture(1)
	if frame == nil || st.Len() != 1 {
		t.Error(fail.Msg{
			Func: "ST.Capture",
			Msg:  "should capture a Frame and add it to the stack",
		})
	}

	frame = st.Capture(1)
	if frame == nil || st.Len() != 2 {
		t.Error(fail.Msg{
			Func: "ST.Capture",
			Msg:  "should capture a second Frame and add it to the stack",
		})
	}

	frame = st.Capture(9999)
	if frame != nil || st.Len() != 2 {
		t.Error(fail.Msg{
			Func: "ST.Capture",
			Msg:  "should only add a Frame to the stack when its not empty",
		})
	}
}

func TestST_Frames(t *testing.T) {
	st := NewStackTrace()
	want := []Frame{*st.Capture(1)}
	have := st.Frames()

	if !reflect.DeepEqual(have, want) {
		t.Error(fail.Diff{
			Func: "ST.Frames",
			Msg:  "should return the captured frame(s)",
			Have: have,
			Want: want,
		})
	}
}
