// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.20
// +build go1.20

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewf2(t *testing.T) {
	t.Run("with causes", func(t *testing.T) {
		cause1 := New("some err")
		cause2 := New("another err")
		have := Errorf("whoops: %w and %w", cause1, cause2)
		assert.ErrorIs(t, have, cause1)
		assert.ErrorIs(t, have, cause2)
		assert.Equal(t, []error{cause1, cause2}, have.(MultiError).Unwrap())
	})
}
