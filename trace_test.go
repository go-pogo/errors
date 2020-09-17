package errors

import (
	stderrors "errors"
	"os"
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
		"with primitive": {
			err:     stderrors.New(""),
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

func TestTraceErr_Is(t *testing.T) {
	err := stderrors.New("original err")
	traced := Trace(err).(*traceErr)

	assert.True(t, traced.Is(err))
	assert.True(t, stderrors.Is(traced, err))
}

func TestTraceErr_As(t *testing.T) {
	err := new(os.PathError)
	trace := Trace(err).(*traceErr)

	var t1 *os.PathError
	assert.True(t, trace.As(&t1))
	assert.Exactly(t, err, t1)

	var t2 *os.PathError
	assert.True(t, stderrors.As(trace, &t2))
	assert.Exactly(t, t1, t2)
}
