package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, Trace(nil))
	})

	tests := map[string]struct {
		err     error
		wantLen int
	}{
		"with error": {
			err:     New("", ""),
			wantLen: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Trace(tc.err).(ErrorWithFrames)
			assert.Len(t, *err.Frames(), tc.wantLen)
		})
	}
}
