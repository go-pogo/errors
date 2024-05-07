// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithStack(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, WithStack(nil))
	})
	t.Run("with Msg", func(t *testing.T) {
		assert.PanicsWithValue(t, panicUseNewInstead, func() {
			_ = WithStack(Msg("panic!"))
		})
	})
	t.Run("with std error", func(t *testing.T) {
		err := stderrors.New("some err")
		have := WithStack(err)
		assert.ErrorIs(t, have, err)

		if !internal.TraceStack {
			return
		}
		if !assert.Len(t, have.StackTrace().Frames(), 1) {
			fmt.Printf("\n%+v\n", have)
		}
	})
	t.Run("with error", func(t *testing.T) {
		err := New("some err")
		have := WithStack(err)
		assert.Same(t, err, have)
	})
}
