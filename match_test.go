// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

type errChainHelper []error

//goland:noinspection GoMixedReceiverTypes
func (h *errChainHelper) append(err error) error {
	*h = append(*h, err)
	return err
}

//goland:noinspection GoMixedReceiverTypes
func (h *errChainHelper) prepend(err error) error {
	*h = append([]error{err}, *h...)
	return err
}

//goland:noinspection GoMixedReceiverTypes
func (h errChainHelper) last() error { return h[len(h)-1] }

// test if the root cause error matches all wrapping errors in the chain
func TestIs(t *testing.T) {
	baseErr := New("root cause")
	wrapErr := Wrap(baseErr, "its a wrap")
	stdBase := stderrors.New("root cause")
	stdWrap := fmt.Errorf("error: %w", stdBase)

	chains := map[string]errChainHelper{
		"base":      {baseErr},
		"wrap":      {baseErr, wrapErr},
		"std error": {stdBase},
		"std wrap":  {stdBase, stdWrap},

		"std with stack": {stdBase, WithStack(stdBase)},
	}

	for group, wrapFn := range provideEmbedders() {
		t.Run(group, func(t *testing.T) {
			for name, chain := range chains {
				t.Run(name, func(t *testing.T) {
					// pass the last error to the function we'd like  to test
					err := chain.append(wrapFn(chain.last()))
					assert.Same(t, chain[0], Cause(err))

					for i, target := range chain {
						for j := i; j < len(chain); j++ {
							err = chain[j]
							assert.ErrorIs(t, err, target)
						}
					}
				})
			}
		})
	}

	t.Run("manual", func(t *testing.T) {
		cause := stderrors.New("root cause")
		wrapped := Wrap(cause, Msg("second"))
		withFormatter := WithFormatter(wrapped)

		// both upgrades should match with the original error
		assert.ErrorIs(t, wrapped, cause)
		assert.ErrorIs(t, withFormatter, cause)

		// both upgrades should match with each other
		assert.ErrorIs(t, wrapped, withFormatter)
		assert.ErrorIs(t, withFormatter, wrapped)
	})

	t.Run("multi", func(t *testing.T) {
		internal.DisableTraceStack()
		defer internal.EnableTraceStack()

		err1 := stderrors.New("some err")
		err2 := New("whoops")
		multi := newMultiErr([]error{err2, err1}, 0)
		assert.True(t, Is(multi, err1))
		assert.True(t, Is(multi, err2))
		assert.False(t, Is(multi, stderrors.New("some err")))
	})
}

type customError struct{}

func (ce *customError) Error() string { return "this is a custom error" }

func TestAs(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	var customErr *customError
	var pathErrPtr *os.PathError
	_, pathErr := os.Open("non-existing")

	tests := map[string]struct {
		error  error
		target interface{}
		wantFn func(err error) interface{}
	}{
		"traced os.PathError": {
			error:  WithStack(pathErr),
			target: &pathErrPtr,
			wantFn: func(err error) interface{} {
				return pathErr
			},
		},
		"traced custom error": {
			error:  WithStack(&customError{}),
			target: &customErr,
			wantFn: func(err error) interface{} {
				return &customError{}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			val := reflect.ValueOf(tc.target)
			val.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))

			assert.True(t, As(tc.error, tc.target))

			got := val.Elem().Interface()
			want := tc.wantFn(tc.error)

			assert.Equal(t, want, got)
		})
	}

	t.Run("multi", func(t *testing.T) {
		internal.DisableTraceStack()
		defer internal.EnableTraceStack()

		err1 := stderrors.New("some err")
		err2 := New("whoops")
		multi := newMultiErr([]error{err2, err1}, 0)

		var have commonError
		assert.True(t, As(multi, &have))
		assert.Equal(t, err2, &have)

		var have2 Msg
		assert.False(t, As(multi, &have2))
	})
}
