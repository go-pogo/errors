package errs

import (
	"testing"

	"github.com/roeldev/go-fail"
)

func TestCaptureFrames(t *testing.T) {
	tests := map[string]struct {
		n int
	}{
		"n=0": {n: 0},
		"n=1": {n: 1},
		"n=2": {n: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			frames := CaptureFrames(tc.n, 1)
			if len(frames) != tc.n {
				t.Error(fail.Diff{
					Func: "CaptureFrames",
					Msg:  "should capture n frames",
					Have: len(frames),
					Want: tc.n,
				})
			}
		})
	}
}
