// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFormatter(t *testing.T) {
	t.Run("std error", func(t *testing.T) {
		rootCause := stderrors.New("root cause error")
		have := WithFormatter(rootCause)

		t.Run("add", func(t *testing.T) {
			want := &formatterErr{error: rootCause}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithFormatter(have)
			want := &formatterErr{error: rootCause}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
		})
	})

	t.Run("common error", func(t *testing.T) {
		rootCause := New("root cause error")
		have := WithFormatter(rootCause)

		t.Run("set", func(t *testing.T) {
			want := toCommonErr(Original(rootCause), true)
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithFormatter(have)
			want := toCommonErr(Original(rootCause), true)
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
		})
	})

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithFormatter(nil))
	})
}

func TestFormatError(t *testing.T) {
	tests := map[string]struct {
		setup      func() error
		traceLines []int
	}{
		"error": {
			setup: func() error { return New("some err") },
		},
		"traced primitive": {
			setup:      func() error { return Trace(stderrors.New("primitive")) },
			traceLines: traceHelper(-1, 1),
		},
		"traced error": {
			setup:      func() error { return Trace(New("another err")) },
			traceLines: traceHelper(-1, 1),
		},
		"multi error": {
			setup: func() error {
				err1 := New("err1")
				err2 := Trace(New("err2"))
				return Trace(Combine(err1, err2))
			},
			traceLines: traceHelper(-3, 2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.setup()
			str := fmt.Sprintf("%+v", err)

			for _, line := range tc.traceLines {
				assert.Contains(t, str, "format_test.go:"+strconv.Itoa(line))
			}
		})
	}
}

func traceHelper(offset int, total int) []int {
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return nil
	}

	line += offset

	res := make([]int, 0, total)
	for i := 0; i < total; i++ {
		res = append(res, line+i)
	}
	return res
}
