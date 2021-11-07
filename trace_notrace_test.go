// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build notrace
// +build notrace

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStackFrames(t *testing.T) {
	tests := map[string]error{
		"with nil":          nil,
		"with primitive":    stderrors.New(""),
		"with error":        New(""),
		"with traced error": Trace(New("")),
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Nil(t, GetStackFrames(tc))
		})
	}
}
