package errs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCaptureFrames1(n, s uint) Frames { return CaptureFrames(n, s) }
func testCaptureFrames2(n, s uint) Frames { return testCaptureFrames1(n, s) }
func testCaptureFrames3(n, s uint) Frames { return testCaptureFrames2(n, s) }

func TestCaptureFrames(t *testing.T) {
	tests := map[string]struct {
		n       uint
		skip    uint
		wantLen int
	}{
		"n=0": {n: 0, wantLen: 0},
		"n=1": {n: 1, wantLen: 1},
		"n=2": {n: 2, wantLen: 2},
		"n=6": {n: 6, wantLen: 5},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			frames := testCaptureFrames3(tc.n, tc.skip+1)
			assert.Len(t, frames, tc.wantLen)
			assert.Equal(t, uint(tc.wantLen), frames.Len())

			i := tc.wantLen
			if i > 3 {
				i = 3
			}

			s := frames.String()
			for ; i > int(tc.skip); i-- {
				assert.Contains(t, s, "testCaptureFrames"+strconv.Itoa(i))
			}
		})
	}
}

func TestGetFrames(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		f := *GetFrames(New("", ""))
		assert.Len(t, f, 1)
		assert.Contains(t, f.String(), "frames_test.go:")
	})
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, GetFrames(nil))
	})
}
