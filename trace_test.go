package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, Trace(nil))
	})
	t.Run("cast to original", func(t *testing.T) {
		err := errors.New("original err")
		traced := Trace(err)

		assert.IsType(t, new(traceErr), traced)
		assert.True(t, errors.Is(traced, err))
	})

	tests := map[string]struct {
		err     error
		wantLen int
	}{
		"with error": {
			err:     New("", ""),
			wantLen: 2,
		},
		"with primitive": {
			err:     errors.New(""),
			wantLen: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Trace(tc.err).(ErrorWithFrames)
			assert.Len(t, *err.Frames(), tc.wantLen)
		})
	}
}
