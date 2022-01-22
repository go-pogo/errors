// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !notrace
// +build !notrace

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStackTrace(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, GetStackTrace(nil))
	})
	t.Run("with std error", func(t *testing.T) {
		assert.Nil(t, GetStackTrace(stderrors.New("err")))
	})

	tests := map[string]error{
		"with error":                New("err"),
		"with std error with stack": WithStack(stderrors.New("err")),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			f := GetStackTrace(err)
			assert.Len(t, f.Frames(), 1)
			assert.Contains(t, f.String(), "trace_test.go:")
		})
	}
}
