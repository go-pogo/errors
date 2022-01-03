// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errChainHelper []error

func (h *errChainHelper) append(err error) error {
	*h = append(*h, err)
	return err
}

func (h *errChainHelper) prepend(err error) error {
	*h = prepend(*h, err)
	return err
}

func (h errChainHelper) last() error { return h[len(h)-1] }

func wrappers() map[string]func(parent error) error {
	res := embedders()
	res["WithKind"] = func(parent error) error {
		return WithKind(parent, "some kind")
	}
	return res
}

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

		"base with kind": {baseErr, WithKind(baseErr, "kind")},
		"std with kind":  {stdBase, WithKind(stdBase, "kind")},
	}

	for group, wrapFn := range wrappers() {
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
		rootCause := stderrors.New("root cause")
		withKind := WithKind(rootCause, "some kind")
		withFormatter := WithFormatter(withKind)

		// both upgrades should match with the original error
		assert.ErrorIs(t, withKind, rootCause)
		assert.ErrorIs(t, withFormatter, rootCause)

		// both upgrades should match with each other
		assert.ErrorIs(t, withKind, withFormatter)
		assert.ErrorIs(t, withFormatter, withKind)
	})
}

type customError struct{}

func (ce *customError) Error() string { return "this is a custom error" }

func TestAs(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

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
}
