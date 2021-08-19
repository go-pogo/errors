// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
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
		"with primitive": {
			err:     stderrors.New(""),
			wantLen: 1,
		},
		"with error": {
			err:     New(""),
			wantLen: 1,
		},
		"with traced error": {
			err:     Trace(New("")),
			wantLen: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Trace(tc.err).(StackTracer)
			assert.Len(t, *err.StackFrames(), tc.wantLen)
		})
	}
}

func TestGetFrames(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		f := *GetStackFrames(New(""))
		assert.Len(t, f, 0)
		assert.NotContains(t, f.String(), "trace_test.go:")
	})
	t.Run("with traced error", func(t *testing.T) {
		f := *GetStackFrames(Trace(New("")))
		assert.Len(t, f, 1)
		assert.Contains(t, f.String(), "trace_test.go:")
	})
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, GetStackFrames(nil))
	})
}
