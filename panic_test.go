// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
)

func panicOnSomething() {
	panic("panic!")
}

func TestWrapPanic(t *testing.T) {
	t.Run("without panic", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		defer WrapPanic("wrapped")
	})

	t.Run("with panic", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "wrapped: panic!", recover())
		}()

		defer WrapPanic("wrapped")
		panicOnSomething()
	})
}

func TestMust(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		var err error
		Must(true, err)
	})

	t.Run("panic on error", func(t *testing.T) {
		errStr := "foo error"
		defer func() {
			assert.Contains(t, recover(), errStr)
		}()

		Must(false, New(errStr))
	})
}

func TestCatchPanic(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	val := struct{ val string }{val: "some value"}

	tests := map[string]struct {
		panic   interface{}
		wantMsg string
	}{
		"with string": {
			panic:   "paniek!",
			wantMsg: "panic: paniek!",
		},
		"with struct": {
			panic:   val,
			wantMsg: fmt.Sprintf("panic: %v", val),
		},
		"with stderror": {
			panic:   stderrors.New("nooo!"),
			wantMsg: "panic: nooo!",
		},
		"with error": {
			panic:   New("panic error"),
			wantMsg: "panic: panic error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var have error
			defer func() {
				assert.Equal(t, newCommonErr(&panicError{tc.panic}, false, 1), have)
				assert.Equal(t, tc.wantMsg, have.Error())
			}()
			defer CatchPanic(&have)
			panic(tc.panic)
		})
	}
}
